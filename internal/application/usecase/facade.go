package usecase

import (
	"context"
	appdto "team_service/internal/application/common/dto"
	istore "team_service/internal/application/common/interface/store"
	appvalidation "team_service/internal/application/common/validation"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/proto/common"

	"github.com/wagslane/go-rabbitmq"
)

type (
	GroupUseCase interface {
		CreateGroup(ctx context.Context, req *appdto.CreateGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError)
		Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, errorbase.AppError)
		GetGroupRequest(ctx context.Context, req *appdto.GetGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError)
	}

	UserUseCase interface {
		SyncUserData(ctx context.Context) func(d rabbitmq.Delivery) rabbitmq.Action
	}
)

func NewGroupUseCase(
	store istore.Store,
	validator *appvalidation.GroupValidator,
) GroupUseCase {
	return &groupUseCase{
		store:     store,
		groupRepo: store.GroupRepository(),
		userRepo:  store.UserRepository(),
		validator: validator,
	}
}

func NewUserUseCase(store istore.Store) UserUseCase {
	return &userUseCase{
		store: store,
	}
}
