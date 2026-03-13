package appdto

import "time"

type UserSummaryDTO struct {
	ID     string
	Name   string
	Avatar *string
}

type GetUserRequest struct {
	UserID string
}

type GetUserResponse struct {
	ID     string
	Email  string
	Status string

	TimeZone string

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
