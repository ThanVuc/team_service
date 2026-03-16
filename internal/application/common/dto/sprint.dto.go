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

type SprintResponse struct {
	ID      string
	GroupID string

	Name string
	Goal *string

	Status enum.SprintStatus

	StartDate time.Time
	EndDate   time.Time

	VelocityWork     int
	VelocityEstimate float64
	WorkDeleted      int

	CreatedAt time.Time
	UpdatedAt time.Time
}
