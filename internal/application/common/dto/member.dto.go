package appdto

import (
	"team_service/internal/domain/enum"
	"time"
)

type ListMembersRequest struct {
	GroupID string
}

type UpdateMemberRoleRequest struct {
	GroupID  string
	MemberId string
	Role     enum.GroupRole
}

type RemoveMemberRequest struct {
	GroupID  string
	MemberId string
}

type MemberResponse struct {
	ID     string
	Name   string
	Email  string
	Avatar *string

	Role enum.GroupRole

	JoinedAt time.Time
}

type ListMembersResponse struct {
	Members []MemberResponse
	Total   int
}

type RemoveMemberResponse struct {
	Success bool
}
