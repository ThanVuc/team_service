package usecase

import (
	"context"
	istore "team_service/internal/application/common/interface/store"
	"team_service/proto/common"

	"github.com/wagslane/go-rabbitmq"
)

type (
	GroupUseCase interface {
		CreateGroup(ctx context.Context) error
		Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, error)
	}

	UserUseCase interface {
		SyncUserData(ctx context.Context) func(d rabbitmq.Delivery) rabbitmq.Action
	}
)

func NewGroupUseCase(store istore.Store) GroupUseCase {
	return &groupUseCase{
		store: store,
	}
}

func NewUserUseCase(store istore.Store) UserUseCase {
	return &userUseCase{
		store: store,
	}
}
