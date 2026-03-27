package usecase

import (
	"context"
	"fmt"
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
	"team_service/proto/common"
	"time"

	"github.com/google/uuid"
)

type groupUseCase struct {
	store              istore.Store
	groupRepo          irepository.GroupRepository
	userRepo           irepository.UserRepository
	validator          *appvalidation.GroupValidator
	authHelper         *apphelper.AuthHelper
	notificationHelper *apphelper.NotificationHelper
}

func (uc *groupUseCase) CreateGroup(ctx context.Context, req *appdto.CreateGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError) {
	// validation error -> error in response (422)
	group, user, err := uc.validator.ValidateCreateGroup(ctx, req)
	if err != nil {
		return &appdto.BaseResponse[appdto.GroupResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		group, err = repo.GroupRepository().CreateGroup(ctx, group, user.ID)
		if err != nil {
			return err
		}

		if group == nil {
			return errorbase.New(errdict.ErrNotFound)
		}

		user, err = repo.UserRepository().GetUserByID(ctx, user.ID)
		if err != nil {
			return err
		}

		if user == nil {
			return errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("The user is not found"))
		}

		groupMember, err := entity.NewGroupMember(
			uuid.NewString(), group.ID, user.ID, enum.GroupRoleOwner, time.Now(),
		)

		err = repo.GroupRepository().AddGroupMember(ctx, groupMember)

		if err != nil {
			return err
		}

		return nil
	})

	// handle error from transaction (500)
	if err != nil {
		println("ERROR IN USE CASE")
		return nil, err
	}

	groupM := appmapper.ToGroupResponse(
		group,
		user,
		enum.GroupRoleOwner,
		nil,
		1,
	)

	domain := utils.GetBaseURLFromIncomingContext(ctx)
	if domain == "" {
		domain = "https://www.schedulr.site" // fallback chuẩn hơn
	}

	link := fmt.Sprintf("%s/groups/%s", domain, group.ID)

	uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
		EventType:   appconstant.EventTypeGroupCreated,
		SenderID:    user.ID,
		ReceiverIDs: []string{user.ID},
		Payload: appdto.TeamNotificationMessagePayload{
			Title:           appconstant.GetDisplayTitle(appconstant.EventTypeGroupCreated),
			Message:         fmt.Sprintf("Bạn đã tạo nhóm %s thành công", group.Name),
			Link:            utils.Ptr(link),
			ImageURL:        nil,
			CorrelationID:   group.ID,
			CorrelationType: int(appconstant.CorrelationTypeGroup),
		},
		Metadata: appdto.TeamNotificationMessageMetadata{
			IsSentMail:           false,
			NonExistentReceivers: []string{},
		},
	}, &appdto.UserWithPermission{
		ID:                   user.ID,
		HasEmailNotification: user.HasEmailNotification,
		HasPushNotification:  user.HasPushNotification,
	})

	return &appdto.BaseResponse[appdto.GroupResponse]{
		Data:  groupM,
		Error: nil,
	}, nil
}

func (uc *groupUseCase) Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, errorbase.AppError) {
	// Implement the logic for the Ping method
	userID := utils.GetUserIDFromOutgoingContext(ctx)
	println("Ping received from user:", userID)
	return &common.EmptyResponse{}, nil
}

func MapRoleToString(role string) enum.GroupRole {
	switch role {
	case "owner":
		return enum.GroupRoleOwner
	case "manager":
		return enum.GroupRoleManager
	case "member":
		return enum.GroupRoleMember
	default:
		return enum.GroupRoleOwner
	}
}

