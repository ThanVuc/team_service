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

type UserController struct {
	team_service.UnimplementedUserServiceServer
	userUseCase usecase.UserUseCase
	logger      log.LoggerV2
}

func NewUserController(
	userUseCase usecase.UserUseCase,
	logger log.LoggerV2,
) *UserController {
	return &UserController{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

func (c *UserController) GetUserInfo(ctx context.Context, req *common.EmptyRequest) (*team_service.GetUserInfoResponse, error) {
	getUserInfoReq := adapermapper.ToGetUserInfoDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, getUserInfoReq, c.userUseCase.GetUserInfo)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToGetUserInfoGrpcResponse(resp), nil
}

func (c *UserController) NotificationConfiguration(ctx context.Context, req *team_service.NotificationConfigurationRequest) (*team_service.NotificationConfigurationResponse, error) {
	configureNotificationReq := adapermapper.ToConfigureNotificationDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, configureNotificationReq, c.userUseCase.NotificationConfiguration)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToConfigureNotificationGrpcResponse(resp), nil
}
