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
)

type sprintUseCase struct {
	store      istore.Store
	sprintRepo irepository.SprintRepository
	validator  *appvalidation.SprintValidator
	authHelper *apphelper.AuthHelper
}

func (uc *sprintUseCase) CreateSprint(ctx context.Context, req *appdto.CreateSprintRequest) (*appdto.BaseResponse[appdto.SprintResponse], errorbase.AppError) {
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
	_ = ctx
	_ = req

	return &appdto.BaseResponse[appdto.SprintResponse]{
		Data:  nil,
		Error: sprintNotImplementedError(),
	}, nil
}

func (uc *sprintUseCase) ListSprints(ctx context.Context, req *appdto.ListSprintsRequest) (*appdto.BaseResponse[appdto.ListSprintsResponse], errorbase.AppError) {
	_ = ctx
	_ = req

	return &appdto.BaseResponse[appdto.ListSprintsResponse]{
		Data:  nil,
		Error: sprintNotImplementedError(),
	}, nil
}

func (uc *sprintUseCase) UpdateSprint(ctx context.Context, req *appdto.UpdateSprintRequest) (*appdto.BaseResponse[appdto.SprintResponse], errorbase.AppError) {
	_ = ctx
	_ = req

	return &appdto.BaseResponse[appdto.SprintResponse]{
		Data:  nil,
		Error: sprintNotImplementedError(),
	}, nil
}

func (uc *sprintUseCase) UpdateSprintStatus(ctx context.Context, req *appdto.UpdateSprintStatusRequest) (*appdto.BaseResponse[appdto.UpdateSprintStatusResponse], errorbase.AppError) {
	_ = ctx
	_ = req

	return &appdto.BaseResponse[appdto.UpdateSprintStatusResponse]{
		Data:  nil,
		Error: sprintNotImplementedError(),
	}, nil
}

func (uc *sprintUseCase) DeleteSprint(ctx context.Context, req *appdto.DeleteSprintRequest) (*appdto.BaseResponse[appdto.DeleteSprintResponse], errorbase.AppError) {
	_ = ctx
	_ = req

	return &appdto.BaseResponse[appdto.DeleteSprintResponse]{
		Data:  nil,
		Error: sprintNotImplementedError(),
	}, nil
}

func sprintNotImplementedError() *appdto.ErrorResponse {
	return &appdto.ErrorResponse{
		Code:    errdict.ErrUnprocessable.Code,
		Message: "Sprint use case is not implemented yet",
		Detail:  domainhelper.Ptr("Controller/DTO/DI wiring is ready; business logic was intentionally omitted."),
	}
}
