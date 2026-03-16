package repository

import (
	"context"
	"fmt"

	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/infrastructure/persistence/db/database"

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
