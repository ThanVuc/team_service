package appdto

import (
	"team_service/internal/domain/enum"
	"time"
)

type ListMembersRequest struct {
	GroupID string
}

type UpdateMemberRoleRequest struct {
	GroupID string
	UserID  string
	Role    enum.GroupRole
}

type RemoveMemberRequest struct {
	GroupID string
	UserID  string
}

type MemberResponse struct {
	ID     string
	Name   string
	Email  string
	Avatar *string

	Role enum.GroupRole

	JoinedAt time.Time
}
