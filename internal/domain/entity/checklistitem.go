package entity

import (
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"time"
)

type ChecklistItem struct {
	ID          string
	WorkID      string
	Name        string
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewChecklistItem(
	id string,
	workID string,
	name string,
	now time.Time,
) (*ChecklistItem, errorbase.AppError) {
	if name == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("name is required"))
	}

	return &ChecklistItem{
		ID:          id,
		WorkID:      workID,
		Name:        name,
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (c *ChecklistItem) ToggleComplete(now time.Time) {
	c.IsCompleted = !c.IsCompleted
	c.UpdatedAt = now
}
