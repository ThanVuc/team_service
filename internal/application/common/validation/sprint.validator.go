package appvalidation

import (
	"context"
	"strings"
	appdto "team_service/internal/application/common/dto"
	irepository "team_service/internal/application/common/interface/repository"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/infrastructure/share/utils"
	"time"

	"github.com/google/uuid"
)

type SprintValidator struct {
	sprintRepo irepository.SprintRepository
	userRepo   irepository.UserRepository
	groupRepo  irepository.GroupRepository
}

func NewSprintValidator(
	sprintRepo irepository.SprintRepository,
	userRepo irepository.UserRepository,
	groupRepo irepository.GroupRepository,
) *SprintValidator {
	return &SprintValidator{
		sprintRepo: sprintRepo,
		userRepo:   userRepo,
		groupRepo:  groupRepo,
	}
}

func (v *SprintValidator) ValidateCreateSprint(
	ctx context.Context,
	req *appdto.CreateSprintRequest,
) (*entity.Sprint, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	groupID := strings.TrimSpace(req.GroupID)
	if groupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	if _, err := uuid.Parse(groupID); err != nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id must be a valid UUID"))
	}

	userID := strings.TrimSpace(utils.GetUserIDFromOutgoingContext(ctx))
	if userID == "" {
		return nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("missing user context"))
	}

	if _, err := uuid.Parse(userID); err != nil {
		return nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("user id in context is invalid"))
	}

	user, err := v.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("user not found"))
	}

	groupExists, err := v.groupRepo.CheckGroupExists(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if !groupExists {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group not found"))
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint name is required"))
	}

	if len(name) > 255 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint name must be between 1 and 255 characters"))
	}

	if req.StartDate.IsZero() || req.EndDate.IsZero() {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("start date and end date are required"))
	}

	startDate := normalizeDateUTC(req.StartDate)
	endDate := normalizeDateUTC(req.EndDate)
	today := normalizeDateUTC(time.Now().UTC())

	if startDate.Before(today) {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("start date cannot be in the past"))
	}

	if !startDate.Before(endDate) {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("end date must be after start date"))
	}

	if endDate.Sub(startDate).Hours() > 24*30 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint duration must not exceed 30 days"))
	}

	var goal *string
	if req.Goal != nil {
		goalValue := strings.TrimSpace(*req.Goal)
		if goalValue == "" {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("goal must not be empty when provided"))
		}

		if len(goalValue) > 5000 {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("goal must not exceed 5000 characters"))
		}

		goal = &goalValue
	}

	isOverlap, err := v.sprintRepo.IsSprintOverlap(ctx, groupID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if isOverlap {
		return nil, errorbase.New(errdict.ErrConflict, errorbase.WithDetail("sprint date range overlaps an existing sprint"))
	}

	sprints, err := v.sprintRepo.GetSprintsByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	for _, sprint := range sprints {
		if strings.EqualFold(strings.TrimSpace(sprint.Name), name) {
			return nil, errorbase.New(errdict.ErrConflict, errorbase.WithDetail("sprint name already exists in this group"))
		}
	}

	newSprint, err := entity.NewSprint(
		uuid.NewString(),
		groupID,
		name,
		startDate,
		endDate,
		today,
	)
	if err != nil {
		return nil, err
	}

	newSprint.Goal = goal

	return newSprint, nil
}

func normalizeDateUTC(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
