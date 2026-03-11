package repository

import (
	"context"
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

func (r *GroupRepository) CreateGroup(ctx context.Context, req *team_service.CreateGroupRequest, userID string) (*database.Group, error) {
	groupID := pgtype.UUID{
		Bytes: uuid.New(),
		Valid: true,
	}
	var ownerID pgtype.UUID
	if err := ownerID.Scan(userID); err != nil {
		return nil, err
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
		return nil, err
	}

	return &group, nil
}

func (r *GroupRepository) CountGroupsByOwner(ctx context.Context, ownerID string) (int64, error) {
	var ownerIDUUID pgtype.UUID
	if err := ownerIDUUID.Scan(ownerID); err != nil {
		return 0, err
	}

	return r.q.CountGroupsByOwner(ctx, ownerIDUUID)
}

func (r *GroupRepository) GetUserByID(ctx context.Context, userID string) (*database.GetUserByIDRow, error) {
	var userIDUUID pgtype.UUID
	if err := userIDUUID.Scan(userID); err != nil {
		return nil, err
	}
	user, err := r.q.GetUserByID(ctx, userIDUUID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *GroupRepository) AddGroupMember(ctx context.Context, arg database.CreateGroupMemberParams) error {
	return r.q.CreateGroupMember(ctx, arg)
}
