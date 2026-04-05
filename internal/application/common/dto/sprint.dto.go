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

type GenerateSprintRequest struct {
	Name              string                   `json:"name"`
	Goal              string                   `json:"goal,omitempty"`
	StartDate         string                   `json:"start_date"`
	EndDate           string                   `json:"end_date"`
	AdditionalContext *string                  `json:"additional_context,omitempty"`
	Files             []AISprintGenerationFile `json:"files,omitempty"`
}

type AISprintGenerationFile struct {
	ObjectKey string `json:"object_key"`
	Size      int64  `json:"size"`
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

type GenerateSprintResponse struct {
	Message string
}

type AISprintGenerationRequestedMessage struct {
	EventType string                             `json:"event_type"`
	JobID     string                             `json:"job_id"`
	GroupID   string                             `json:"group_id"`
	SenderID  string                             `json:"sender_id"`
	Payload   AISprintGenerationRequestedPayload `json:"payload"`
}

type AISprintGenerationResultMessage struct {
	EventType string                          `json:"event_type"`
	JobID     string                          `json:"job_id"`
	GroupID   string                          `json:"group_id"`
	SenderID  string                          `json:"sender_id"`
	Payload   AISprintGenerationResultPayload `json:"payload"`
}

type AISprintGenerationRequestedPayload struct {
	Sprint            AISprintGenerationSprint `json:"sprint"`
	Files             []AISprintGenerationFile `json:"files"`
	AdditionalContext *string                  `json:"additional_context,omitempty"`
}

type AISprintGenerationResultPayload struct {
	Status string                   `json:"status"`
	Sprint AISprintGenerationSprint `json:"sprint"`
	Tasks  []AISprintGeneratedTask  `json:"tasks"`
	Error  *AISprintGenerationError `json:"error,omitempty"`
}

type AISprintGenerationSprint struct {
	Name      string `json:"name"`
	Goal      string `json:"goal,omitempty"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type AISprintGeneratedTask struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Priority    *string `json:"priority,omitempty"`
	StoryPoint  *int    `json:"story_point,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
}

type AISprintGenerationError struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Detail  *string `json:"detail,omitempty"`
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
