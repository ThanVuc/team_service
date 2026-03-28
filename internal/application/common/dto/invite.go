package appdto

import (
	"team_service/internal/domain/enum"
	"time"
)

type CreateInviteRequest struct {
	GroupID string
	Role    enum.GroupRole
	Email   *string
}

type InviteResponse struct {
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type AcceptInviteRequest struct {
	Code string
}

type AcceptInviteResponse struct {
	Location string
}
