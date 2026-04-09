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
	groupRepo  irepository.GroupRepository
	userRepo   irepository.UserRepository
	inviteRepo irepository.InviteRepository
}

func NewGroupValidator(
	groupRepo irepository.GroupRepository,
	userRepo irepository.UserRepository,
	inviteRepo irepository.InviteRepository,
) *GroupValidator {
	return &GroupValidator{
		groupRepo:  groupRepo,
		userRepo:   userRepo,
		inviteRepo: inviteRepo,
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

// Note: validator signatures below now accept actor to avoid fetching roles from DB inside validators.
func (v *GroupValidator) ValidateUpdateGroup(ctx context.Context, req *appdto.UpdateGroupRequest, actor *appdto.UserWithPermission) (*entity.Group, errorbase.AppError) {
	userID := ""
	if actor != nil {
		userID = actor.ID
	}

	if req.GroupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	groupExists, err := v.groupRepo.CheckGroupExists(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	if !groupExists {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group not found"))
	}

	if !(actor.Role == enum.GroupRoleOwner || actor.Role == enum.GroupRoleManager) {
		forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("user does not have permission to update the group"))
		return nil, forbidden
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

func (v *GroupValidator) ValidateDeleteGroup(ctx context.Context, req *appdto.DeleteGroupRequest, actor *appdto.UserWithPermission) errorbase.AppError {
	if req.GroupID == "" {
		return errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
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

func (v *GroupValidator) ValidateGetListGroupMembers(ctx context.Context, req *appdto.ListMembersRequest, actor *appdto.UserWithPermission) errorbase.AppError {
	// Use actor to determine requester identity and group association.
	if actor == nil || actor.ID == "" {
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

	// Ensure the actor is operating within the same group
	if actor.GroupId != req.GroupID {
		return errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("user is not a member of the group"))
	}

	return nil
}

func (v *GroupValidator) ValidateUpdateMemberRole(ctx context.Context, req *appdto.UpdateMemberRoleRequest, actor *appdto.UserWithPermission) (*entity.User, errorbase.AppError) {
	// Validate request fields and group existence. Role checks and permission enforcement moved to usecase.
	if actor == nil || actor.ID == "" {
		return nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("missing user context"))
	}

	if req.GroupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	checkGroupExist, err := v.groupRepo.CheckGroupExists(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	if !checkGroupExist {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group is not found"))
	}

	if req.MemberId == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("member id is required"))
	}

	if actor.Role != enum.GroupRoleOwner && actor.Role != enum.GroupRoleManager {
		forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("user does not have permission to update member role"))
		return nil, forbidden
	}

	if actor.ID == req.MemberId {
		forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("cannot update own role"))
		return nil, forbidden
	}

	memberRole, err := v.groupRepo.GetRoleByUserIDAndGroupID(ctx, req.MemberId, req.GroupID)
	if err != nil {
		return nil, err
	}

	if memberRole == "" {
		notFound := errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("member is not found in the group"))
		return nil, notFound
	}

	if memberRole == string(enum.GroupRoleOwner) {
		forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("cannot update owner role"))
		return nil, forbidden
	}

	// Manager cannot assign owner role
	if actor.Role == enum.GroupRoleManager && req.Role == enum.GroupRoleOwner {
		forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("manager cannot assign owner role"))
		return nil, forbidden
	}

	// count constraint when promoting from viewer to manager/member
	count, err := v.groupRepo.CountManagerAndMemberByGroupID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	if memberRole == string(enum.GroupRoleViewer) {
		if count >= 9 {
			forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("the group has reached maximum number of managers and members"))
			return nil, forbidden
		}
	}

	member, err := v.userRepo.GetUserByID(ctx, req.MemberId)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (v *GroupValidator) ValidateRemoveMember(ctx context.Context, req *appdto.RemoveMemberRequest, actor *appdto.UserWithPermission) (*entity.User, errorbase.AppError) {
	if actor == nil || actor.ID == "" {
		return nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("missing user context"))
	}

	if req.GroupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	checkGroupExist, err := v.groupRepo.CheckGroupExists(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	if !checkGroupExist {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group is not found"))
	}

	if req.MemberId == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("member id is required"))
	}

	if actor.ID == req.MemberId {
		forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("cannot remove own membership"))
		return nil, forbidden
	}

	memberRole, err := v.groupRepo.GetRoleByUserIDAndGroupID(ctx, req.MemberId, req.GroupID)
	if err != nil {
		return nil, err
	}

	if memberRole == "" {
		notFound := errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("member is not found in the group"))
		return nil, notFound
	}

	if memberRole == string(enum.GroupRoleOwner) {
		forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("cannot remove owner"))
		return nil, forbidden
	}
	member, err := v.userRepo.GetUserByID(ctx, req.MemberId)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (v *GroupValidator) ValidateCreateInvitation(ctx context.Context, req *appdto.CreateInviteRequest, actor *appdto.UserWithPermission) (*entity.Invite, errorbase.AppError) {
	if actor == nil || actor.ID == "" {
		return nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("missing user context"))
	}

	if req.GroupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	checkGroupExist, err := v.groupRepo.CheckGroupExists(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	if !checkGroupExist {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group is not found"))
	}

	countMember, err := v.groupRepo.CountManagerAndMemberByGroupID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	if countMember >= 9 {
		return nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("the group has reached maximum number of managers and members"))
	}

	if req.Email != nil {
		if actor.Email == *req.Email {
			forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("cannot invite oneself"))
			return nil, forbidden
		}
	}

	// enforce business rules: Manager can only invite Member or Viewer; only Owner can invite Manager
	if actor.Role == enum.GroupRoleManager {
		if req.Role.String() == string(enum.GroupRoleOwner) || req.Role.String() == string(enum.GroupRoleManager) {
			forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("manager cannot assign owner or manager role"))
			return nil, forbidden
		}
	}

	if req.Role.String() == string(enum.GroupRoleOwner) {
		forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("cannot assign owner role"))
		return nil, forbidden
	}
	var email *string
	if req.Email != nil {
		email = req.Email
	}

	// Business permission checks (who can invite which role) moved to usecase.
	newInvite, err := entity.NewInvite(
		uuid.NewString(),
		req.GroupID,
		uuid.NewString(),
		req.Role,
		email,
		time.Now().Add(7*24*time.Hour),
		actor.ID,
		time.Now(),
	)

	return newInvite, err
}

