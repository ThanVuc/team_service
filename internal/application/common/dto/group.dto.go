package appdto

import (
	"team_service/internal/domain/enum"
	"time"
)

type CreateGroupRequest struct {
	Name        string
	Description *string
}

type UpdateGroupRequest struct {
	GroupID     string
	Name        *string
	Description *string
	AvatarURL   *string
}

type DeleteGroupRequest struct {
	GroupID string
}

type GetGroupRequest struct {
	GroupID string
}

type ListGroupsRequest struct{}

type ListGroupItem struct {
	ID          string
	Name        string
	Owner       OwnerDTO
	MyRole      enum.GroupRole
	MemberTotal int
	AvatarURL   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ListGroupsResponse struct {
	Items []ListGroupItem
	Total int
}

type GroupResponse struct {
	ID          string
	Name        string
	Description *string

	Owner OwnerDTO

	MyRole       enum.GroupRole
	ActiveSprint *string
	MemberTotal  int

	AvatarURL *string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type OwnerDTO struct {
	ID     string
	Email  string
	Avatar *string
}

type DeleteGroupResponse struct {
	Success bool
}

type PresignFileItem struct {
	Index       int
	ContentType string
	FileName    string
}

type GeneratePresignedURLsRequest struct {
	Files []PresignFileItem
}

type PresignedFileItem struct {
	Index      int
	PresignUrl string
}

type GeneratePresignedURLsResponse struct {
	Files []PresignedFileItem
}
