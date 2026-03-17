package usecase

import (
	"context"
	appdto "team_service/internal/application/common/dto"
	apphelper "team_service/internal/application/common/helper"
	istore "team_service/internal/application/common/interface/store"
	appvalidation "team_service/internal/application/common/validation"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/proto/common"

	"github.com/thanvuc/go-core-lib/log"
	"github.com/wagslane/go-rabbitmq"
)

type (
	GroupUseCase interface {
		CreateGroup(ctx context.Context, req *appdto.CreateGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError)
		Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, errorbase.AppError)
		GetGroupRequest(ctx context.Context, req *appdto.GetGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError)
		UpdateGroup(ctx context.Context, req *appdto.UpdateGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError)
		DeleteGroup(ctx context.Context, req *appdto.DeleteGroupRequest) (*appdto.BaseResponse[appdto.DeleteGroupResponse], errorbase.AppError)
	}

	SprintUseCase interface {
		CreateSprint(ctx context.Context, req *appdto.CreateSprintRequest) (*appdto.BaseResponse[appdto.SprintResponse], errorbase.AppError)
		GetSprint(ctx context.Context, req *appdto.GetSprintRequest) (*appdto.BaseResponse[appdto.SprintResponse], errorbase.AppError)
		ListSprints(ctx context.Context, req *appdto.ListSprintsRequest) (*appdto.BaseResponse[appdto.ListSprintsResponse], errorbase.AppError)
		UpdateSprint(ctx context.Context, req *appdto.UpdateSprintRequest) (*appdto.BaseResponse[appdto.SprintResponse], errorbase.AppError)
		UpdateSprintStatus(ctx context.Context, req *appdto.UpdateSprintStatusRequest) (*appdto.BaseResponse[appdto.UpdateSprintStatusResponse], errorbase.AppError)
		DeleteSprint(ctx context.Context, req *appdto.DeleteSprintRequest) (*appdto.BaseResponse[appdto.DeleteSprintResponse], errorbase.AppError)
	}

	UserUseCase interface {
		SyncUserData(ctx context.Context) func(d rabbitmq.Delivery) rabbitmq.Action
	}
)

func NewGroupUseCase(
	store istore.Store,
	validator *appvalidation.GroupValidator,
	authHelper *apphelper.AuthHelper,
) GroupUseCase {
	return &groupUseCase{
		store:     store,
		groupRepo: store.GroupRepository(),
		userRepo:  store.UserRepository(),
		validator: validator,
	}
}

func NewSprintUseCase(
	store istore.Store,
	validator *appvalidation.SprintValidator,
	authHelper *apphelper.AuthHelper,
) SprintUseCase {
	return &sprintUseCase{
		store:      store,
		sprintRepo: store.SprintRepository(),
		validator:  validator,
		authHelper: authHelper,
	}
}

func NewUserUseCase(
	store istore.Store,
	logger log.LoggerV2,
) UserUseCase {
	return &userUseCase{
		store:  store,
		logger: logger,
	}
}
