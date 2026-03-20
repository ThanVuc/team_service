package appmapper

import (
	appdto "team_service/internal/application/common/dto"
	errorbase "team_service/internal/domain/common/apperror"
)

func ToErrorResponse(err errorbase.AppError) *appdto.ErrorResponse {
	if err == nil {
		return nil
	}

	info := err.ErrorInfo()
	return &appdto.ErrorResponse{
		Code:    info.Code,
		Message: info.Title,
		Detail:  info.Detail,
	}
}
