package entity

import (
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/enum"
	"time"
)

type Invite struct {
	ID        string
	GroupID   string
	Token     string
	Role      string
	Email     *string
	ExpiresAt time.Time
	CreatedBy string
	CreatedAt time.Time
}

func NewInvite(
	id string,
	groupID string,
	token string,
	role enum.GroupRole,
	email *string,
	expiresAt time.Time,
	createdBy string,
	now time.Time,
) (*Invite, errorbase.AppError) {

	if groupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest)
	}

	if token == "" {
		return nil, errorbase.New(errdict.ErrBadRequest)
	}

	if !role.IsValid() {
		return nil, errorbase.New(errdict.ErrBadRequest)
	}

	if expiresAt.Before(now) {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("invalid expiration"))
	}

	return &Invite{
		ID:        id,
		GroupID:   groupID,
		Token:     token,
		Role:      string(role),
		Email:     email,
		ExpiresAt: expiresAt,
		CreatedBy: createdBy,
		CreatedAt: now,
	}, nil
}

func (i *Invite) IsExpired(now time.Time) bool {
	return now.After(i.ExpiresAt)
}
