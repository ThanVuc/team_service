package repository

import (
	"context"
	"database/sql"
	appdto "team_service/internal/application/common/dto"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/internal/infrastructure/share/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	q *database.Queries
}

func NewUserRepository(
	q *database.Queries,
) *UserRepository {
	return &UserRepository{
		q: q,
	}
}

func (r *UserRepository) GetUserByID(
	ctx context.Context,
	userID string,
) (*entity.User, errorbase.AppError) {
	var uuid pgtype.UUID
	if err := uuid.Scan(userID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse user id"),
		)
	}

	u, err := r.q.GetUserByID(ctx, uuid)

	if err != nil {
		return nil, errorbase.Wrap(err, errdict.ErrBadRequest)
	}

	return &entity.User{
		ID:        u.ID.String(),
		Email:     u.Email,
		Status:    enum.UserStatus(u.Status),
		CreatedAt: u.CreatedAt.Time,
		AvatarURL: utils.Ptr(u.AvatarUrl.String),
	}, nil
}

func (r *UserRepository) GetUserWithPermissionByID(
	ctx context.Context,
	groupId string,
	userID string,
) (*appdto.UserWithPermission, errorbase.AppError) {
	userUUID, err := utils.ToUUID(userID)
	if err != nil || !userUUID.Valid {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse user id to UUID"),
		)
	}

	groupUUID, err := utils.ToUUID(groupId)
	if err != nil || !groupUUID.Valid {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id to UUID"),
		)
	}

	user, err := r.q.GetUserWithPermissionByID(
		ctx,
		database.GetUserWithPermissionByIDParams{
			GroupID: groupUUID,
			ID:      userUUID,
		},
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorbase.New(
				errdict.ErrUnauthorized,
				errorbase.WithDetail("user is unauthorized to access this group or user does not exist"),
			)
		}

		return nil, errorbase.Wrap(
			err,
			errdict.ErrBadRequest,
			errorbase.WithDetail("failed to get user with permission by id in repository"),
		)
	}

	return &appdto.UserWithPermission{
		ID:       user.ID.String(),
		Email:    user.Email,
		Status:   enum.UserStatus(user.Status),
		Role:     enum.GroupRole(user.Role),
		GroupId:  groupId,
		JoinedAt: user.JoinedAt.Time,
	}, nil
}

func (r *UserRepository) UpsertUser(
	ctx context.Context,
	user *entity.User,
) errorbase.AppError {
	userUUID, err := utils.ToUUID(user.ID)
	if err != nil || !userUUID.Valid {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse user id to UUID"),
		)
	}

	params := database.UpsertUserParams{
		ID:                   userUUID,
		Email:                user.Email,
		Status:               string(user.Status),
		CreatedAt:            pgtype.Timestamptz{Time: user.CreatedAt, Valid: true},
		HasEmailNotification: user.HasEmailNotification,
		HasPushNotification:  user.HasPushNotification,
	}

	if user.AvatarURL != nil {
		params.AvatarUrl = pgtype.Text{
			String: *user.AvatarURL,
			Valid:  true,
		}
	}

	err = r.q.UpsertUser(ctx, params)
	if err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail(err.Error()),
		)
	}

	return nil
}

func (r *UserRepository) GetListMembersByGroupID(ctx context.Context, groupID string) (*appdto.ListMembersResponse, errorbase.AppError) {
	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	members, err := r.q.GetListGroupMembers(ctx, groupUUID)
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail("failed to get list members by group id in repository"),
		)
	}

	if len(members) == 0 {
		return &appdto.ListMembersResponse{
			Members: []appdto.MemberResponse{},
			Total:   0,
		}, nil
	}

	memberVals := make([]appdto.MemberResponse, len(members))
	for i, member := range members {
		memberVals[i] = appdto.MemberResponse{
			ID:       member.ID.String(),
			Email:    member.Email,
			Role:     enum.GroupRole(member.Role),
			Avatar:   utils.Ptr(member.Avatar.String),
			JoinedAt: member.JoinedAt.Time,
		}
	}

	return &appdto.ListMembersResponse{
		Members: memberVals,
		Total:   int(len(memberVals)),
	}, nil
}
