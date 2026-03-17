package appvalidation

import (
	"context"
	appdto "team_service/internal/application/common/dto"
	irepository "team_service/internal/application/common/interface/repository"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"team_service/internal/infrastructure/share/utils"
	"time"

	"github.com/google/uuid"
)

type GroupValidator struct {
	groupRepo irepository.GroupRepository
	userRepo  irepository.UserRepository
}

func NewGroupValidator(
	groupRepo irepository.GroupRepository,
	userRepo irepository.UserRepository,
) *GroupValidator {
	return &GroupValidator{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

func (v *GroupValidator) ValidateCreateGroup(ctx context.Context, req *appdto.CreateGroupRequest) (*entity.Group, *entity.User, errorbase.AppError) {
	userID := utils.GetUserIDFromOutgoingContext(ctx)
	group, err := entity.NewGroup(
		uuid.NewString(),
		req.Name,
		userID,
		req.Description,
		time.Now(),
	)

	if err != nil {
		return nil, nil, err
	}

	count, err := v.groupRepo.CountGroupsByOwner(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	if count >= 10 {
		return nil, nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("user has reached the maximum number of groups"))
	}

	user, err := v.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	if group == nil {
		return nil, nil, errorbase.New(errdict.ErrInternal, errorbase.WithDetail("group is nil"))
	}

	if user == nil {
		return nil, nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("user not found"))
	}

	return group, user, nil
}

func (v *GroupValidator) ValidateUpdateGroup(ctx context.Context, req *appdto.UpdateGroupRequest) (*entity.Group, errorbase.AppError) {
	userID := utils.GetUserIDFromOutgoingContext(ctx)
	if req.GroupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	role, err := v.groupRepo.GetRoleByUserIDAndGroupID(ctx, userID, req.GroupID)
	if err != nil {
		return nil, err
	}

	if role != string(enum.GroupRoleOwner) && role != string(enum.GroupRoleManager) {
		return nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("user does not have permission to update the group"))
	}

	groupExists, err := v.groupRepo.CheckGroupExists(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	if !groupExists {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group not found"))
	}

	name := ""
	if req.Name != nil {
		name = *req.Name
	}

	group, err := entity.NewGroup(
		req.GroupID,
		name,
		userID,
		req.Description,
		time.Now(),
	)

	if err != nil {
		return nil, err
	}

	err = group.Update(req.Name, req.Description, time.Now())

	if err != nil {
		return nil, err
	}

	return group, nil
}

func (v *GroupValidator) ValidateDeleteGroup(ctx context.Context, req *appdto.DeleteGroupRequest) errorbase.AppError {
	userID := utils.GetUserIDFromOutgoingContext(ctx)

	if req.GroupID == "" {
		return errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	role, err := v.groupRepo.GetRoleByUserIDAndGroupID(ctx, userID, req.GroupID)
	if err != nil {
		return err
	}

	if role != string(enum.GroupRoleOwner) {
		return errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("user does not have permission to delete the group"))
	}

	groupExists, err := v.groupRepo.CheckGroupExists(ctx, req.GroupID)
	if err != nil {
		return err
	}

	if !groupExists {
		return errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group not found"))
	}

	return nil
}
