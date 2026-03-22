package repository

import (
	"context"
	"fmt"

	appdto "team_service/internal/application/common/dto"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/internal/infrastructure/share/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

type GroupRepository struct {
	q *database.Queries
}

func NewGroupRepository(
	q *database.Queries,
) *GroupRepository {
	return &GroupRepository{q: q}
}

func (r *GroupRepository) CreateGroup(
	ctx context.Context,
	group *entity.Group,
	userID string,
) (*entity.Group, errorbase.AppError) {
	var groupID pgtype.UUID
	if err := groupID.Scan(group.ID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	var ownerID pgtype.UUID
	if err := ownerID.Scan(userID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse owner id"),
		)
	}

	var desc pgtype.Text
	if group.Description != nil {
		desc = pgtype.Text{
			String: *group.Description,
			Valid:  true,
		}
	}

	dbGroup, err := r.q.CreateGroup(ctx, database.CreateGroupParams{
		ID:          groupID,
		Name:        group.Name,
		Description: desc,
		OwnerID:     ownerID,
	})

	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to create group name=%s owner=%s", group.Name, userID)),
		)
	}

	return &entity.Group{
		ID:          dbGroup.ID.String(),
		Name:        dbGroup.Name,
		Description: group.Description,
		OwnerID:     dbGroup.OwnerID.String(),
		CreatedAt:   dbGroup.CreatedAt.Time,
		UpdatedAt:   dbGroup.UpdatedAt.Time,
	}, nil
}

func (r *GroupRepository) CountGroupsByOwner(
	ctx context.Context,
	ownerID string,
) (int64, errorbase.AppError) {

	var ownerUUID pgtype.UUID
	if err := ownerUUID.Scan(ownerID); err != nil {
		return 0, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse owner id"),
		)
	}

	count, err := r.q.CountGroupsByOwner(ctx, ownerUUID)
	if err != nil {
		return 0, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to count groups for owner=%s", ownerID)),
		)
	}

	return count, nil
}

func (r *GroupRepository) AddGroupMember(
	ctx context.Context,
	member *entity.GroupMember,
) errorbase.AppError {

	var id pgtype.UUID
	if err := id.Scan(member.ID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group member id"),
		)
	}

	var groupID pgtype.UUID
	if err := groupID.Scan(member.GroupID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	var userID pgtype.UUID
	if err := userID.Scan(member.UserID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse user id"),
		)
	}

	err := r.q.CreateGroupMember(ctx, database.CreateGroupMemberParams{
		ID:      id,
		GroupID: groupID,
		UserID:  userID,
		Role:    string(member.Role),
	})

	if err != nil {
		return errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to add member user=%s to group=%s", member.UserID, member.GroupID)),
		)
	}

	return nil
}

func (r *GroupRepository) GetGroupsByUserID(
	ctx context.Context,
	userID string,
) (*appdto.ListGroupsResponse, errorbase.AppError) {
	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse user id"),
		)
	}

	rows, err := r.q.GetGroupsByUserID(ctx, userUUID)
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to list groups for user=%s", userID)),
		)
	}

	items := make([]appdto.ListGroupItem, 0, len(rows))
	for _, row := range rows {
		owner, err := r.q.GetOwnerByGroupID(ctx, row.ID)
		if err != nil {
			return nil, errorbase.Wrap(
				err,
				errdict.ErrInternal,
				errorbase.WithDetail(fmt.Sprintf("failed to get owner for group=%s", row.ID.String())),
			)
		}

		var ownerAvatar *string
		if owner.OwnerImage != "" {
			ownerAvatar = utils.Ptr(owner.OwnerImage)
		}

		items = append(items, appdto.ListGroupItem{
			ID:   row.ID.String(),
			Name: row.Name,
			Owner: appdto.OwnerDTO{
				ID:     owner.OwnerID.String(),
				Email:  owner.OwnerEmail,
				Avatar: ownerAvatar,
			},
			MyRole:      enum.GroupRole(row.MyRole),
			MemberTotal: int(row.MemberTotal),
			AvatarURL:   row.AvatarUrl,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		})
	}

	return &appdto.ListGroupsResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func (r *GroupRepository) GetGroupByID(
	ctx context.Context,
	groupID string,
) (*entity.Group, int32, string, errorbase.AppError) {

	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return nil, 0, "", errorbase.Wrap(
			err,
			errdict.ErrBadRequest,
			errorbase.WithDetail(fmt.Sprintf("invalid group id=%s", groupID)),
		)
	}

	g, err := r.q.GetGroupByID(ctx, groupUUID)
	if err != nil {
		return nil, 0, "", errorbase.Wrap(
			err,
			errdict.ErrNotFound,
			errorbase.WithDetail(fmt.Sprintf("group not found id=%s", groupID)),
		)
	}

	memberCount, err := r.q.CountGroupMembersByGroupID(ctx, groupUUID)
	if err != nil {
		return nil, 0, "", errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to count members for group=%s", groupID)),
		)
	}

	var sprintName string
	sprint, err := r.q.GetSprintByGroupID(ctx, groupUUID)
	if err == nil {
		sprintName = sprint.Name
	}

	var description *string
	if g.Description.Valid {
		description = &g.Description.String
	}

	group := &entity.Group{
		ID:          g.ID.String(),
		Name:        g.Name,
		Description: description,
		OwnerID:     g.OwnerID.String(),
		CreatedAt:   g.CreatedAt.Time,
		UpdatedAt:   g.UpdatedAt.Time,
	}

	return group, int32(memberCount), sprintName, nil
}

