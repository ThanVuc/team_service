package usecase

import (
	"context"
	istore "team_service/internal/application/common/interface/store"
	"team_service/proto/common"
)

type (
	GroupUseCase interface {
		CreateGroup(ctx context.Context) error
		Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, error)
	}
)

func NewGroupUseCase(store istore.Store) GroupUseCase {
	return &groupUseCase{
		store: store,
	}
}
