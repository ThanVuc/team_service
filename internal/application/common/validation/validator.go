package appvalidation

import (
	"context"
	"strings"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/proto/team_service"
)

// type Validator struct {}

func ValidateGroup(ctx context.Context, req *team_service.CreateGroupRequest) errorbase.AppError {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errorbase.New(errdict.ErrBadRequest)
	}
	if len(name) < 3 {
		return errorbase.New(errdict.ErrBadRequest)
	}

	if len(name) > 100 {
		return errorbase.New(errdict.ErrBadRequest)
	}

	if strings.Contains(name, "<") || strings.Contains(name, ">") {
		return errorbase.New(errdict.ErrBadRequest)
	}

	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)

		if desc != "" {
			if len(desc) > 500 {
				return errorbase.New(errdict.ErrBadRequest)
			}
		}
	}

	return nil

}