func (v *GroupValidator) ValidateAcceptInvitation(ctx context.Context, req *appdto.AcceptInviteRequest) (*entity.Invite, *entity.User, errorbase.AppError) {
	userID := utils.GetUserIDFromOutgoingContext(ctx)
	if userID == "" {
		return nil, nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("missing user context"))
	}

	if req.Code == "" {
		return nil, nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("code is required"))
	}

	invite, err := v.inviteRepo.GetInviteByToken(ctx, req.Code)
	if err != nil {
		return nil, nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("invite not found"))
	}

	if invite == nil {
		return nil, nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("invite not found"))
	}

	var user *entity.User
	if invite.Email != nil {
		user, err = v.userRepo.GetUserByEmail(ctx, *invite.Email)
		if err != nil {
			return nil, nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("user not found with the email associated with the invite"))
		}

		if invite.Email != nil && user.Email != *invite.Email {
			return nil, nil, errorbase.New(errdict.ErrEmailNotMatched, errorbase.WithDetail("the invite is not associated with the user's email"))
		}
	} else {
		user, err = v.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			return nil, nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("user not found"))
		}
	}

	if user == nil {
		return nil, nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("user not found"))
	}

	if time.Now().UTC().After(invite.ExpiresAt) {
		return nil, nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("invite is expired"))
	}

	existingRole, err := v.groupRepo.GetRoleByUserIDAndGroupID(ctx, user.ID, invite.GroupID)
	if err != nil {
		return nil, nil, errorbase.New(errdict.ErrInternal, errorbase.WithDetail("failed to check existing membership"))
	}

	if existingRole != "" {
		forbidden := errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("user is already a member of the group"))
		return nil, nil, forbidden
	}

	// Note: membership/role checks moved to usecase where RequireRole is not used for AcceptInvite (public).
	checkGroupExist, err := v.groupRepo.CheckGroupExists(ctx, invite.GroupID)
	if err != nil {
		return nil, nil, err
	}

	if !checkGroupExist {
		return nil, nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group is not found"))
	}

	countMember, err := v.groupRepo.CountManagerAndMemberByGroupID(ctx, invite.GroupID)
	if err != nil {
		return nil, nil, err
	}

	if countMember >= 9 {
		return nil, nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("the group has reached maximum number of members"))
	}

	countViewer, err := v.groupRepo.CountViewerByGroupID(ctx, invite.GroupID)
	if err != nil {
		return nil, nil, err
	}

	if invite.Role == string(enum.GroupRoleViewer) && countViewer >= 10 {
		return nil, nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("the group has reached maximum number of viewers"))
	}

	return invite, user, nil
}

func (v *GroupValidator) ValidatePresignURLsRequest(ctx context.Context, req *appdto.GeneratePresignedURLsRequest) errorbase.AppError {
	if len(req.Files) == 0 {
		return errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("at least one file is required"))
	}

	for i, file := range req.Files {
		if file.ContentType == "" {
			return errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("content type is required for file at index "+string(i)))
		}
	}

	return nil

}
