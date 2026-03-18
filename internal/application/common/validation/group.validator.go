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

func (v *GroupValidator) ValidateGetListGroupMembers(ctx context.Context, req *appdto.ListMembersRequest) errorbase.AppError {
	userID := utils.GetUserIDFromOutgoingContext(ctx)
	if userID == "" {
		return errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("missing user context"))
	}

	if req.GroupID == "" {
		return errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	check, err := v.groupRepo.CheckGroupExists(ctx, req.GroupID)
	if err != nil {
		return err
	}

	if !check {
		return errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group is not found"))
	}

	role, err := v.groupRepo.GetRoleByUserIDAndGroupID(ctx, userID, req.GroupID)
	if err != nil {
		return err
	}

	if role == "" {
		return errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("user is not a member of the group"))
	}

	return nil
}

func (v *GroupValidator) ValidateUpdateMemberRole(ctx context.Context, req *appdto.UpdateMemberRoleRequest) errorbase.AppError {
	userID := utils.GetUserIDFromOutgoingContext(ctx)
	if userID == "" {
		return errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("missing user context"))
	}

	if req.GroupID == "" {
		return errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	checkGroupExist, err := v.groupRepo.CheckGroupExists(ctx, req.GroupID)
	if err != nil {
		return err
	}

	if !checkGroupExist {
		return errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group is not found"))
	}

	myRole, err := v.groupRepo.GetRoleByUserIDAndGroupID(ctx, userID, req.GroupID)
	if err != nil {
		return err
	}

	if myRole == "" {
		return errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("user is not a member of the group"))
	}

	if myRole != string(enum.GroupRoleOwner) && myRole != string(enum.GroupRoleManager) {
		return errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("user does not have permission to update member role"))
	}

	if myRole == string(enum.GroupRoleManager) && req.Role == enum.GroupRoleOwner {
		return errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("manager cannot assign owner role"))
	}

	memberRole, err := v.groupRepo.GetRoleByUserIDAndGroupID(ctx, req.MemberId, req.GroupID)

	if err != nil {
		return err
	}

	if memberRole == "" {
		return errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("member is not found in the group"))
	}

	if memberRole == string(enum.GroupRoleOwner) {
		return errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("cannot update owner role"))
	}

	if userID == req.MemberId {
		return errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("cannot update own role"))
	}

	count, err := v.groupRepo.CountManagerAndMemberByGroupID(ctx, req.GroupID)
	if err != nil {
		return err
	}

	if memberRole == string(enum.GroupRoleViewer) {
		if count >= 9 {
			return errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("group has reached the maximum number of members"))
		}
	}

	return nil

}
