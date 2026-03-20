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

func MapProtoGroupRole(role team_service.GroupRole) enum.GroupRole {
	switch role {
	case team_service.GroupRole_GROUP_ROLE_OWNER:
		return enum.GroupRoleOwner
	case team_service.GroupRole_GROUP_ROLE_MANAGER:
		return enum.GroupRoleManager
	case team_service.GroupRole_GROUP_ROLE_MEMBER:
		return enum.GroupRoleMember
	case team_service.GroupRole_GROUP_ROLE_VIEWER:
		return enum.GroupRoleViewer
	default:
		return ""
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

func MapWorkStatus(status enum.WorkStatus) team_service.WorkStatus {
	switch status {
	case enum.WorkStatusTodo:
		return team_service.WorkStatus_WORK_STATUS_TODO
	case enum.WorkStatusInProgress:
		return team_service.WorkStatus_WORK_STATUS_IN_PROGRESS
	case enum.WorkStatusInReview:
		return team_service.WorkStatus_WORK_STATUS_IN_REVIEW
	case enum.WorkStatusDone:
		return team_service.WorkStatus_WORK_STATUS_DONE
	default:
		return team_service.WorkStatus_WORK_STATUS_UNSPECIFIED
	}
}

func MapProtoWorkStatus(status team_service.WorkStatus) enum.WorkStatus {
	switch status {
	case team_service.WorkStatus_WORK_STATUS_TODO:
		return enum.WorkStatusTodo
	case team_service.WorkStatus_WORK_STATUS_IN_PROGRESS:
		return enum.WorkStatusInProgress
	case team_service.WorkStatus_WORK_STATUS_IN_REVIEW:
		return enum.WorkStatusInReview
	case team_service.WorkStatus_WORK_STATUS_DONE:
		return enum.WorkStatusDone
	default:
		return ""
	}
}

func MapWorkPriority(priority enum.WorkPriority) team_service.WorkPriority {
	switch priority {
	case enum.WorkPriorityLow:
		return team_service.WorkPriority_WORK_PRIORITY_LOW
	case enum.WorkPriorityMedium:
		return team_service.WorkPriority_WORK_PRIORITY_MEDIUM
	case enum.WorkPriorityHigh:
		return team_service.WorkPriority_WORK_PRIORITY_HIGH
	default:
		return team_service.WorkPriority_WORK_PRIORITY_UNSPECIFIED
	}
}

func MapProtoWorkPriority(priority team_service.WorkPriority) enum.WorkPriority {
	switch priority {
	case team_service.WorkPriority_WORK_PRIORITY_LOW:
		return enum.WorkPriorityLow
	case team_service.WorkPriority_WORK_PRIORITY_MEDIUM:
		return enum.WorkPriorityMedium
	case team_service.WorkPriority_WORK_PRIORITY_HIGH:
		return enum.WorkPriorityHigh
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
