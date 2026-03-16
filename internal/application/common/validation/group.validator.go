package appvalidation

import (
	"context"
	appdto "team_service/internal/application/common/dto"
	irepository "team_service/internal/application/common/interface/repository"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
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
