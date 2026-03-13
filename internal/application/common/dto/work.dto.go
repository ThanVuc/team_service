package appdto

import (
	"team_service/internal/domain/enum"
	"time"
)

type CreateWorkRequest struct {
	GroupID string

	SprintID *string

	Name        string
	Description *string

	Status enum.WorkStatus

	Priority enum.WorkPriority

	AssigneeID *string

	EstimateHours *float64
	StoryPoint    *int

	DueDate *time.Time
}

type UpdateWorkRequest struct {
	WorkID string

	Name        *string
	Description *string

	Status *enum.WorkStatus

	SprintID *string

	AssigneeID *string

	EstimateHours *float64
	StoryPoint    *int

	Priority *enum.WorkPriority

	DueDate *time.Time

	Version int
}

type DeleteWorkRequest struct {
	WorkID string
}

type WorkResponse struct {
	ID       string
	GroupID  string
	SprintID *string

	Name        string
	Description *string

	Status enum.WorkStatus

	Priority enum.WorkPriority

	AssigneeID *string
	CreatorID  string

	EstimateHours *float64
	StoryPoint    *int

	DueDate *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
