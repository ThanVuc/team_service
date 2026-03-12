package usecase

import (
	"context"
	"fmt"
	"log"
	irepository "team_service/internal/application/common/interface/repository"
	istore "team_service/internal/application/common/interface/store"
	appmapper "team_service/internal/application/common/mapper"
	appvalidation "team_service/internal/application/common/validation"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/internal/infrastructure/share/utils"
	"team_service/proto/common"
	"team_service/proto/team_service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type groupUseCase struct {
	store      istore.Store
	mapper     *appmapper.GroupMapper
	groupRepos irepository.GroupRepository
}

func (uc *groupUseCase) CreateGroup(ctx context.Context, req *team_service.CreateGroupRequest) (*team_service.CreateGroupResponse, errorbase.AppError) {
	if uc.store == nil {
		panic("store is nil")
	}

	if uc.mapper == nil {
		panic("mapper is nil")
	}

	if err := appvalidation.ValidateGroup(ctx, req); err != nil {
		return &team_service.CreateGroupResponse{
			Group: nil,
			Error: &team_service.Error{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Details: err.ErrorInfo().Detail,
			},
		}, nil
	}

	userID := utils.GetUserIDFromOutgoingContext(ctx)

	var group *database.Group
	var user *database.GetUserByIDRow

	err := uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		count, err := repo.GroupRepository().CountGroupsByOwner(ctx, userID)
		if err != nil {
			return err
		}

		if count >= 10 {
			return errorbase.New(errdict.ErrNotFound)
		}

		group, err = repo.GroupRepository().CreateGroup(ctx, req, userID)
		if err != nil {
			return err
		}

		user, err = repo.GroupRepository().GetUserByID(ctx, userID)
		if err != nil {
			return err
		}

		gm := pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		}

		err = repo.GroupRepository().AddGroupMember(ctx, database.CreateGroupMemberParams{
			ID:      gm,
			GroupID: group.ID,
			UserID:  user.ID,
			Role:    MapRole(team_service.GroupRole_GROUP_ROLE_OWNER),
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return &team_service.CreateGroupResponse{
			Group: nil,
			Error: &team_service.Error{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Details: err.ErrorInfo().Detail,
			},
		}, nil
	}

	groupM := uc.mapper.MapGroupMessage(group, user)

	return &team_service.CreateGroupResponse{
		Group: groupM,
	}, nil
}

func MapRole(role team_service.GroupRole) string {
	switch role {
	case team_service.GroupRole_GROUP_ROLE_OWNER:
		return "owner"
	case team_service.GroupRole_GROUP_ROLE_MANAGER:
		return "admin"
	case team_service.GroupRole_GROUP_ROLE_MEMBER:
		return "member"
	default:
		return "member"
	}
}

func (uc *groupUseCase) Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, errorbase.AppError) {
	// Implement the logic for the Ping method
	userID := utils.GetUserIDFromOutgoingContext(ctx)
	println("Received Ping request from user ID:", userID)
	return &common.EmptyResponse{}, nil
}

func MapRoleToString(role string) team_service.GroupRole {
	switch role {
	case "owner":
		return team_service.GroupRole_GROUP_ROLE_OWNER
	case "admin":
		return team_service.GroupRole_GROUP_ROLE_MANAGER
	case "member":
		return team_service.GroupRole_GROUP_ROLE_MEMBER
	default:
		return team_service.GroupRole_GROUP_ROLE_MEMBER
	}
}

func (uc *groupUseCase) GetGroup(ctx context.Context, req *common.IDRequest) (*team_service.GetGroupResponse, errorbase.AppError) {
	groupID, errS := uuid.Parse(req.Id)
	if errS != nil {
		return &team_service.GetGroupResponse{
			Group: nil,
			Error: &team_service.Error{
				Code:    errdict.ErrInvalidUUID.Code,
				Message: errdict.ErrInvalidUUID.Title,
				Details: errdict.ErrInvalidUUID.Detail,
			},
		}, nil
	}

	fmt.Printf("Fetching group with ID: %s\n", groupID.String())

	userID := utils.GetUserIDFromOutgoingContext(ctx)

	group, memberCount, sprint, myRole, err := repo.GroupRepository().GetGroupByID(ctx, userID, groupID.String())
	if err != nil {
		log.Printf("Error fetching group: %v\n", err)
		return errorbase.New(errdict.ErrBadRequest)
	}
	owner, err = repo.GroupRepository().GetUserByID(ctx, group.OwnerID.String())
	if err != nil {
		log.Printf("Error fetching group owner: %v\n", err)
		return errorbase.New(errdict.ErrBadRequest)
	}

	myRoleEnum := MapRoleToString(myRole)

	groupDetail := uc.mapper.MapGroupDetail(group, owner, memberCount, myRoleEnum, sprint)

	return nil

	if err != nil {
		return &team_service.GetGroupResponse{
			Group: nil,
			Error: &team_service.Error{
				Code:    err.Error(),
				Message: "Group not found",
			},
		}, nil
	}

	return &team_service.GetGroupResponse{
		Group: groupDetail,
		Error: nil,
	}, nil
}
