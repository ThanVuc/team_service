package entity

import (
	"strings"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/enum"
	"time"
)

type Work struct {
	ID            string
	GroupID       string
	SprintID      string
	Name          string
	Description   *string
	Status        enum.WorkStatus
	AssigneeID    string
	CreatorID     string
	EstimateHours *float64
	StoryPoint    *int32
	Priority      *enum.WorkPriority
	DueDate       *time.Time
	CompletedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewWork(
	id string,
	groupID string,
	sprintID string,
	name string,
	description *string,
	creatorID string,
	assigneeID string,
	estimateHours *float64,
	storyPoint *int32,
	priority *enum.WorkPriority,
	dueDate *time.Time,
	now time.Time,
) (*Work, errorbase.AppError) {

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errorbase.New(
			errdict.ErrBadRequest,
			errorbase.WithDetail("work name is required"),
		)
	}

	if groupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	if creatorID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("creator id is required"))
	}

	if priority != nil && !priority.IsValid() {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("invalid priority"))
	}

	return &Work{
		ID:            id,
		GroupID:       groupID,
		SprintID:      sprintID,
		Name:          name,
		Description:   description,
		Status:        enum.WorkStatusTodo,
		AssigneeID:    assigneeID,
		CreatorID:     creatorID,
		EstimateHours: estimateHours,
		StoryPoint:    storyPoint,
		Priority:      priority,
		DueDate:       dueDate,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (w *Work) Rename(name string, now time.Time) errorbase.AppError {
	name = strings.TrimSpace(name)

	if name == "" {
		return errorbase.New(
			errdict.ErrBadRequest,
			errorbase.WithDetail("name cannot be empty"),
		)
	}

	w.Name = name
	w.UpdatedAt = now

	return nil
}

func (w *Work) UpdateDescription(desc *string, now time.Time) {
	w.Description = desc
	w.UpdatedAt = now
}

func (w *Work) Assign(userID string, now time.Time) errorbase.AppError {

	userID = strings.TrimSpace(userID)

	if userID == "" {
		return errorbase.New(
			errdict.ErrBadRequest,
			errorbase.WithDetail("assignee id is required"),
		)
	}

	w.AssigneeID = userID
	w.UpdatedAt = now

	return nil
}

func (w *Work) SetEstimateHours(hours float64, now time.Time) errorbase.AppError {

	if hours < 0 {
		return errorbase.New(
			errdict.ErrBadRequest,
			errorbase.WithDetail("estimate hours must be positive"),
		)
	}

	w.EstimateHours = &hours
	w.UpdatedAt = now

	return nil
}

func (w *Work) SetStoryPoint(point int32, now time.Time) errorbase.AppError {

	if point < 0 {
		return errorbase.New(
			errdict.ErrBadRequest,
			errorbase.WithDetail("story point must be positive"),
		)
	}

	w.StoryPoint = &point
	w.UpdatedAt = now

	return nil
}

func (w *Work) SetPriority(priority enum.WorkPriority, now time.Time) errorbase.AppError {

	if !priority.IsValid() {
		return errorbase.New(
			errdict.ErrBadRequest,
			errorbase.WithDetail("invalid priority"),
		)
	}

	w.Priority = &priority
	w.UpdatedAt = now

	return nil
}

func (w *Work) ChangeStatus(status enum.WorkStatus, now time.Time) errorbase.AppError {

	if !status.IsValid() {
		return errorbase.New(
			errdict.ErrBadRequest,
			errorbase.WithDetail("invalid status"),
		)
	}

	if status == enum.WorkStatusDone && w.CompletedAt == nil {
		w.CompletedAt = &now
	}

	w.Status = status
	w.UpdatedAt = now

	return nil
}

func (w *Work) MoveToSprint(sprintID string, now time.Time) errorbase.AppError {

	sprintID = strings.TrimSpace(sprintID)

	if sprintID == "" {
		return errorbase.New(
			errdict.ErrBadRequest,
			errorbase.WithDetail("sprint id is required"),
		)
	}

	w.SprintID = sprintID
	w.UpdatedAt = now

	return nil
}

func (w *Work) SetDueDate(dueDate *time.Time, now time.Time) {
	w.DueDate = dueDate
	w.UpdatedAt = now
}
