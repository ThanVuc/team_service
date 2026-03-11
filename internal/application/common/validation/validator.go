package appvalidation

import (
	"context"
	"strings"
	coreerror "team_service/internal/domain/common/apperror"
	"team_service/internal/domain/common/apperror/errordictionary"
	"team_service/proto/team_service"
)

// type Validator struct {}

func ValidateGroup(ctx context.Context, req *team_service.CreateGroupRequest) coreerror.AppError {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errordictionary.ErrGroupBadRequest
	}
	if len(name) < 3 {
		return errordictionary.ErrGroupBadRequest
	}

	if len(name) > 100 {
		return errordictionary.ErrGroupBadRequest
	}

	if strings.Contains(name, "<") || strings.Contains(name, ">") {
		return errordictionary.ErrGroupBadRequest
	}

	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)

		if desc != "" {
			if len(desc) > 500 {
				return errordictionary.ErrGroupBadRequest
			}
		}
	}

	return nil

}
