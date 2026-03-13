package adapermapper

import (
	appdto "team_service/internal/application/common/dto"
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
