package enum

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive,
		UserStatusInactive:
		return true
	}
	return false
}
