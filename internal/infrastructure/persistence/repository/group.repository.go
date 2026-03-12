package repository

import (
	"context"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/proto/team_service"

	"github.com/google/uuid"
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

func (r *GroupRepository) CreateGroup(ctx context.Context, req *team_service.CreateGroupRequest, userID string) (*database.Group, errorbase.AppError) {
	groupID := pgtype.UUID{
		Bytes: uuid.New(),
		Valid: true,
	}
	var ownerID pgtype.UUID
	if err := ownerID.Scan(userID); err != nil {
		return nil, errorbase.New(errdict.ErrInternal)
	}

	var desc pgtype.Text
	if req.Description != nil {
		desc = pgtype.Text{
			String: *req.Description,
			Valid:  true,
		}
	}

	group, err := r.q.CreateGroup(ctx, database.CreateGroupParams{
		ID:          groupID,
		Name:        req.Name,
		Description: desc,
		OwnerID:     ownerID,
	})
	if err != nil {
		return nil, errorbase.New(errdict.ErrInternal)
	}

	return &group, nil
}

func (r *GroupRepository) CountGroupsByOwner(ctx context.Context, ownerID string) (int64, errorbase.AppError) {
	var ownerIDUUID pgtype.UUID
	if err := ownerIDUUID.Scan(ownerID); err != nil {
		return 0, errorbase.New(errdict.ErrInternal)
	}

	count, err := r.q.CountGroupsByOwner(ctx, ownerIDUUID)
	if err != nil {
		return 0, errorbase.Wrap(
			err,
			errdict.ErrBadRequest,
		)
	}
	return count, nil
}

func (r *GroupRepository) GetUserByID(ctx context.Context, userID string) (*database.GetUserByIDRow, errorbase.AppError) {
	var userIDUUID pgtype.UUID
	if err := userIDUUID.Scan(userID); err != nil {
		return nil, errorbase.New(errdict.ErrInternal)
	}
	user, err := r.q.GetUserByID(ctx, userIDUUID)
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrBadRequest,
		)
	}

	return &user, nil
}

func (r *GroupRepository) AddGroupMember(ctx context.Context, arg database.CreateGroupMemberParams) errorbase.AppError {
	err := r.q.CreateGroupMember(ctx, arg)
	if err != nil {
		return errorbase.Wrap(
			err,
			errdict.ErrInternal,
		)
	}
	return nil
}

func (r *GroupRepository) GetGroupByID(ctx context.Context, user, groupID string) (*database.GetGroupByIDRow, int32, string, string, errorbase.AppError) {
	var groupIDUUID pgtype.UUID
	if err := groupIDUUID.Scan(groupID); err != nil {
		return nil, 0, "", "", errorbase.Wrap(err , errdict.ErrBadRequest)
	}

	var userIDUUID pgtype.UUID
	if err := userIDUUID.Scan(user); err != nil {
		return nil, 0, "", "", errorbase.Wrap(err , errdict.ErrBadRequest)
	}

	group, err := r.q.GetGroupByID(ctx, groupIDUUID)
	if err != nil {
		return nil, 0, "", "", errorbase.Wrap(err , errdict.ErrBadRequest)
	}

	count, err := r.q.CountGroupMembersByGroupID(ctx, groupIDUUID)
	if err != nil {
		return nil, 0, "", "", errorbase.Wrap(err , errdict.ErrBadRequest)
	}

	sprintName := ""
	sprint, err := r.q.GetSprintByGroupID(ctx, groupIDUUID)
	if err == nil {
		sprintName = sprint.Name
	}

	payload := database.GetRoleByGroupIDAndUserIDParams{
		GroupID: groupIDUUID,
		UserID:  userIDUUID,
	}

	myRole, err := r.q.GetRoleByGroupIDAndUserID(ctx, payload)
	if err != nil {
		return nil, 0, "", "", errorbase.Wrap(err , errdict.ErrBadRequest)
	}

	return &group, int32(count), sprintName, myRole, nil
}
