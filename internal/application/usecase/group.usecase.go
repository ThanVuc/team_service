package usecase

import (
	"context"
	appdto "team_service/internal/application/common/dto"
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
	store     istore.Store
	groupRepo irepository.GroupRepository
	userRepo  irepository.UserRepository
	validator *appvalidation.GroupValidator
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

func (uc *groupUseCase) GetGroupRequest(ctx context.Context, req *appdto.GetGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError) {
	userID := utils.GetUserIDFromOutgoingContext(ctx)

	group, memmberCount, sprint, err := uc.groupRepo.GetGroupByID(ctx, req.GroupID)
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
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	myRole, err := uc.groupRepo.GetRoleByUserIDAndGroupID(ctx, userID, req.GroupID)
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
	userID := utils.GetUserIDFromOutgoingContext(ctx)

	group, err := uc.validator.ValidateUpdateGroup(ctx, req)
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
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	myRole, err := uc.groupRepo.GetRoleByUserIDAndGroupID(ctx, userID, req.GroupID)
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

	groupM := appmapper.ToGroupResponse(
		updatedGroup,
		owner,
		MapRoleToString(myRole),
		nil,
		0,
	)

	return &appdto.BaseResponse[appdto.GroupResponse]{
		Data:  groupM,
		Error: nil,
	}, nil

}

func (uc *groupUseCase) DeleteGroup(ctx context.Context, req *appdto.DeleteGroupRequest) (*appdto.BaseResponse[appdto.DeleteGroupResponse], errorbase.AppError) {
	err := uc.validator.ValidateDeleteGroup(ctx, req)

	if err != nil {
		return &appdto.BaseResponse[appdto.DeleteGroupResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {

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

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &appdto.BaseResponse[appdto.DeleteGroupResponse]{
		Data: &appdto.DeleteGroupResponse{
			Success: true,
		},
		Error: nil,
	}, nil
}
