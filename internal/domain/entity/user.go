package entity

import (
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/enum"
	"time"
)

type User struct {
	ID                   string
	Email                string
	Status               enum.UserStatus
	CreatedAt            time.Time
	AvatarURL            *string
	HasEmailNotification bool
	HasPushNotification  bool
}

func (u *User) SetNotificationPreference(email bool, push bool) {
	u.HasEmailNotification = email
	u.HasPushNotification = push
}

func CreateUser(
	id string,
	email string,
	now time.Time,
	status enum.UserStatus,
	avatarURL *string,
) (*User, errorbase.AppError) {
	if !status.IsValid() {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("invalid user status"))
	}

	return &User{
		ID:                   id,
		Email:                email,
		Status:               enum.UserStatusActive,
		CreatedAt:            now,
		HasEmailNotification: false,
		HasPushNotification:  true,
		AvatarURL:            avatarURL,
	}, nil
}
