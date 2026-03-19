package usecase

import (
	"context"
	appdto "team_service/internal/application/common/dto"
	apphelper "team_service/internal/application/common/helper"
	irepository "team_service/internal/application/common/interface/repository"
	istore "team_service/internal/application/common/interface/store"
	appmapper "team_service/internal/application/common/mapper"
	appvalidation "team_service/internal/application/common/validation"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	domainhelper "team_service/internal/domain/common/helper"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
)

type sprintUseCase struct {
	store      istore.Store
	sprintRepo irepository.SprintRepository
	validator  *appvalidation.SprintValidator
	authHelper *apphelper.AuthHelper
}

func (uc *sprintUseCase) CreateSprint(ctx context.Context, req *appdto.CreateSprintRequest) (*appdto.BaseResponse[appdto.SprintResponse], errorbase.AppError) {
	_, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleManager)
	if err != nil {
		return &appdto.BaseResponse[appdto.SprintResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	sprint, err := uc.validator.ValidateCreateSprint(ctx, req)
	if err != nil {
		return &appdto.BaseResponse[appdto.SprintResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	var createdSprint *entity.Sprint
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		createdSprint, err = repo.SprintRepository().CreateSprint(ctx, sprint)
		if err != nil {
			return err
		}

		if createdSprint == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("create sprint returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &appdto.BaseResponse[appdto.SprintResponse]{
		Data:  appmapper.ToSprintResponse(createdSprint),
		Error: nil,
	}, nil
}

func (uc *sprintUseCase) GetSprint(ctx context.Context, req *appdto.GetSprintRequest) (*appdto.BaseResponse[appdto.SprintResponse], errorbase.AppError) {
	_, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleViewer)
	if err != nil {
		return &appdto.BaseResponse[appdto.SprintResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	sprintID, err := uc.validator.ValidateGetSprint(ctx, req)
	if err != nil {
		return &appdto.BaseResponse[appdto.SprintResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	sprint, err := uc.sprintRepo.GetSprintByID(ctx, sprintID)
	if err != nil {
		return &appdto.BaseResponse[appdto.SprintResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	if sprint == nil {
		return &appdto.BaseResponse[appdto.SprintResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    errdict.ErrNotFound.Code,
				Message: errdict.ErrNotFound.Title,
				Detail:  domainhelper.Ptr("sprint not found"),
			},
		}, nil
	}

	if err != nil {
		return &appdto.BaseResponse[appdto.SprintResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	return &appdto.BaseResponse[appdto.SprintResponse]{
		Data:  appmapper.ToSprintResponse(sprint),
		Error: nil,
	}, nil
}

func (uc *sprintUseCase) ListSprints(ctx context.Context, req *appdto.ListSprintsRequest) (*appdto.BaseResponse[appdto.ListSprintsResponse], errorbase.AppError) {
	_, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleViewer)
	if err != nil {
		return &appdto.BaseResponse[appdto.ListSprintsResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	groupID, err := uc.validator.ValidateListSprints(ctx, req)
	if err != nil {
		return &appdto.BaseResponse[appdto.ListSprintsResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	sprints, err := uc.sprintRepo.GetSprintsByGroupID(ctx, groupID)
	if err != nil {
		return &appdto.BaseResponse[appdto.ListSprintsResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	sprintResponses := make([]appdto.SprintResponse, 0, len(sprints))
	for _, sprint := range sprints {
		mapped := appmapper.ToSprintResponse(sprint)
		if mapped == nil {
			continue
		}

		sprintResponses = append(sprintResponses, *mapped)
	}

	return &appdto.BaseResponse[appdto.ListSprintsResponse]{
		Data: &appdto.ListSprintsResponse{
			Sprints: sprintResponses,
			Total:   int32(len(sprintResponses)),
		},
		Error: nil,
	}, nil
}

func (uc *sprintUseCase) UpdateSprint(ctx context.Context, req *appdto.UpdateSprintRequest) (*appdto.BaseResponse[appdto.SprintResponse], errorbase.AppError) {
	_, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleManager)
	if err != nil {
		return &appdto.BaseResponse[appdto.SprintResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	payload, err := uc.validator.ValidateUpdateSprint(ctx, req)
	if err != nil {
		return &appdto.BaseResponse[appdto.SprintResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	var updatedSprint *entity.Sprint
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		updatedSprint, err = repo.SprintRepository().UpdateSprint(
			ctx,
			payload.SprintID,
			payload.Name,
			payload.Goal,
			payload.StartDate,
			payload.EndDate,
		)
		if err != nil {
			return err
		}

		if updatedSprint == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("update sprint returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &appdto.BaseResponse[appdto.SprintResponse]{
		Data:  appmapper.ToSprintResponse(updatedSprint),
		Error: nil,
	}, nil
}

func (uc *sprintUseCase) UpdateSprintStatus(ctx context.Context, req *appdto.UpdateSprintStatusRequest) (*appdto.BaseResponse[appdto.UpdateSprintStatusResponse], errorbase.AppError) {
	_, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleManager)
	if err != nil {
		return &appdto.BaseResponse[appdto.UpdateSprintStatusResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	payload, err := uc.validator.ValidateUpdateSprintStatus(ctx, req)
	if err != nil {
		return &appdto.BaseResponse[appdto.UpdateSprintStatusResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	var updatedSprint *entity.Sprint
	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		updatedSprint, err = repo.SprintRepository().UpdateSprintStatus(ctx, payload.SprintID, payload.Status)
		if err != nil {
			return err
		}

		if updatedSprint == nil {
			return errorbase.New(errdict.ErrInternal, errorbase.WithDetail("update sprint status returned nil"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &appdto.BaseResponse[appdto.UpdateSprintStatusResponse]{
		Data: &appdto.UpdateSprintStatusResponse{
			SprintID: updatedSprint.ID,
			Status:   updatedSprint.Status,
		},
		Error: nil,
	}, nil
}

func (uc *sprintUseCase) DeleteSprint(ctx context.Context, req *appdto.DeleteSprintRequest) (*appdto.BaseResponse[appdto.DeleteSprintResponse], errorbase.AppError) {
	_, err := uc.authHelper.RequireRole(ctx, enum.GroupRoleManager)
	if err != nil {
		return &appdto.BaseResponse[appdto.DeleteSprintResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	payload, err := uc.validator.ValidateDeleteSprint(ctx, req)
	if err != nil {
		return &appdto.BaseResponse[appdto.DeleteSprintResponse]{
			Data: nil,
			Error: &appdto.ErrorResponse{
				Code:    err.ErrorInfo().Code,
				Message: err.ErrorInfo().Title,
				Detail:  err.ErrorInfo().Detail,
			},
		}, nil
	}

	err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
		return repo.SprintRepository().DeleteSprint(ctx, payload.SprintID)
	})
	if err != nil {
		return nil, err
	}

	return &appdto.BaseResponse[appdto.DeleteSprintResponse]{
		Data:  &appdto.DeleteSprintResponse{Success: true},
		Error: nil,
	}, nil
}
