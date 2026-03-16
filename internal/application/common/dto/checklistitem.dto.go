package appdto

import "time"

type CreateChecklistItemRequest struct {
	WorkID string
	Name   string
}

type UpdateChecklistItemRequest struct {
	ItemID      string
	Name        *string
	IsCompleted *bool
	Version     int
}

type DeleteChecklistItemRequest struct {
	ItemID string
}

type ChecklistItemResponse struct {
	ID          string
	WorkID      string
	Name        string
	IsCompleted bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
