package usecase

import (
	"context"
	istore "team_service/internal/application/common/interface/store"
	appmapper "team_service/internal/application/common/mapper"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/proto/common"

	"team_service/proto/team_service"

	"github.com/wagslane/go-rabbitmq"
)

type (
	GroupUseCase interface {
		CreateGroup(ctx context.Context, req *team_service.CreateGroupRequest) (*team_service.CreateGroupResponse, errorbase.AppError)
		Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, errorbase.AppError)
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
		store:     store,
		mapper:    mapper,
		groupRepo: store.GroupRepository(),
	}
}

func NewUserUseCase(store istore.Store) UserUseCase {
	return &userUseCase{
		store: store,
	}
}
