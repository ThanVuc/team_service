package entity

import (
	"strings"
	"time"

	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
)

type Group struct {
	ID          string
	Name        string
	Description *string
	OwnerID     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	AvatarURL   *string
}

func NewGroup(
	id string,
	name string,
	ownerID string,
	description *string,
	now time.Time,
) (*Group, errorbase.AppError) {

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errorbase.New(
			errdict.ErrBadRequest,
			errorbase.WithDetail("group name is required"),
		)
	}

	if len(name) < 3 || len(name) > 100 {
		return nil, errorbase.New(
			errdict.ErrBadRequest,
			errorbase.WithDetail("group name must be between 3 and 100 characters"),
		)
	}

	if description != nil {
		desc := strings.TrimSpace(*description)

		if len(desc) > 500 {
			return nil, errorbase.New(
				errdict.ErrBadRequest,
				errorbase.WithDetail("description must be at most 500 characters"),
			)
		}

		description = &desc
	}

	if ownerID == "" {
		return nil, errorbase.New(
			errdict.ErrBadRequest,
			errorbase.WithDetail("owner id is required"),
		)
	}

	return &Group{
		ID:          id,
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (g *Group) Update(
	name *string,
	description *string,
	now time.Time,
) errorbase.AppError {

	if g.DeletedAt != nil {
		return errorbase.New(errdict.ErrUnprocessable, errorbase.WithDetail("group deleted"))
	}

	if name != nil {
		n := strings.TrimSpace(*name)

		if len(n) < 3 || len(n) > 100 {
			return errorbase.New(
				errdict.ErrBadRequest,
				errorbase.WithDetail("group name must be between 3 and 100 characters"),
			)
		}

		g.Name = n
	}

	if description != nil {
		d := strings.TrimSpace(*description)

		if len(d) > 500 {
			return errorbase.New(
				errdict.ErrBadRequest,
				errorbase.WithDetail("description must be <= 500 characters"),
			)
		}

		g.Description = &d
	}

	g.UpdatedAt = now

	return nil
}

func (g *Group) Delete(now time.Time) errorbase.AppError {

	if g.DeletedAt != nil {
		return errorbase.New(errdict.ErrConflict, errorbase.WithDetail("group already deleted"))
	}

	g.DeletedAt = &now
	g.UpdatedAt = now

	return nil
}
