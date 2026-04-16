package usecase

import (
	"context"
	"fmt"
	adapterdomain "team_service/internal/adapter/constant/domain"
	appconstant "team_service/internal/application/common/constant"
	appdto "team_service/internal/application/common/dto"
	apphelper "team_service/internal/application/common/helper"
	irepository "team_service/internal/application/common/interface/repository"
	istore "team_service/internal/application/common/interface/store"
	appmapper "team_service/internal/application/common/mapper"
	appvalidation "team_service/internal/application/common/validation"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"team_service/internal/infrastructure/share/utils"
)

type workUseCase struct {
	store              istore.Store
	workRepo           irepository.WorkRepository
	groupRepo          irepository.GroupRepository
	validator          *appvalidation.WorkValidator
	authHelper         *apphelper.AuthHelper
	notificationHelper *apphelper.NotificationHelper
}

func (uc *workUseCase) CreateWork(ctx context.Context, req *appdto.CreateWorkRequest) (*appdto.BaseResponse[appdto.WorkResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleMember)
	if err != nil {
		return &appdto.BaseResponse[appdto.WorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	work, err := uc.validator.ValidateCreateWork(ctx, req)
	if err != nil {
		return &appdto.BaseResponse[appdto.WorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var createdWork *entity.Work
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		createdWork, err = repo.WorkRepository().CreateWork(ctx, work)
		if err != nil {
			return err
		}

		if createdWork == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("create work returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	link := fmt.Sprintf("%s/groups/%s/works/%s", adapterdomain.Domain, createdWork.GroupID, createdWork.ID)
	receivers := []string{actor.ID}
	if createdWork.AssigneeID != "" && createdWork.AssigneeID != actor.ID {
		receivers = append(receivers, createdWork.AssigneeID)
	}
	_ = uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
		EventType:   appconstant.EventTypeWorkCreated,
		SenderID:    actor.ID,
		ReceiverIDs: receivers,
		Payload: appdto.TeamNotificationMessagePayload{
			Title:           appconstant.GetDisplayTitle(appconstant.EventTypeWorkCreated),
			Message:         fmt.Sprintf("Bạn đã tạo công việc %s", createdWork.Name),
			Link:            utils.Ptr(link),
			ImageURL:        nil,
			CorrelationID:   createdWork.GroupID,
			CorrelationType: int(appconstant.CorrelationTypeWork),
		},
		Metadata: appdto.TeamNotificationMessageMetadata{
			IsSentMail:           false,
			NonExistentReceivers: []string{},
		},
	}, &appdto.UserWithPermission{
		ID:                   actor.ID,
		HasEmailNotification: actor.HasEmailNotification,
		HasPushNotification:  actor.HasPushNotification,
	})

	return &appdto.BaseResponse[appdto.WorkResponse]{
		Data:  appmapper.ToWorkResponse(createdWork),
		Error: nil,
	}, nil
}

func (uc *workUseCase) GetWork(ctx context.Context, req *appdto.GetWorkRequest) (*appdto.BaseResponse[appdto.WorkResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleViewer)
	if err != nil {
		return &appdto.BaseResponse[appdto.WorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	payload, err := uc.validator.ValidateGetWork(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.WorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	work, err := uc.workRepo.GetWorkAggregation(ctx, payload.WorkID)
	if err != nil {
		return &appdto.BaseResponse[appdto.WorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	if work == nil || work.GroupID != payload.GroupID {
		notFoundErr := errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("work not found"))
		return &appdto.BaseResponse[appdto.WorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(notFoundErr),
		}, nil
	}

	return &appdto.BaseResponse[appdto.WorkResponse]{
		Data:  work,
		Error: nil,
	}, nil
}

func (uc *workUseCase) ListWorks(ctx context.Context, req *appdto.ListWorksRequest) (*appdto.BaseResponse[appdto.ListWorksResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleViewer)
	if err != nil {
		return &appdto.BaseResponse[appdto.ListWorksResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	payload, err := uc.validator.ValidateListWorks(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.ListWorksResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	works, err := uc.workRepo.GetWorksBySprint(ctx, payload.GroupID, payload.SprintID, payload.AssigneeID)
	if err != nil {
		return &appdto.BaseResponse[appdto.ListWorksResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	if works == nil {
		works = make([]appdto.WorkResponse, 0)
	}

	return &appdto.BaseResponse[appdto.ListWorksResponse]{
		Data: &appdto.ListWorksResponse{
			Works: works,
		},
		Error: nil,
	}, nil
}

func (uc *workUseCase) UpdateWork(ctx context.Context, req *appdto.UpdateWorkRequest) (*appdto.BaseResponse[appdto.WorkResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleMember)
	if err != nil {
		return &appdto.BaseResponse[appdto.WorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	payload, err := uc.validator.ValidateUpdateWork(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.WorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	previousWork, err := uc.workRepo.GetWorkAggregation(ctx, payload.WorkID)
	if err != nil {
		return &appdto.BaseResponse[appdto.WorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	if previousWork == nil {
		notFoundErr := errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("work not found"))
		return &appdto.BaseResponse[appdto.WorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(notFoundErr),
		}, nil
	}

	var updatedWork *appdto.WorkResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		updatedWorkResp, err := repo.WorkRepository().UpdateWork(ctx, payload.Request)
		if err != nil {
			return err
		}

		if updatedWorkResp == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("update work returned nil"))
		}

		updatedWork, err = repo.WorkRepository().GetWorkAggregation(ctx, payload.WorkID)
		if err != nil {
			return err
		}

		if updatedWork == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("get updated work returned nil"))
		}

		return nil
	})

	if err != nil {
		return &appdto.BaseResponse[appdto.WorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	sprintID := ""
	if updatedWork.SprintID != nil {
		sprintID = *updatedWork.SprintID
	}

	link := apphelper.BuildWorkUpdateLink(ctx, updatedWork.GroupID, sprintID, updatedWork.ID)
	previousAssigneeID := ""
	if previousWork.AssigneeID != nil {
		previousAssigneeID = *previousWork.AssigneeID
	}

	currentAssigneeID := ""
	if updatedWork.AssigneeID != nil {
		currentAssigneeID = *updatedWork.AssigneeID
	}

	if currentAssigneeID != "" && currentAssigneeID != previousAssigneeID {
		_ = uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
			EventType:   appconstant.EventTypeWorkAssigned,
			SenderID:    actor.ID,
			ReceiverIDs: []string{currentAssigneeID},
			Payload: appdto.TeamNotificationMessagePayload{
				Title:           appconstant.GetDisplayTitle(appconstant.EventTypeWorkAssigned),
				Message:         fmt.Sprintf("Bạn đã được giao thực hiện công việc '%s'.", updatedWork.Name),
				Link:            utils.Ptr(link),
				ImageURL:        nil,
				CorrelationID:   updatedWork.ID,
				CorrelationType: int(appconstant.CorrelationTypeWork),
			},
			Metadata: appdto.TeamNotificationMessageMetadata{
				IsSentMail:           true,
				NonExistentReceivers: []string{},
			},
		}, nil)
	}

	if updatedWork.Status != previousWork.Status && currentAssigneeID != "" {
		_ = uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
			EventType:   appconstant.EventTypeWorkStatusChanged,
			SenderID:    actor.ID,
			ReceiverIDs: []string{currentAssigneeID},
			Payload: appdto.TeamNotificationMessagePayload{
				Title:           appconstant.GetDisplayTitle(appconstant.EventTypeWorkStatusChanged),
				Message:         fmt.Sprintf("Công việc '%s' đã chuyển sang trạng thái %s.", updatedWork.Name, updatedWork.Status),
				Link:            utils.Ptr(link),
				ImageURL:        nil,
				CorrelationID:   updatedWork.ID,
				CorrelationType: int(appconstant.CorrelationTypeWork),
			},
			Metadata: appdto.TeamNotificationMessageMetadata{
				IsSentMail:           false,
				NonExistentReceivers: []string{},
			},
		}, nil)
	}

	var message string
	message = fmt.Sprintf("Công việc %s đã được cập nhật", updatedWork.Name)

	receivers := []string{actor.ID}
	if updatedWork.AssigneeID != nil {
		if *updatedWork.AssigneeID != "" && *updatedWork.AssigneeID != actor.ID {
			message = fmt.Sprintf("Công việc %s đã được cập nhật và được giao cho %s", updatedWork.Name, updatedWork.Assignee.Email)
			receivers = append(receivers, *updatedWork.AssigneeID)
		}
	}

	_ = uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
		EventType:   appconstant.EventTypeWorkUpdated,
		SenderID:    actor.ID,
		ReceiverIDs: receivers,
		Payload: appdto.TeamNotificationMessagePayload{
			Title:           appconstant.GetDisplayTitle(appconstant.EventTypeWorkUpdated),
			Message:         message,
			Link:            utils.Ptr(link),
			ImageURL:        nil,
			CorrelationID:   updatedWork.GroupID,
			CorrelationType: int(appconstant.CorrelationTypeWork),
		},
		Metadata: appdto.TeamNotificationMessageMetadata{
			IsSentMail:           false,
			NonExistentReceivers: []string{},
		},
	}, &appdto.UserWithPermission{
		ID:                   actor.ID,
		HasEmailNotification: actor.HasEmailNotification,
		HasPushNotification:  actor.HasPushNotification,
	})

	return &appdto.BaseResponse[appdto.WorkResponse]{
		Data:  updatedWork,
		Error: nil,
	}, nil
}

func (uc *workUseCase) DeleteWork(ctx context.Context, req *appdto.DeleteWorkRequest) (*appdto.BaseResponse[appdto.DeleteWorkResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleManager)
	if err != nil {
		return &appdto.BaseResponse[appdto.DeleteWorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	payload, err := uc.validator.ValidateDeleteWork(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.DeleteWorkResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var deletedResult *appdto.DeleteWorkResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		deletedResult, err = repo.WorkRepository().DeleteWork(ctx, payload.WorkID)
		if err != nil {
			return err
		}

		if deletedResult == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("delete work returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	link := fmt.Sprintf("%s/groups/%s/works", adapterdomain.Domain, actor.GroupId)
	_ = uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
		EventType:   appconstant.EventTypeWorkDeleted,
		SenderID:    actor.ID,
		ReceiverIDs: []string{actor.ID},
		Payload: appdto.TeamNotificationMessagePayload{
			Title:           appconstant.GetDisplayTitle(appconstant.EventTypeWorkDeleted),
			Message:         "Công việc đã bị xóa",
			Link:            utils.Ptr(link),
			ImageURL:        nil,
			CorrelationID:   "",
			CorrelationType: int(appconstant.CorrelationTypeWork),
		},
		Metadata: appdto.TeamNotificationMessageMetadata{
			IsSentMail:           false,
			NonExistentReceivers: []string{},
		},
	}, &appdto.UserWithPermission{
		ID:                   actor.ID,
		HasEmailNotification: actor.HasEmailNotification,
		HasPushNotification:  actor.HasPushNotification,
	})

	return &appdto.BaseResponse[appdto.DeleteWorkResponse]{
		Data:  deletedResult,
		Error: nil,
	}, nil
}

func (uc *workUseCase) CreateChecklistItem(ctx context.Context, req *appdto.CreateChecklistItemRequest) (*appdto.BaseResponse[appdto.ChecklistItemResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleMember)
	if err != nil {
		return &appdto.BaseResponse[appdto.ChecklistItemResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	payload, err := uc.validator.ValidateCreateChecklistItem(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.ChecklistItemResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var createdItem *appdto.ChecklistItemResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		createdItem, err = repo.WorkRepository().CreateChecklistItem(ctx, payload.Item)
		if err != nil {
			return err
		}

		if createdItem == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("create checklist item returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &appdto.BaseResponse[appdto.ChecklistItemResponse]{
		Data:  createdItem,
		Error: nil,
	}, nil
}

func (uc *workUseCase) UpdateChecklistItem(ctx context.Context, req *appdto.UpdateChecklistItemRequest) (*appdto.BaseResponse[appdto.ChecklistItemResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleMember)
	if err != nil {
		return &appdto.BaseResponse[appdto.ChecklistItemResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	payload, err := uc.validator.ValidateUpdateChecklistItem(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.ChecklistItemResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var updatedItem *appdto.ChecklistItemResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		updatedItem, err = repo.WorkRepository().UpdateChecklistItem(ctx, payload.Request)
		if err != nil {
			return err
		}

		if updatedItem == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("update checklist item returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &appdto.BaseResponse[appdto.ChecklistItemResponse]{
		Data:  updatedItem,
		Error: nil,
	}, nil
}

func (uc *workUseCase) DeleteChecklistItem(ctx context.Context, req *appdto.DeleteChecklistItemRequest) (*appdto.BaseResponse[appdto.ChecklistItemResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleMember)
	if err != nil {
		return &appdto.BaseResponse[appdto.ChecklistItemResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	payload, err := uc.validator.ValidateDeleteChecklistItem(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.ChecklistItemResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var deletedItem *appdto.ChecklistItemResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		deletedItem, err = repo.WorkRepository().DeleteChecklistItem(ctx, payload.ItemID)
		if err != nil {
			return err
		}

		if deletedItem == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("delete checklist item returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &appdto.BaseResponse[appdto.ChecklistItemResponse]{
		Data:  deletedItem,
		Error: nil,
	}, nil
}

func (uc *workUseCase) CreateComment(ctx context.Context, req *appdto.CreateCommentRequest) (*appdto.BaseResponse[appdto.CommentListResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleMember)
	if err != nil {
		return &appdto.BaseResponse[appdto.CommentListResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	payload, err := uc.validator.ValidateCreateComment(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.CommentListResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var createdComment *appdto.CommentListResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		createdComment, err = repo.WorkRepository().CreateComment(ctx, payload.Comment)
		if err != nil {
			return err
		}

		if createdComment == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("create comment returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	work, err := uc.workRepo.GetWorkAggregation(ctx, payload.WorkID)
	if err != nil {
		return nil, err
	}

	if work != nil {
		receiverIDs := make([]string, 0)
		if work.AssigneeID != nil {
			receiverIDs = append(receiverIDs, *work.AssigneeID)
		}

		receiverIDs = append(receiverIDs, apphelper.CollectDiscussionParticipantIDs(work)...)
		receiverIDs = apphelper.ExcludeID(apphelper.UniqueIDs(receiverIDs...), actor.ID)

		sprintID := ""
		if work.SprintID != nil {
			sprintID = *work.SprintID
		}

		workLink := apphelper.BuildWorkUpdateLink(ctx, work.GroupID, sprintID, work.ID)
		_ = uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
			EventType:   appconstant.EventTypeWorkCommented,
			SenderID:    actor.ID,
			ReceiverIDs: receiverIDs,
			Payload: appdto.TeamNotificationMessagePayload{
				Title:           appconstant.GetDisplayTitle(appconstant.EventTypeWorkCommented),
				Message:         fmt.Sprintf("%s đã bình luận trong '%s': %s", actor.Email, work.Name, payload.Comment.Content),
				Link:            utils.Ptr(workLink),
				ImageURL:        nil,
				CorrelationID:   work.ID,
				CorrelationType: int(appconstant.CorrelationTypeWork),
			},
			Metadata: appdto.TeamNotificationMessageMetadata{
				IsSentMail:           false,
				NonExistentReceivers: []string{},
			},
		}, nil)
	}

	return &appdto.BaseResponse[appdto.CommentListResponse]{
		Data:  createdComment,
		Error: nil,
	}, nil
}

func (uc *workUseCase) UpdateComment(ctx context.Context, req *appdto.UpdateCommentRequest) (*appdto.BaseResponse[appdto.CommentListResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleMember)
	if err != nil {
		return &appdto.BaseResponse[appdto.CommentListResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	payload, err := uc.validator.ValidateUpdateComment(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.CommentListResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var updatedComment *appdto.CommentListResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		updatedComment, err = repo.WorkRepository().UpdateComment(ctx, payload.Request)
		if err != nil {
			return err
		}

		if updatedComment == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("update comment returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &appdto.BaseResponse[appdto.CommentListResponse]{
		Data:  updatedComment,
		Error: nil,
	}, nil
}

func (uc *workUseCase) DeleteComment(ctx context.Context, req *appdto.DeleteCommentRequest) (*appdto.BaseResponse[appdto.CommentListResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleMember)
	if err != nil {
		return &appdto.BaseResponse[appdto.CommentListResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	payload, err := uc.validator.ValidateDeleteComment(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.CommentListResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var deletedComment *appdto.CommentListResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		deletedComment, err = repo.WorkRepository().DeleteComment(ctx, payload.CommentID)
		if err != nil {
			return err
		}

		if deletedComment == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("delete comment returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &appdto.BaseResponse[appdto.CommentListResponse]{
		Data:  deletedComment,
		Error: nil,
	}, nil
}
