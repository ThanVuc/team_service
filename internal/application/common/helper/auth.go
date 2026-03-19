package apphelper

import (
	"context"
	"fmt"
	appconstant "team_service/internal/application/common/constant"
	appdto "team_service/internal/application/common/dto"
	icacherepository "team_service/internal/application/common/interface/cacherepository"
	irepository "team_service/internal/application/common/interface/repository"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/enum"
	"team_service/internal/infrastructure/share/utils"

	"github.com/thanvuc/go-core-lib/log"
)

type AuthHelper struct {
	userRepo  irepository.UserRepository
	cacheRepo icacherepository.CacheRepository
	logger    log.LoggerV2
}

func NewAuthHelper(
	userRepo irepository.UserRepository,
	cacheRepo icacherepository.CacheRepository,
	logger log.LoggerV2,
) *AuthHelper {
	return &AuthHelper{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
		logger:    logger,
	}
}

func (h *AuthHelper) RequireRole(ctx context.Context, expectedRole enum.GroupRole) (*appdto.UserWithPermission, errorbase.AppError) {
	userID := utils.GetUserIDFromOutgoingContext(ctx)
	groupId := utils.GetGroupIDFromContext(ctx)
	key := appconstant.CacheUserWithRolePrefix + userID

	if userID == "" {
		return nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("user id is required in context"))
	}

	if groupId == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required in context"))
	}

	if err := h.cacheRepo.Get(ctx, key, &expectedRole); err == nil {
		if !expectedRole.HasPermission(expectedRole) {
			return nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("User does not have permission!"))
		}

		return &appdto.UserWithPermission{
			ID:   userID,
			Role: expectedRole,
		}, nil
	}

	user, err := h.userRepo.GetUserWithPermissionByID(ctx, groupId, userID)
	if err != nil {
		return nil, err
	}

	if cacheErr := h.cacheRepo.Set(ctx, key, user.Role.String(), 800); cacheErr != nil {
		h.logger.Error(fmt.Sprintf("Failed to set user role in cache for user %s: %v", userID, cacheErr))
	}

	if !user.Role.HasPermission(expectedRole) {
		return nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("User does not have permission!"))
	}

	return user, nil
}
