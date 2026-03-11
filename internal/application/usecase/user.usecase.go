package usecase

import (
	"context"
	istore "team_service/internal/application/common/interface/store"

	"github.com/wagslane/go-rabbitmq"
)

type userUseCase struct {
	store istore.Store
}

func (uc *userUseCase) SyncUserData(ctx context.Context) func(d rabbitmq.Delivery) rabbitmq.Action {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		return rabbitmq.Ack
	}
}
