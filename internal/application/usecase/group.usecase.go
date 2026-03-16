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
	"team_service/internal/infrastructure/share/utils"
	"team_service/proto/common"
	"time"

	"github.com/google/uuid"
)

type groupUseCase struct {
	store      istore.Store
	groupRepo  irepository.GroupRepository
	validator  *appvalidation.GroupValidator
	authHelper *apphelper.AuthHelper
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
		println("repo:", repo)
		println("group repo:", repo.GroupRepository())
		println("user repo:", repo.UserRepository())
		println("group:", group)
		println("user:", user)
		group, err = repo.GroupRepository().CreateGroup(ctx, group, user.ID)
		println("PASS 1")
		if err != nil {
			return err
		}

		println("PASS 2")
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
		return nil, err
	}

	groupM := appmapper.ToGroupResponse(
		group,
		user,
		enum.GroupRoleOwner,
		nil,
		1,
	)
	println("CreatedAt: " + groupM.CreatedAt.String())
	println("UpdatedAt: " + groupM.UpdatedAt.String())
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

// func MapRoleToString(role string) team_service.GroupRole {
// 	switch role {
// 	case "owner":
// 		return team_service.GroupRole_GROUP_ROLE_OWNER
// 	case "admin":
// 		return team_service.GroupRole_GROUP_ROLE_MANAGER
// 	case "member":
// 		return team_service.GroupRole_GROUP_ROLE_MEMBER
// 	default:
// 		return team_service.GroupRole_GROUP_ROLE_MEMBER
// 	}
// }

// func (uc *groupUseCase) GetGroup(ctx context.Context, req *common.IDRequest) (*team_service.GetGroupResponse, errorbase.AppError) {
// 	groupID, errS := uuid.Parse(req.Id)
// 	if errS != nil {
// 		return &team_service.GetGroupResponse{
// 			Group: nil,
// 			Error: &team_service.Error{
// 				Code:    errdict.ErrInvalidUUID.Code,
// 				Message: errdict.ErrInvalidUUID.Title,
// 				Details: errdict.ErrInvalidUUID.Detail,
// 			},
// 		}, nil
// 	}

// 	fmt.Printf("Fetching group with ID: %s\n", groupID.String())

// 	userID := utils.GetUserIDFromOutgoingContext(ctx)

// 	group, memberCount, sprint, myRole, err := uc.groupRepo.GetGroupByID(ctx, userID, groupID.String())
// 	if err != nil {
// 		log.Printf("Error fetching group: %v\n", err)
// 		return &team_service.GetGroupResponse{
// 			Group: nil,
// 			Error: &team_service.Error{
// 				Code:    err.ErrorInfo().Code,
// 				Message: err.ErrorInfo().Title,
// 				Details: err.ErrorInfo().Detail,
// 			},
// 		}, nil
// 	}

// 	owner, err := uc.groupRepo.GetUserByID(ctx, group.OwnerID.String())
// 	if err != nil {
// 		log.Printf("Error fetching group owner: %v\n", err)
// 		return &team_service.GetGroupResponse{
// 			Group: nil,
// 			Error: &team_service.Error{
// 				Code:    err.ErrorInfo().Code,
// 				Message: err.ErrorInfo().Title,
// 				Details: err.ErrorInfo().Detail,
// 			},
// 		}, nil
// 	}

// 	myRoleEnum := MapRoleToString(myRole)

// 	groupDetail := appmapper.MapGroupDetail(group, owner, memberCount, myRoleEnum, sprint)

// 	return &team_service.GetGroupResponse{
// 		Group: groupDetail,
// 		Error: nil,
// 	}, nil
// }

// Auth helper Example with delete
func (uc *groupUseCase) DeleteGroup(ctx context.Context, req *common.IDRequest) (*common.EmptyResponse, errorbase.AppError) {
	// Example: Check if the user has the required role to delete the group with req.id = groupId
	uc.authHelper.RequireRole(ctx, req.Id, enum.GroupRoleOwner)

	// Implement the logic for deleting a group
	return &common.EmptyResponse{}, nil
}
