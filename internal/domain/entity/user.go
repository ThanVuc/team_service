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
	TimeZone             string
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
	timeZone string,
	now time.Time,
	status enum.UserStatus,
) (*User, errorbase.AppError) {
	if !status.IsValid() {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("invalid user status"))
	}

	return &User{
		ID:                   id,
		Email:                email,
		Status:               enum.UserStatusActive,
		TimeZone:             timeZone,
		CreatedAt:            now,
		HasEmailNotification: false,
		HasPushNotification:  true,
	}, nil
}
