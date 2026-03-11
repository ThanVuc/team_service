package usecase

import (
	"context"
	istore "team_service/internal/application/common/interface/store"
	appmapper "team_service/internal/application/common/mapper"
	"team_service/proto/common"

	"github.com/wagslane/go-rabbitmq"
	"team_service/proto/team_service"
)

type (
	GroupUseCase interface {
		CreateGroup(ctx context.Context, req *team_service.CreateGroupRequest) (*team_service.CreateGroupResponse, error)
		Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, error)
	}

	UserUseCase interface {
		SyncUserData(ctx context.Context) func(d rabbitmq.Delivery) rabbitmq.Action
	}
)

func NewGroupUseCase(
	store istore.Store,
	mapper *appmapper.GroupMapper,
) GroupUseCase {
	return &groupUseCase{
		store:  store,
		mapper: mapper,
	}
}

func NewUserUseCase(store istore.Store) UserUseCase {
	return &userUseCase{
		store: store,
	}
}
