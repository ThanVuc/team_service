package entity

import (
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"time"
)

type Comment struct {
	ID        string
	WorkID    string
	CreatorID string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Comment) UpdateContent(
	userID string,
	content string,
	now time.Time,
) errorbase.AppError {

	if c.CreatorID != userID {
		return errorbase.New(errdict.ErrForbidden)
	}

	c.Content = content
	c.UpdatedAt = now

	return nil
}

func NewComment(
	id string,
	workID string,
	creatorID string,
	content string,
	now time.Time,
) (*Comment, errorbase.AppError) {
	if content == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("content is required"))
	}

	return &Comment{
		ID:        id,
		WorkID:    workID,
		CreatorID: creatorID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
