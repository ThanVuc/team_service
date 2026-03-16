package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	appdto "team_service/internal/application/common/dto"
	istore "team_service/internal/application/common/interface/store"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"team_service/proto/common"
	"time"

	"github.com/thanvuc/go-core-lib/log"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type userUseCase struct {
	store  istore.Store
	logger log.LoggerV2
}

func (uc *userUseCase) SyncUserData(ctx context.Context) func(d rabbitmq.Delivery) rabbitmq.Action {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		user := &entity.User{}

		outbox := &common.Outbox{}
		err := proto.Unmarshal(d.Body, outbox)
		if err != nil {
			uc.logger.Error("failed to unmarshal outbox message", log.WithFields(zap.Error(err)))
			return rabbitmq.NackDiscard
		}

		var userPayload appdto.UserOutboxPayload
		if err := json.Unmarshal(outbox.Payload, &userPayload); err != nil {
			uc.logger.Error("failed to unmarshal user payload", log.WithFields(zap.Error(err)))
			return rabbitmq.NackDiscard
		}

		user, err = entity.CreateUser(
			userPayload.UserID,
			userPayload.Email,
			time.UnixMilli(userPayload.CreatedAt),
			enum.UserStatusActive,
			&userPayload.AvatarUrl,
		)

		uc.logger.Info(fmt.Sprintf(
			"creating user: user_id=%s email=%s created_at_raw=%d created_at=%s avatar_url=%s",
			user.ID,
			user.Email,
			user.CreatedAt,
			user.CreatedAt.Format(time.RFC3339),
			*user.AvatarURL,
		))

		err = uc.store.ExecTx(ctx, func(repo istore.RepositoryContainer) errorbase.AppError {
			err := repo.UserRepository().UpsertUser(ctx, user)
			if err != nil {
				return errorbase.Wrap(err, errdict.ErrUnprocessable, errorbase.WithDetail("failed to upsert user in database"))
			}
			return nil
		})

		if err != nil {
			uc.logger.Error("failed to upsert user in database", log.WithFields(zap.Error(err)))
			return rabbitmq.NackDiscard
		}

		return rabbitmq.Ack
	}
}
