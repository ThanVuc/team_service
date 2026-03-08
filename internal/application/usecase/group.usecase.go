package usecase

import (
	"context"
	istore "team_service/internal/application/common/interface/store"
	"team_service/internal/infrastructure/share/utils"
	"team_service/proto/common"
)

type groupUseCase struct {
	// Add any dependencies or services needed for the Group use case here
	// Example: GroupRepository
	store istore.Store
}

func (uc *groupUseCase) CreateGroup(ctx context.Context) error {
	// Implement the logic to create a group using uc.groupRepo
	return uc.store.ExecTx(ctx, func(store istore.RepositoryContainer) error {
		// Example: Call the CreateGroup method of the GroupRepository
		return store.GroupRepository().CreateGroup()
	})
}

func (uc *groupUseCase) Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, error) {
	// Implement the logic for the Ping method
	userID := utils.GetUserIDFromOutgoingContext(ctx)
	println("Received Ping request from user ID:", userID)
	return &common.EmptyResponse{}, nil
}
