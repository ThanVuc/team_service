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

type WorkController struct {
	team_service.UnimplementedWorkServiceServer
	workUseCase usecase.WorkUseCase
	logger      log.LoggerV2
}

func NewWorkController(
	workUseCase usecase.WorkUseCase,
	logger log.LoggerV2,
) *WorkController {
	return &WorkController{
		workUseCase: workUseCase,
		logger:      logger,
	}
}

func (c *WorkController) CreateWork(ctx context.Context, req *team_service.CreateWorkRequest) (*team_service.CreateWorkResponse, error) {
	createWorkReq := adapermapper.ToCreateWorkDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, createWorkReq, c.workUseCase.CreateWork)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToCreateWorkGrpcResponse(resp), nil
}

func (c *WorkController) GetWork(ctx context.Context, req *common.IDRequest) (*team_service.GetWorkResponse, error) {
	getWorkReq := adapermapper.ToGetWorkDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, getWorkReq, c.workUseCase.GetWork)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToGetWorkGrpcResponse(resp), nil
}

func (c *WorkController) ListWorks(ctx context.Context, req *team_service.ListWorksRequest) (*team_service.ListWorksResponse, error) {
	listWorksReq := adapermapper.ToListWorksDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, listWorksReq, c.workUseCase.ListWorks)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToListWorksGrpcResponse(resp), nil
}

func (c *WorkController) UpdateWork(ctx context.Context, req *team_service.UpdateWorkRequest) (*team_service.UpdateWorkResponse, error) {
	updateWorkReq := adapermapper.ToUpdateWorkDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, updateWorkReq, c.workUseCase.UpdateWork)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToUpdateWorkGrpcResponse(resp), nil
}

func (c *WorkController) DeleteWork(ctx context.Context, req *common.IDRequest) (*team_service.DeleteWorkResponse, error) {
	deleteWorkReq := adapermapper.ToDeleteWorkDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, deleteWorkReq, c.workUseCase.DeleteWork)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToDeleteWorkGrpcResponse(resp), nil
}

func (c *WorkController) CreateChecklistItem(ctx context.Context, req *team_service.CreateChecklistItemRequest) (*team_service.CreateChecklistItemResponse, error) {
	createChecklistReq := adapermapper.ToCreateChecklistItemDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, createChecklistReq, c.workUseCase.CreateChecklistItem)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToCreateChecklistItemGrpcResponse(resp), nil
}

func (c *WorkController) UpdateChecklistItem(ctx context.Context, req *team_service.UpdateChecklistItemRequest) (*team_service.UpdateChecklistItemResponse, error) {
	updateChecklistReq := adapermapper.ToUpdateChecklistItemDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, updateChecklistReq, c.workUseCase.UpdateChecklistItem)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToUpdateChecklistItemGrpcResponse(resp), nil
}

func (c *WorkController) DeleteChecklistItem(ctx context.Context, req *common.IDRequest) (*team_service.DeleteChecklistItemResponse, error) {
	deleteChecklistReq := adapermapper.ToDeleteChecklistItemDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, deleteChecklistReq, c.workUseCase.DeleteChecklistItem)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToDeleteChecklistItemGrpcResponse(resp), nil
}

func (c *WorkController) CreateComment(ctx context.Context, req *team_service.CreateCommentRequest) (*team_service.CreateCommentResponse, error) {
	createCommentReq := adapermapper.ToCreateCommentDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, createCommentReq, c.workUseCase.CreateComment)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToCreateCommentGrpcResponse(resp), nil
}

func (c *WorkController) UpdateComment(ctx context.Context, req *team_service.UpdateCommentRequest) (*team_service.UpdateCommentResponse, error) {
	updateCommentReq := adapermapper.ToUpdateCommentDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, updateCommentReq, c.workUseCase.UpdateComment)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToUpdateCommentGrpcResponse(resp), nil
}

func (c *WorkController) DeleteComment(ctx context.Context, req *common.IDRequest) (*team_service.DeleteCommentResponse, error) {
	deleteCommentReq := adapermapper.ToDeleteCommentDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, deleteCommentReq, c.workUseCase.DeleteComment)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToDeleteCommentGrpcResponse(resp), nil
}