func (r *GroupRepository) GetRoleByUserIDAndGroupID(
	ctx context.Context,
	userID string,
	groupID string,
) (string, errorbase.AppError) {
	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return "", errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse user id"),
		)
	}

	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return "", errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	role, err := r.q.GetRoleByGroupIDAndUserID(ctx, database.GetRoleByGroupIDAndUserIDParams{
		GroupID: groupUUID,
		UserID:  userUUID,
	})

	if err != nil {
		return "", errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to get role for user=%s in group=%s", userID, groupID)),
		)
	}

	return role, nil
}

func (r *GroupRepository) CheckGroupExists(
	ctx context.Context,
	groupID string,
) (bool, errorbase.AppError) {
	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return false, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	exists, err := r.q.CheckGroupExists(ctx, groupUUID)
	if err != nil {
		return false, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to check if group exists id=%s", groupID)),
		)
	}

	return exists, nil
}

func (r *GroupRepository) UpdateGroup(
	ctx context.Context,
	group *entity.Group,
) (*entity.Group, errorbase.AppError) {
	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(group.ID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	var desc pgtype.Text
	if group.Description != nil {
		desc = pgtype.Text{
			String: *group.Description,
			Valid:  true,
		}
	}

	var name pgtype.Text
	name = pgtype.Text{
		String: group.Name,
		Valid:  true,
	}

	dbGroup, err := r.q.UpdateGroup(ctx, database.UpdateGroupParams{
		ID:          groupUUID,
		Name:        name,
		Description: desc,
	})

	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to update group id=%s", group.ID)),
		)
	}

	return &entity.Group{
		ID:          dbGroup.ID.String(),
		Name:        dbGroup.Name,
		Description: &dbGroup.Description.String,
		OwnerID:     dbGroup.OwnerID.String(),
		CreatedAt:   dbGroup.CreatedAt.Time,
		UpdatedAt:   dbGroup.UpdatedAt.Time,
	}, nil

}

func (r *GroupRepository) DeleteGroup(
	ctx context.Context,
	groupID string,
) errorbase.AppError {
	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	err := r.q.DeleteGroup(ctx, groupUUID)
	if err != nil {
		return errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to delete group id=%s", groupID)),
		)
	}

	return nil
}

func (r *GroupRepository) CountManagerAndMemberByGroupID(
	ctx context.Context,
	groupID string,
) (int64, errorbase.AppError) {
	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return 0, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	count, err := r.q.CountManagerAndMemberByGroupID(ctx, groupUUID)
	if err != nil {
		return 0, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to count manager and member for group=%s", groupID)),
		)
	}

	return count, nil
}

func (r *GroupRepository) UpdateMemberRole(
	ctx context.Context,
	userID string,
	groupID string,
	newRole string,
) (*appdto.MemberResponse, errorbase.AppError) {
	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse user id"),
		)
	}

	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	member, err := r.q.UpdateRoleMember(ctx, database.UpdateRoleMemberParams{
		UserID:  userUUID,
		GroupID: groupUUID,
		Role:    newRole,
	})

	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to update member role user=%s in group=%s", userID, groupID)),
		)
	}

	return &appdto.MemberResponse{
		ID:       member.ID.String(),
		Email:    member.Email,
		Avatar:   utils.Ptr(member.AvatarUrl.String),
		Role:     enum.GroupRole(member.Role),
		JoinedAt: member.JoinedAt.Time,
	}, nil
}

func (r *GroupRepository) RemoveMember(
	ctx context.Context,
	groupID string,
	userID string,
) errorbase.AppError {
	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse user id"),
		)
	}

	err := r.q.RemoveMember(ctx, database.RemoveMemberParams{
		GroupID: groupUUID,
		UserID:  userUUID,
	})

	if err != nil {
		return errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to remove member user=%s from group=%s", userID, groupID)),
		)
	}

	return nil
}

func (r *GroupRepository) CheckMemberExistsByEmail(
	ctx context.Context,
	groupID string,
	email string,
) (bool, errorbase.AppError) {
	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return false, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	exists, err := r.q.CheckMemberExistsByEmail(ctx, database.CheckMemberExistsByEmailParams{
		GroupID: groupUUID,
		Email:   email,
	})

	if err != nil {
		return false, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to check if member with email=%s exists in group=%s", email, groupID)),
		)
	}

	return exists, nil
}

func (r *GroupRepository) GetSimpleUsersByGroupID(
	ctx context.Context,
	groupID string,
) ([]*appdto.SimpleUserResponse, errorbase.AppError) {
	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	rows, err := r.q.GetSimpleUserByGroupID(ctx, groupUUID)
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to get simple users for group=%s", groupID)),
		)
	}

	var users []*appdto.SimpleUserResponse
	for _, row := range rows {
		users = append(users, &appdto.SimpleUserResponse{
			ID:        row.ID.String(),
			Email:     row.Email,
			AvatarURL: utils.Ptr(row.AvatarUrl.String),
		})
	}

	return users, nil
}
