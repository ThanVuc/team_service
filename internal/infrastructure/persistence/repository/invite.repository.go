package repository

import (
	"context"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/infrastructure/persistence/db/database"

	"github.com/jackc/pgx/v5/pgtype"
)

type InviteRepository struct {
	q *database.Queries
}

func NewInviteRepository(
	q *database.Queries,
) *InviteRepository {
	return &InviteRepository{
		q: q,
	}
}

func (r *InviteRepository) CreateInvite(ctx context.Context, invite *entity.Invite) (*entity.Invite, errorbase.AppError) {
	var id pgtype.UUID
	if err := id.Scan(invite.ID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse invite id"),
		)
	}

	var groupID pgtype.UUID
	if err := groupID.Scan(invite.GroupID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	var email pgtype.Text
	if invite.Email != nil {
		email = pgtype.Text{
			String: *invite.Email,
			Valid:  true,
		}
	}

	var senderID pgtype.UUID
	if err := senderID.Scan(invite.CreatedBy); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse sender id"),
		)
	}

	var expiresAt pgtype.Timestamptz
	if err := expiresAt.Scan(invite.ExpiresAt); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse expires at"),
		)
	}

	dbInvite, err := r.q.CreateInvite(ctx, database.CreateInviteParams{
		ID:        id,
		GroupID:   groupID,
		Token:     invite.Token,
		Role:      invite.Role,
		Email:     email,
		ExpiresAt: expiresAt,
		CreatedBy: senderID,
	})

	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail("failed to create invite"),
		)
	}

	newInvite := &entity.Invite{
		ID:        invite.ID,
		GroupID:   invite.GroupID,
		Token:     dbInvite.Token,
		Role:      invite.Role,
		Email:     invite.Email,
		ExpiresAt: dbInvite.ExpiresAt.Time,
		CreatedBy: invite.CreatedBy,
		CreatedAt: dbInvite.CreatedAt.Time,
	}

	return newInvite, nil
}

func (r InviteRepository) CheckPendingInviteExists(ctx context.Context, groupID string, email string) (bool, errorbase.AppError) {
	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return false, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	var emailText pgtype.Text
	if email != "" {
		emailText = pgtype.Text{
			String: email,
			Valid:  true,
		}
	}

	exists, err := r.q.CheckPendingInvite(ctx, database.CheckPendingInviteParams{
		GroupID: groupUUID,
		Email:   emailText,
	})

	if err != nil {
		return false, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail("failed to check pending invite"),
		)
	}

	return exists, nil

}

func (r InviteRepository) GetInviteByToken(ctx context.Context, token string) (*entity.Invite, errorbase.AppError) {
	dbInvite, err := r.q.GetInviteByToken(ctx, token)
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail("failed to get invite by token"),
		)
	}

	invite := &entity.Invite{
		ID:        dbInvite.ID.String(),
		GroupID:   dbInvite.GroupID.String(),
		Token:     dbInvite.Token,
		Role:      dbInvite.Role,
		Email:     nil,
		ExpiresAt: dbInvite.ExpiresAt.Time,
		CreatedBy: dbInvite.CreatedBy.String(),
		CreatedAt: dbInvite.CreatedAt.Time,
	}

	if dbInvite.Email.Valid {
		invite.Email = &dbInvite.Email.String
	}

	return invite, nil
}
