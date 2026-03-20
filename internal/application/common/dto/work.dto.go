package appdto

import (
	"team_service/internal/domain/enum"
	"time"
)

type CreateWorkRequest struct {
	SprintID    *string
	Name        string
	Description *string
}

type GetWorkRequest struct {
	WorkID string
}

type ListWorksRequest struct {
	SprintID *string
}

type UpdateWorkRequest struct {
	WorkID string

	Name        *string
	Description *string

	Status *enum.WorkStatus

	SprintID *string

	AssigneeID *string

	StoryPoint *int32

	Priority *enum.WorkPriority

	DueDate *time.Time

	Version int32
}

type DeleteWorkRequest struct {
	WorkID string
}

type DeleteWorkResponse struct {
	Success bool
}

type ListWorksResponse struct {
	Works []WorkResponse
}

type SimpleSprintDTO struct {
	ID   string
	Name string
}

type SimpleUserDTO struct {
	ID     string
	Email  string
	Avatar *string
}

type WorkResponse struct {
	ID       string
	GroupID  string
	SprintID *string

	Name        string
	Description *string

	Status enum.WorkStatus

	Priority enum.WorkPriority
	Sprint   *SimpleSprintDTO
	Assignee *SimpleUserDTO

	AssigneeID *string
	CreatorID  string

	EstimateHours *float64
	StoryPoint    *int32

	DueDate *time.Time

	CheckList *ChecklistSummaryResponse
	Comments  *CommentListResponse

	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int32
}
