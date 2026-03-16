package adapermapper

import (
	appdto "team_service/internal/application/common/dto"
	"team_service/internal/domain/enum"
	"team_service/proto/team_service"
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
