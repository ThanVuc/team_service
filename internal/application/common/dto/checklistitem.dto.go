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

type ChecklistItemMeta struct {
	ID     string
	WorkID string
}

type ChecklistSummaryResponse struct {
	Total     int32
	Completed int32
	Items     []ChecklistItemResponse
}
