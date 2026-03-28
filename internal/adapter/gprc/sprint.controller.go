package grpccontroller

import (
	"context"
	adapermapper "team_service/internal/adapter/mapper"
	"team_service/internal/application/usecase"
	"team_service/internal/infrastructure/share/utils"
	"team_service/proto/common"
	"team_service/proto/team_service"

	"github.com/thanvuc/go-core-lib/log"
)

type SprintController struct {
	team_service.UnimplementedSprintServiceServer
	sprintUseCase usecase.SprintUseCase
	logger        log.LoggerV2
}

func NewSprintController(
	sprintUseCase usecase.SprintUseCase,
	logger log.LoggerV2,
) *SprintController {
	return &SprintController{
		sprintUseCase: sprintUseCase,
		logger:        logger,
	}
}

func (c *SprintController) CreateSprint(ctx context.Context, req *team_service.CreateSprintRequest) (*team_service.CreateSprintResponse, error) {
	createSprintReq := adapermapper.ToCreateSprintDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, createSprintReq, c.sprintUseCase.CreateSprint)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToCreateSprintGrpcResponse(resp), nil
}

func (c *SprintController) GetSprint(ctx context.Context, req *common.IDRequest) (*team_service.GetSprintResponse, error) {
	getSprintReq := adapermapper.ToGetSprintDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, getSprintReq, c.sprintUseCase.GetSprint)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToGetSprintGrpcResponse(resp), nil
}

func (c *SprintController) ListSprints(ctx context.Context, req *team_service.ListSprintsRequest) (*team_service.ListSprintsResponse, error) {
	listSprintReq := adapermapper.ToListSprintsDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, listSprintReq, c.sprintUseCase.ListSprints)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToListSprintsGrpcResponse(resp), nil
}

func (c *SprintController) GetSimpleSprints(ctx context.Context, req *common.IDRequest) (*team_service.GetSimpleSprintsResponse, error) {
	getSimpleSprintsReq := adapermapper.ToGetSimpleSprintsDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, getSimpleSprintsReq, c.sprintUseCase.GetSimpleSprints)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToGetSimpleSprintsGrpcResponse(resp), nil
}

func (c *SprintController) UpdateSprint(ctx context.Context, req *team_service.UpdateSprintRequest) (*team_service.UpdateSprintResponse, error) {
	updateSprintReq := adapermapper.ToUpdateSprintDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, updateSprintReq, c.sprintUseCase.UpdateSprint)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToUpdateSprintGrpcResponse(resp), nil
}

func (c *SprintController) UpdateSprintStatus(ctx context.Context, req *team_service.UpdateSprintStatusRequest) (*team_service.UpdateSprintStatusResponse, error) {
	updateSprintStatusReq := adapermapper.ToUpdateSprintStatusDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, updateSprintStatusReq, c.sprintUseCase.UpdateSprintStatus)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToUpdateSprintStatusGrpcResponse(resp), nil
}

func (c *SprintController) DeleteSprint(ctx context.Context, req *common.IDRequest) (*team_service.DeleteSprintResponse, error) {
	deleteSprintReq := adapermapper.ToDeleteSprintDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, deleteSprintReq, c.sprintUseCase.DeleteSprint)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToDeleteSprintGrpcResponse(resp), nil
}

func (c *SprintController) ExportSprint(ctx context.Context, req *common.IDRequest) (*team_service.ExportSprintResponse, error) {
	exportSprintReq := adapermapper.ToExportSprintDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, exportSprintReq, c.sprintUseCase.ExportSprint)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToExportSprintGrpcResponse(resp), nil
}
