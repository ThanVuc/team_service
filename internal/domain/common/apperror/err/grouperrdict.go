package errdict

import (
	errorbase "team_service/internal/domain/common/apperror"
	domainhelper "team_service/internal/domain/common/helper"
)

var (
	ErrInvalidUUID = errorbase.ErrorInfo{
		Code:   "ts.validation.invalid-uuid",
		Title:  "Invalid UUID",
		Detail: domainhelper.Ptr("The provided UUID is not valid. Please ensure it is in the correct format."),
	}

	ErrInviteExpired = errorbase.ErrorInfo{
		Code:   "ts.validation.invite-expired",
		Title:  "Invite Expired",
		Detail: domainhelper.Ptr("The invite has expired."),
	}

	ErrGroupMaxMembersReached = errorbase.ErrorInfo{
		Code:   "ts.validation.group-max-members-reached",
		Title:  "Group Member Limit Reached",
		Detail: domainhelper.Ptr("The group has reached maximum number of members."),
	}
)
