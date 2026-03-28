package appdto

import (
	appconstant "team_service/internal/application/common/constant"
	"team_service/internal/domain/entity"
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

type ExportSprintRequest struct {
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

type ExportSprintResponse struct {
	FileName    string
	File        []byte
	ContentType string
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

type SprintExportInput struct {
	Sprint   *entity.Sprint
	Members  []*entity.User
	Works    []*entity.Work
	FileName string
}

type SprintExportOutput struct {
	FileName string
	Content  []byte
}

type SprintTaskView struct {
	Name         string
	AssigneeName string
	StoryPoint   *int
	CompletedAt  *time.Time
	Status       appconstant.SprintTaskStatus
}

type SprintStatistic struct {
	TotalEstimated   int
	TotalCompleted   int
	CompletionRate   float64
	Spillover        int
	UnestimatedCount int
}

type SprintExportStyles struct {
	Base       int
	Bold       int
	BoldBorder int
	Header     int
	TotalLabel int
	TotalData  int
	DoneEarly  int
	DoneOnTime int
	DoneLate   int
}

type SprintExportCellColor struct {
	Cell   string
	Status appconstant.SprintTaskStatus
}