func (uc *groupUseCase) ListGroups(ctx context.Context, req *appdto.ListGroupsRequest) (*appdto.BaseResponse[appdto.ListGroupsResponse], errorbase.AppError) {
	_ = req

	userID := utils.GetUserIDFromOutgoingContext(ctx)
	groups, err := uc.groupRepo.GetGroupsByUserID(ctx, userID)
	if err != nil {
		fmt.Println("Error fetching groups for user:", err)
		return &appdto.BaseResponse[appdto.ListGroupsResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	userId := utils.GetUserIDFromOutgoingContext(ctx)
	fmt.Printf("ListGroups called for userID: %s", userId)
	if groups == nil {
		groups = &appdto.ListGroupsResponse{
			Items: []appdto.ListGroupItem{},
			Total: 0,
		}
	}

	return &appdto.BaseResponse[appdto.ListGroupsResponse]{
		Data:  groups,
		Error: nil,
	}, nil
}

func (uc *groupUseCase) GetGroup(ctx context.Context, req *appdto.GetGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleViewer)
	if err != nil {
		return &appdto.BaseResponse[appdto.GroupResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	// validate request and existence
	err = uc.validator.ValidateGetListGroupMembers(ctx, &appdto.ListMembersRequest{GroupID: req.GroupID}, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.GroupResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	userID := actor.ID

	group, memmberCount, sprint, err := uc.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return &appdto.BaseResponse[appdto.GroupResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	if group == nil {
		return &appdto.BaseResponse[appdto.GroupResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    errdict.ErrNotFound.Code,
				Message: errdict.ErrNotFound.Title,
				Detail:  errdict.ErrNotFound.Detail,
			},
		}, nil
	}

	owner, err := uc.userRepo.GetUserByID(ctx, group.OwnerID)
	if err != nil {
		return &appdto.BaseResponse[appdto.GroupResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	myRole, err := uc.groupRepo.GetRoleByUserIDAndGroupID(ctx, userID, req.GroupID)
	if err != nil {
		return &appdto.BaseResponse[appdto.GroupResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	groupM := appmapper.ToGroupResponse(
		group,
		owner,
		MapRoleToString(myRole),
		&sprint,
		int(memmberCount),
	)

	return &appdto.BaseResponse[appdto.GroupResponse]{
		Data:  groupM,
		Error: nil,
	}, nil
}

func (uc *groupUseCase) UpdateGroup(ctx context.Context, req *appdto.UpdateGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleManager)
	if err != nil {
		return &appdto.BaseResponse[appdto.GroupResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	group, err := uc.validator.ValidateUpdateGroup(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.GroupResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var updatedGroup *entity.Group
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		updatedGroup, err = repo.GroupRepository().UpdateGroup(ctx, group)
		if err != nil {
			return err
		}

		if updatedGroup == nil {
			return errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group not found"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	owner, err := uc.userRepo.GetUserByID(ctx, updatedGroup.OwnerID)
	if err != nil {
		return &appdto.BaseResponse[appdto.GroupResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	myRole := actor.Role

	groupM := appmapper.ToGroupResponse(
		updatedGroup,
		owner,
		myRole,
		nil,
		0,
	)

	var usersID []string
	usersID, err = uc.groupRepo.GetListUserIDByGroupID(ctx, updatedGroup.ID)
	if err != nil {
		return nil, err
	}

	domain := utils.GetBaseURLFromIncomingContext(ctx)
	if domain == "" {
		domain = "https://www.schedulr.site"
	}
	link := fmt.Sprintf("%s/groups/%s", domain, group.ID)

	uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
		EventType:   appconstant.EventTypeGroupUpdated,
		SenderID:    actor.ID,
		ReceiverIDs: usersID,
		Payload: appdto.TeamNotificationMessagePayload{
			Title:           appconstant.GetDisplayTitle(appconstant.EventTypeGroupUpdated),
			Message:         fmt.Sprintf("Nhóm %s đã được cập nhật", updatedGroup.Name),
			Link:            utils.Ptr(link),
			ImageURL:        nil,
			CorrelationID:   updatedGroup.ID,
			CorrelationType: int(appconstant.CorrelationTypeGroup),
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

	return &appdto.BaseResponse[appdto.GroupResponse]{
		Data:  groupM,
		Error: nil,
	}, nil
}

func (uc *groupUseCase) DeleteGroup(ctx context.Context, req *appdto.DeleteGroupRequest) (*appdto.BaseResponse[appdto.DeleteGroupResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleOwner)
	if err != nil {
		return &appdto.BaseResponse[appdto.DeleteGroupResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	err = uc.validator.ValidateDeleteGroup(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.DeleteGroupResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var groupName string

	var deleted bool
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		group, _, _, err := uc.groupRepo.GetGroupByID(ctx, req.GroupID)
		groupName = group.Name

		sprints, err := repo.SprintRepository().GetSprintsByGroupID(ctx, req.GroupID)
		if err != nil {
			return err
		}
		for _, sprint := range sprints {
			if sprint.Status == enum.SprintStatusDraft {
				err = repo.SprintRepository().DeleteDraftSprint(ctx, sprint.ID)
				if err != nil {
					return err
				}
			}

			if sprint.Status == enum.SprintStatusActive {
				err = repo.SprintRepository().CancelSprint(ctx, sprint.ID)
				if err != nil {
					return err
				}
			}
		}

		err = repo.GroupRepository().DeleteGroup(ctx, req.GroupID)
		if err != nil {
			return err
		}

		deleted = true
		return nil
	})

	if err != nil {
		return nil, err
	}

	domain := utils.GetBaseURLFromIncomingContext(ctx)
	if domain == "" {
		domain = "https://www.schedulr.site"
	}
	link := fmt.Sprintf("%s/groups/", domain)

	var usersID []string
	usersID, err = uc.groupRepo.GetListUserIDByGroupID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
		EventType:   appconstant.EventTypeGroupDeleted,
		SenderID:    actor.ID,
		ReceiverIDs: usersID,
		Payload: appdto.TeamNotificationMessagePayload{
			Title:           appconstant.GetDisplayTitle(appconstant.EventTypeGroupDeleted),
			Message:         fmt.Sprintf("Nhóm %s đã bị xóa", groupName),
			Link:            utils.Ptr(link),
			ImageURL:        nil,
			CorrelationID:   req.GroupID,
			CorrelationType: int(appconstant.CorrelationTypeGroup),
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

	return &appdto.BaseResponse[appdto.DeleteGroupResponse]{
		Data: &appdto.DeleteGroupResponse{
			Success: deleted,
		},
		Error: nil,
	}, nil
}

func (uc *groupUseCase) GetListGroupMembers(ctx context.Context, req *appdto.ListMembersRequest) (*appdto.BaseResponse[appdto.ListMembersResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleViewer)
	if err != nil {
		return &appdto.BaseResponse[appdto.ListMembersResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	err = uc.validator.ValidateGetListGroupMembers(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.ListMembersResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	members, err := uc.store.UserRepository().GetListMembersByGroupID(ctx, req.GroupID)
	if err != nil {
		return &appdto.BaseResponse[appdto.ListMembersResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	return &appdto.BaseResponse[appdto.ListMembersResponse]{
		Data:  members,
		Error: nil,
	}, nil
}

func (uc *groupUseCase) GetSimpleUserByGroupID(ctx context.Context, req *appdto.ListMembersRequest) (*appdto.BaseResponse[[]appdto.SimpleUserResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleViewer)
	if err != nil {
		return &appdto.BaseResponse[[]appdto.SimpleUserResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	err = uc.validator.ValidateGetListGroupMembers(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[[]appdto.SimpleUserResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	users, err := uc.groupRepo.GetSimpleUsersByGroupID(ctx, req.GroupID)
	if err != nil {
		return &appdto.BaseResponse[[]appdto.SimpleUserResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var list []appdto.SimpleUserResponse
	if users != nil {
		list = make([]appdto.SimpleUserResponse, 0, len(users))
		for _, u := range users {
			if u == nil {
				continue
			}
			list = append(list, *u)
		}
	} else {
		list = make([]appdto.SimpleUserResponse, 0)
	}

	return &appdto.BaseResponse[[]appdto.SimpleUserResponse]{
		Data:  &list,
		Error: nil,
	}, nil
}

func (uc *groupUseCase) UpdateMemberRole(ctx context.Context, req *appdto.UpdateMemberRoleRequest) (*appdto.BaseResponse[appdto.MemberResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleManager)
	fmt.Println("actor in UpdateMemberRole:", actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.MemberResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var member *entity.User
	member, err = uc.validator.ValidateUpdateMemberRole(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.MemberResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var updatedMember *appdto.MemberResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		if req.Role == enum.GroupRoleOwner && actor.Role == enum.GroupRoleOwner {
			updatedMember, err = repo.GroupRepository().UpdateMemberRole(ctx, req.MemberId, req.GroupID, string(req.Role))
			if err != nil {
				return err
			}

			_, err := repo.GroupRepository().UpdateMemberRole(ctx, actor.ID, req.GroupID, string(enum.GroupRoleMember))
			if err != nil {
				return err
			}

			return nil
		}
		updatedMember, err = repo.GroupRepository().UpdateMemberRole(ctx, req.MemberId, req.GroupID, string(req.Role))
		if err != nil {
			return err
		}

		if updatedMember == nil {
			return errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("member not found"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	domain := utils.GetBaseURLFromIncomingContext(ctx)
	if domain == "" {
		domain = "https://www.schedulr.site"
	}
	link := fmt.Sprintf("%s/groups/%s", domain, req.GroupID)

	var usersID []string
	usersID, err = uc.groupRepo.GetListUserIDByGroupID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
		EventType:   appconstant.EventTypeMemberRoleUpdated,
		SenderID:    actor.ID,
		ReceiverIDs: usersID,
		Payload: appdto.TeamNotificationMessagePayload{
			Title:           appconstant.GetDisplayTitle(appconstant.EventTypeMemberRoleUpdated),
			Message:         fmt.Sprintf("Vai trò của %s trong nhóm  đã được cập nhật thành %s", member.Email, req.Role.String()),
			Link:            utils.Ptr(link),
			ImageURL:        nil,
			CorrelationID:   req.GroupID,
			CorrelationType: int(appconstant.CorrelationTypeGroup),
		},
		Metadata: appdto.TeamNotificationMessageMetadata{
			IsSentMail:           false,
			NonExistentReceivers: []string{},
		},
	}, &appdto.UserWithPermission{
		ID:                   req.MemberId,
		HasEmailNotification: actor.HasEmailNotification,
		HasPushNotification:  actor.HasPushNotification,
	})

	return &appdto.BaseResponse[appdto.MemberResponse]{
		Data:  updatedMember,
		Error: nil,
	}, nil
}

func (uc *groupUseCase) RemoveMember(ctx context.Context, req *appdto.RemoveMemberRequest) (*appdto.BaseResponse[appdto.RemoveMemberResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleManager)
	if err != nil {
		return &appdto.BaseResponse[appdto.RemoveMemberResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}
	var member *entity.User
	member, err = uc.validator.ValidateRemoveMember(ctx, req, actor)
	if err != nil {
		return &appdto.BaseResponse[appdto.RemoveMemberResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var deleted *appdto.RemoveMemberResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		err = repo.WorkRepository().UnassignWorksByMember(ctx, req.GroupID, req.MemberId)
		if err != nil {
			return err
		}

		err = repo.GroupRepository().RemoveMember(ctx, req.GroupID, req.MemberId)
		if err != nil {
			return err
		}

		deleted = &appdto.RemoveMemberResponse{
			Success: true,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	domain := utils.GetBaseURLFromIncomingContext(ctx)
	if domain == "" {
		domain = "https://www.schedulr.site"
	}
	link := fmt.Sprintf("%s/groups/%s", domain, req.GroupID)

	var usersID []string
	usersID, err = uc.groupRepo.GetListUserIDByGroupID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
		EventType:   appconstant.EventTypeMemberRemoved,
		SenderID:    actor.ID,
		ReceiverIDs: usersID,
		Payload: appdto.TeamNotificationMessagePayload{
			Title:           appconstant.GetDisplayTitle(appconstant.EventTypeMemberRemoved),
			Message:         fmt.Sprintf("Thành viên %s đã rời khỏi nhóm ", member.Email),
			Link:            utils.Ptr(link),
			ImageURL:        nil,
			CorrelationID:   req.GroupID,
			CorrelationType: int(appconstant.CorrelationTypeGroup),
		},
		Metadata: appdto.TeamNotificationMessageMetadata{
			IsSentMail:           false,
			NonExistentReceivers: []string{},
		},
	}, &appdto.UserWithPermission{
		ID:                   req.MemberId,
		HasEmailNotification: actor.HasEmailNotification,
		HasPushNotification:  actor.HasPushNotification,
	})

	return &appdto.BaseResponse[appdto.RemoveMemberResponse]{
		Data:  deleted,
		Error: nil,
	}, nil
}

func (uc *groupUseCase) CreateInvite(ctx context.Context, req *appdto.CreateInviteRequest) (*appdto.BaseResponse[appdto.InviteResponse], errorbase.AppError) {
	actor, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleManager)
	if err != nil {
		return &appdto.BaseResponse[appdto.InviteResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	invite, err := uc.validator.ValidateCreateInvitation(ctx, req, actor)
	fmt.Println("invite after validation:", invite)
	if err != nil {
		return &appdto.BaseResponse[appdto.InviteResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	var createdInvite *entity.Invite
	createdInvite, err = uc.store.InviteRepository().CreateInvite(ctx, invite)
	if err != nil {
		return nil, err
	}

	var usersID []string
	if createdInvite.Email != nil {
		user, err := uc.userRepo.GetUserByEmail(ctx, *createdInvite.Email)
		if err != nil {
			return nil, err
		}
		if user != nil {
			usersID = append(usersID, user.ID)
		}
	}

	if createdInvite == nil {
		return &appdto.BaseResponse[appdto.InviteResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    errdict.ErrInternal.Code,
				Message: errdict.ErrInternal.Title,
				Detail:  errdict.ErrInternal.Detail,
			},
		}, nil
	}
	var inviteLink string

	inviteLink = fmt.Sprintf(
		"https://www.schedulr.site/api/v1/ts/invitation/acceptance?code=%s",
		createdInvite.Token,
	)

	uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
		EventType:   appconstant.EventTypeInviteCreated,
		SenderID:    invite.CreatedBy,
		ReceiverIDs: usersID,
		Payload: appdto.TeamNotificationMessagePayload{
			Title:           appconstant.GetDisplayTitle(appconstant.EventTypeInviteCreated),
			Message:         fmt.Sprintf("Bạn đã tạo một lời mời tham gia nhóm thành công. Mã lời mời sẽ hết hạn vào %s", invite.ExpiresAt.Format("02/01/2006 15:04:05")),
			Link:            utils.Ptr(inviteLink),
			ImageURL:        nil,
			CorrelationID:   invite.GroupID,
			CorrelationType: int(appconstant.CorrelationTypeGroup),
		},
		Metadata: appdto.TeamNotificationMessageMetadata{
			IsSentMail:           false,
			NonExistentReceivers: []string{},
		},
	}, &appdto.UserWithPermission{
		ID:                   actor.ID,
		HasEmailNotification: true,
		HasPushNotification:  true,
	})

	return &appdto.BaseResponse[appdto.InviteResponse]{
		Data: &appdto.InviteResponse{
			Code:      createdInvite.Token,
			ExpiresAt: createdInvite.ExpiresAt,
			CreatedAt: createdInvite.CreatedAt,
		},
		Error: nil,
	}, nil

}

func (uc *groupUseCase) AcceptInvite(ctx context.Context, req *appdto.AcceptInviteRequest) (*appdto.BaseResponse[appdto.AcceptInviteResponse], errorbase.AppError) {
	domain := utils.GetBaseURLFromIncomingContext(ctx)
	if domain == "" {
		domain = "https://www.schedulr.site"
	}
	link := fmt.Sprintf("%s/groups/", domain)

	invite, user, err := uc.validator.ValidateAcceptInvitation(ctx, req)
	if err != nil {
		if err.ErrorInfo().Detail != nil && (*err.ErrorInfo().Detail == "invite is expired" || *err.ErrorInfo().Detail == "the group has reached maximum number of members") {

			uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
				EventType:   appconstant.EventTypeInviteError,
				SenderID:    invite.CreatedBy,
				ReceiverIDs: []string{user.ID},
				Payload: appdto.TeamNotificationMessagePayload{
					Title:           appconstant.GetDisplayTitle(appconstant.EventTypeInviteError),
					Message:         fmt.Sprintf("Thành viên %s đã chấp nhận lời mời tham gia nhóm", user.Email),
					Link:            utils.Ptr(link),
					ImageURL:        nil,
					CorrelationID:   invite.GroupID,
					CorrelationType: int(appconstant.CorrelationTypeGroup),
				},
				Metadata: appdto.TeamNotificationMessageMetadata{
					IsSentMail:           false,
					NonExistentReceivers: []string{},
				},
			}, &appdto.UserWithPermission{
				ID:                   invite.CreatedBy,
				HasEmailNotification: true,
				HasPushNotification:  true,
			})

		}
		return &appdto.BaseResponse[appdto.AcceptInviteResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	// create group member
	groupMember, err := entity.NewGroupMember(
		uuid.NewString(),
		invite.GroupID,
		user.ID,
		enum.GroupRole(invite.Role),
		time.Now(),
	)
	if err != nil {
		return &appdto.BaseResponse[appdto.AcceptInviteResponse]{
			Data:  nil,
			Error: appmapper.ToErrorResponse(err),
		}, nil
	}

	var created *appdto.AcceptInviteResponse
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		err = repo.GroupRepository().AddGroupMember(ctx, groupMember)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	link = fmt.Sprintf("%s/groups/%s", domain, invite.GroupID)

	created = &appdto.AcceptInviteResponse{
		Location: link,
	}

	var usersID []string
	usersID, err = uc.groupRepo.GetListUserIDByGroupID(ctx, invite.GroupID)
	if err != nil {
		return nil, err
	}

	uc.notificationHelper.PublishTeamNotificationMessage(ctx, appdto.TeamNotificationMessage{
		EventType:   appconstant.EventTypeInviteAccepted,
		SenderID:    invite.CreatedBy,
		ReceiverIDs: usersID,
		Payload: appdto.TeamNotificationMessagePayload{
			Title:           appconstant.GetDisplayTitle(appconstant.EventTypeInviteAccepted),
			Message:         fmt.Sprintf("Thành viên %s đã chấp nhận lời mời tham gia nhóm", user.Email),
			Link:            utils.Ptr(link),
			ImageURL:        nil,
			CorrelationID:   invite.GroupID,
			CorrelationType: int(appconstant.CorrelationTypeGroup),
		},
		Metadata: appdto.TeamNotificationMessageMetadata{
			IsSentMail:           false,
			NonExistentReceivers: []string{},
		},
	}, &appdto.UserWithPermission{
		ID:                   invite.CreatedBy,
		HasEmailNotification: true,
		HasPushNotification:  true,
	})

	return &appdto.BaseResponse[appdto.AcceptInviteResponse]{
		Data:  created,
		Error: nil,
	}, nil
}
