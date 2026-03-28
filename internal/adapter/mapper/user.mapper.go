package adapermapper

import (
	appdto "team_service/internal/application/common/dto"
	"team_service/proto/common"
	"team_service/proto/team_service"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToGetUserInfoDTO(req *common.EmptyRequest) *appdto.UserInfoRequest {
	return &appdto.UserInfoRequest{}
}

func ToGetUserInfoGrpcResponse(
	resp *appdto.BaseResponse[appdto.UserInfoResponse],
) *team_service.GetUserInfoResponse {
	if resp == nil {
		return &team_service.GetUserInfoResponse{
			Email:                "",
			UseEmailNotification: false,
			UseAppNotification:   false,
			CreatedAt:            nil,
			Error:                ToProtoError(nil),
		}
	}

	return &team_service.GetUserInfoResponse{
		Email:                resp.Data.Email,
		UseEmailNotification: resp.Data.HasEmailNotification,
		UseAppNotification:   resp.Data.HasPushNotification,
		CreatedAt:            timestamppb.New(time.Unix(resp.Data.CreatedAt, 0)),
		Error:                ToProtoError(resp.Error),
	}
}

func ToConfigureNotificationDTO(req *team_service.NotificationConfigurationRequest) *appdto.ConfigureNotificationRequest {
	if req == nil {
		return nil
	}

	return &appdto.ConfigureNotificationRequest{
		UseEmailNotification: req.UseEmailNotification,
		UseAppNotification:   req.UseAppNotification,
	}

}

func ToConfigureNotificationGrpcResponse(
	resp *appdto.BaseResponse[appdto.ConfigureNotificationResponse],
) *team_service.NotificationConfigurationResponse {
	if resp == nil {
		return &team_service.NotificationConfigurationResponse{
			Success: false,
			Error:   ToProtoError(nil),
		}
	}

	return &team_service.NotificationConfigurationResponse{
		Success: resp.Data.Success,
		Error:   ToProtoError(resp.Error),
	}
}
