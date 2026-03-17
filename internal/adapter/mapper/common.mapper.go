package adapermapper

import (
	appdto "team_service/internal/application/common/dto"
	"team_service/internal/domain/enum"
	"team_service/proto/team_service"
	"time"
)

func ToProtoError(err *appdto.ErrorResponse) *team_service.Error {
	if err == nil {
		return nil
	}

	return &team_service.Error{
		Code:    err.Code,
		Message: err.Message,
		Details: err.Detail,
	}
}

func MapGroupRole(role enum.GroupRole) team_service.GroupRole {
	switch role {
	case enum.GroupRoleOwner:
		return team_service.GroupRole_GROUP_ROLE_OWNER
	case enum.GroupRoleManager:
		return team_service.GroupRole_GROUP_ROLE_MANAGER
	case enum.GroupRoleMember:
		return team_service.GroupRole_GROUP_ROLE_MEMBER
	case enum.GroupRoleViewer:
		return team_service.GroupRole_GROUP_ROLE_VIEWER
	default:
		return team_service.GroupRole_GROUP_ROLE_UNSPECIFIED
	}
}

func MapSprintStatus(status enum.SprintStatus) team_service.SprintStatus {
	switch status {
	case enum.SprintStatusDraft:
		return team_service.SprintStatus_SPRINT_STATUS_DRAFT
	case enum.SprintStatusActive:
		return team_service.SprintStatus_SPRINT_STATUS_ACTIVE
	case enum.SprintStatusCompleted:
		return team_service.SprintStatus_SPRINT_STATUS_COMPLETED
	case enum.SprintStatusCancelled:
		return team_service.SprintStatus_SPRINT_STATUS_CANCELLED
	default:
		return team_service.SprintStatus_SPRINT_STATUS_UNSPECIFIED
	}
}

func MapProtoSprintStatus(status team_service.SprintStatus) enum.SprintStatus {
	switch status {
	case team_service.SprintStatus_SPRINT_STATUS_DRAFT:
		return enum.SprintStatusDraft
	case team_service.SprintStatus_SPRINT_STATUS_ACTIVE:
		return enum.SprintStatusActive
	case team_service.SprintStatus_SPRINT_STATUS_COMPLETED:
		return enum.SprintStatusCompleted
	case team_service.SprintStatus_SPRINT_STATUS_CANCELLED:
		return enum.SprintStatusCancelled
	default:
		return ""
	}
}

func FromDateToTime(date *team_service.Date) time.Time {
	if date == nil {
		return time.Time{}
	}

	t := time.Date(int(date.Year), time.Month(date.Month), int(date.Day), 0, 0, 0, 0, time.UTC)
	return t
}

func FromTimeToDate(t time.Time) *team_service.Date {
	if t.IsZero() {
		return nil
	}

	return &team_service.Date{
		Year:  int32(t.Year()),
		Month: int32(t.Month()),
		Day:   int32(t.Day()),
	}
}
