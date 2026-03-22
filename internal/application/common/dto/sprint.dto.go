package appdto

import (
	"team_service/internal/domain/enum"
	"time"
)

type CreateSprintRequest struct {
	GroupID string
	Name    string
	Goal    *string

	StartDate time.Time
	EndDate   time.Time
}

type GetSprintRequest struct {
	SprintID string
}

type ListSprintsRequest struct {
	GroupID string
}

type UpdateSprintRequest struct {
	SprintID string

	Name *string
	Goal *string

	StartDate *time.Time
	EndDate   *time.Time
}

type UpdateSprintStatusRequest struct {
	SprintID string
	Status   enum.SprintStatus
}

type DeleteSprintRequest struct {
	SprintID string
}

type ListSprintsResponse struct {
	Sprints []SprintResponse
	Total   int32
}

type UpdateSprintStatusResponse struct {
	SprintID string
	Status   enum.SprintStatus
}

type DeleteSprintResponse struct {
	Success bool
}

type SprintResponse struct {
	ID      string
	GroupID string

	Name string
	Goal *string

	Status enum.SprintStatus

	StartDate time.Time
	EndDate   time.Time

	TotalWork       int32
	CompletedWork   int32
	ProgressPercent float32

	CreatedAt time.Time
	UpdatedAt time.Time
}

type SimpleSprintResponse struct {
	ID     string
	Name   string
	Status enum.SprintStatus
}
