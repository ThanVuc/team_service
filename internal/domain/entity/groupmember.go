package entity

import (
	"strings"
	"time"

	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/enum"
)

type GroupMember struct {
	ID       string
	GroupID  string
	UserID   string
	Role     enum.GroupRole
	JoinedAt time.Time
}

func NewGroupMember(
	id string,
	groupID string,
	userID string,
	role enum.GroupRole,
	now time.Time,
) (*GroupMember, errorbase.AppError) {
	groupID = strings.TrimSpace(groupID)
	userID = strings.TrimSpace(userID)

	if groupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	if userID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("user id is required"))
	}

	if !role.IsValid() {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("invalid role"))
	}

	return &GroupMember{
		ID:       id,
		GroupID:  groupID,
		UserID:   userID,
		Role:     role,
		JoinedAt: now,
	}, nil
}

func (m *GroupMember) UpdateRole(role enum.GroupRole) errorbase.AppError {

	if !role.IsValid() {
		return errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("invalid role"))
	}

	if role == enum.GroupRoleOwner {
		return errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("cannot assign owner"))
	}

	m.Role = role

	return nil
}
