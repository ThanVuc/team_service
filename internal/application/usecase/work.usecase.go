package usecase

import (
	"context"
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
)

type workUseCase struct {
	store      istore.Store
	workRepo   irepository.WorkRepository
	validator  *appvalidation.WorkValidator
	authHelper *apphelper.AuthHelper
}

func (uc *workUseCase) CreateWork(ctx context.Context, req *appdto.CreateWorkRequest) (*appdto.BaseResponse[appdto.WorkResponse], errorbase.AppError) {
	_, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleMember)
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

	works, err := uc.workRepo.GetWorksBySprint(ctx, payload.GroupID, payload.SprintID)
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

	var updatedWork *appdto.WorkResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		updatedWork, err = repo.WorkRepository().UpdateWork(ctx, payload.Request)
		if err != nil {
			return err
		}

		if updatedWork == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("update work returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	updatedWork, err = uc.workRepo.GetWorkAggregation(ctx, payload.WorkID)
	if err != nil {
		return nil, err
	}

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
