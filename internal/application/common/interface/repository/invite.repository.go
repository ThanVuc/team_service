package irepository

import (
	"context"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/internal/domain/entity"
)

type InviteRepository interface {
	CreateInvite(ctx context.Context, invite *entity.Invite) (*entity.Invite, errorbase.AppError)
}
