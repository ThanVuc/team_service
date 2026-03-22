package appdto

import (
	"team_service/internal/domain/enum"
	"time"
)

type UserSummaryDTO struct {
	ID     string
	Name   string
	Email  string
	Avatar *string
}

type GetUserRequest struct {
	UserID string
}

type GetUserResponse struct {
	ID     string
	Email  string
	Status string

	AvatarURL *string

	CreatedAt time.Time
}

type ConfigureNotificationRequest struct {
	UseEmailNotification bool
	UseAppNotification   *bool
}

type ConfigureNotificationResponse struct {
	Success bool
}

type UserWithPermission struct {
	ID       string
	GroupId  string
	Email    string
	Status   enum.UserStatus
	Role     enum.GroupRole
	JoinedAt time.Time
}

type UserOutboxPayload struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
	Fullname  string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

type SimpleUserResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}
