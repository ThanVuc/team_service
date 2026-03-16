package entity

import (
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/enum"
	"time"
)

type Sprint struct {
	ID               string
	GroupID          string
	Name             string
	Goal             *string
	StartDate        time.Time
	EndDate          time.Time
	Status           enum.SprintStatus
	VelocityWork     *int32
	VelocityEstimate *float64
	WorkDeleted      *int32
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewSprint(
	id string,
	groupID string,
	name string,
	start time.Time,
	end time.Time,
	now time.Time,
) (*Sprint, errorbase.AppError) {

	if name == "" {
		return nil, errorbase.New(errdict.ErrBadRequest)
	}

	if !start.Before(end) {
		return nil, errorbase.New(errdict.ErrBadRequest)
	}

	if end.Sub(start).Hours() > 24*30 {
		return nil, errorbase.New(errdict.ErrBadRequest)
	}

	if start.Before(now) {
		return nil, errorbase.New(errdict.ErrBadRequest)
	}

	if groupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest)
	}

	if id == "" {
		return nil, errorbase.New(errdict.ErrBadRequest)
	}

	return &Sprint{
		ID:               id,
		GroupID:          groupID,
		Name:             name,
		StartDate:        start,
		EndDate:          end,
		Status:           enum.SprintStatusDraft,
		CreatedAt:        now,
		UpdatedAt:        now,
		Goal:             new(string),
		VelocityWork:     new(int32),
		VelocityEstimate: new(float64),
		WorkDeleted:      new(int32),
	}, nil
}

func (s *Sprint) Update(
	name string,
	goal *string,
	start time.Time,
	end time.Time,
	now time.Time,
) errorbase.AppError {

	if s.Status != "draft" {
		return errorbase.New(errdict.ErrUnprocessable, errorbase.WithDetail("only draft editable"))
	}

	s.Name = name
	s.Goal = goal
	s.StartDate = start
	s.EndDate = end
	s.UpdatedAt = now

	return nil
}

func (s *Sprint) ChangeStatus(status enum.SprintStatus) errorbase.AppError {

	if !status.IsValid() {
		return errorbase.New(errdict.ErrBadRequest)
	}

	if s.Status == status {
		return errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("status is the same"))
	}

	switch s.Status {
	case enum.SprintStatusDraft:
		if status != enum.SprintStatusActive && status != enum.SprintStatusCancelled {
			return errorbase.New(errdict.ErrUnprocessable)
		}
	case enum.SprintStatusActive:
		if status != enum.SprintStatusCompleted && status != enum.SprintStatusCancelled {
			return errorbase.New(errdict.ErrUnprocessable)
		}
	default:
		return errorbase.New(errdict.StatusTransitionInvalid)
	}

	s.Status = status

	return nil
}
