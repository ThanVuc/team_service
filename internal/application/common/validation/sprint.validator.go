package appvalidation

import (
	"context"
	"strings"
	appdto "team_service/internal/application/common/dto"
	irepository "team_service/internal/application/common/interface/repository"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"team_service/internal/infrastructure/share/utils"
	"time"

	"github.com/google/uuid"
)

type SprintValidator struct {
	sprintRepo irepository.SprintRepository
	userRepo   irepository.UserRepository
	groupRepo  irepository.GroupRepository
}

type UpdateSprintPayload struct {
	SprintID  string
	Name      *string
	Goal      *string
	StartDate *time.Time
	EndDate   *time.Time
}

type UpdateSprintStatusPayload struct {
	SprintID string
	Status   enum.SprintStatus
}

type DeleteSprintPayload struct {
	SprintID string
}

type ExportSprintPayload struct {
	SprintID string
	GroupID  string
}

type GenerateSprintPayload struct {
	GroupID           string
	Name              string
	Goal              string
	StartDate         string
	EndDate           string
	AdditionalContext *string
	Files             []appdto.AISprintGenerationFile
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
	if req.Goal != nil && strings.TrimSpace(*req.Goal) != "" {
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

func (v *SprintValidator) ValidateGetSprint(
	ctx context.Context,
	req *appdto.GetSprintRequest,
) (string, errorbase.AppError) {
	if req == nil {
		return "", errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	sprintID := strings.TrimSpace(req.SprintID)
	if sprintID == "" {
		return "", errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id is required"))
	}

	if _, err := uuid.Parse(sprintID); err != nil {
		return "", errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id must be a valid UUID"))
	}

	return sprintID, nil
}

func (v *SprintValidator) ValidateListSprints(
	ctx context.Context,
	req *appdto.ListSprintsRequest,
) (string, errorbase.AppError) {
	if req == nil {
		return "", errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	groupID := strings.TrimSpace(req.GroupID)
	if groupID == "" {
		return "", errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required"))
	}

	if _, err := uuid.Parse(groupID); err != nil {
		return "", errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id must be a valid UUID"))
	}

	userID := strings.TrimSpace(utils.GetUserIDFromOutgoingContext(ctx))
	if userID == "" {
		return "", errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("missing user context"))
	}

	if _, err := uuid.Parse(userID); err != nil {
		return "", errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("user id in context is invalid"))
	}

	groupExists, err := v.groupRepo.CheckGroupExists(ctx, groupID)
	if err != nil {
		return "", err
	}

	if !groupExists {
		return "", errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group not found"))
	}

	return groupID, nil
}

func (v *SprintValidator) ValidateUpdateSprint(
	ctx context.Context,
	req *appdto.UpdateSprintRequest,
) (*UpdateSprintPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	sprintID := strings.TrimSpace(req.SprintID)
	if sprintID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id is required"))
	}

	if _, err := uuid.Parse(sprintID); err != nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id must be a valid UUID"))
	}

	if req.Name == nil && req.Goal == nil && req.StartDate == nil && req.EndDate == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("at least one field is required"))
	}

	sprint, err := v.sprintRepo.GetSprintByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	if sprint == nil {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("sprint not found"))
	}

	if sprint.Status != enum.SprintStatusDraft {
		return nil, errorbase.New(errdict.ErrUnprocessable, errorbase.WithDetail("only draft sprint can be updated"))
	}

	effectiveStart := normalizeDateUTC(sprint.StartDate)
	effectiveEnd := normalizeDateUTC(sprint.EndDate)
	today := normalizeDateUTC(time.Now().UTC())

	var name *string
	if req.Name != nil {
		nameValue := strings.TrimSpace(*req.Name)
		if nameValue == "" {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint name must not be empty when provided"))
		}

		if len(nameValue) > 255 {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint name must be between 1 and 255 characters"))
		}

		name = &nameValue
	}

	var goal *string
	if req.Goal != nil && strings.TrimSpace(*req.Goal) != "" {
		goalValue := strings.TrimSpace(*req.Goal)
		if goalValue == "" {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("goal must not be empty when provided"))
		}

		if len(goalValue) > 5000 {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("goal must not exceed 5000 characters"))
		}

		goal = &goalValue
	}

	var startDate *time.Time
	if req.StartDate != nil {
		startValue := normalizeDateUTC(*req.StartDate)
		if startValue.Before(today) {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("start date cannot be in the past"))
		}

		effectiveStart = startValue
		startDate = &startValue
	}

	var endDate *time.Time
	if req.EndDate != nil {
		endValue := normalizeDateUTC(*req.EndDate)
		effectiveEnd = endValue
		endDate = &endValue
	}

	if !effectiveStart.Before(effectiveEnd) {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("end date must be after start date"))
	}

	if effectiveEnd.Sub(effectiveStart).Hours() > 24*30 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint duration must not exceed 30 days"))
	}

	if name != nil || startDate != nil || endDate != nil {
		sprints, err := v.sprintRepo.GetSprintsByGroupID(ctx, sprint.GroupID)
		if err != nil {
			return nil, err
		}

		for _, existing := range sprints {
			if existing == nil || existing.ID == sprintID {
				continue
			}

			if name != nil && strings.EqualFold(strings.TrimSpace(existing.Name), *name) {
				return nil, errorbase.New(errdict.ErrConflict, errorbase.WithDetail("sprint name already exists in this group"))
			}

			if (startDate != nil || endDate != nil) && existing.Status != enum.SprintStatusCancelled {
				existingStart := normalizeDateUTC(existing.StartDate)
				existingEnd := normalizeDateUTC(existing.EndDate)
				if isDateRangeOverlapInclusive(effectiveStart, effectiveEnd, existingStart, existingEnd) {
					return nil, errorbase.New(errdict.ErrConflict, errorbase.WithDetail("sprint date range overlaps an existing sprint"))
				}
			}
		}
	}

	return &UpdateSprintPayload{
		SprintID:  sprintID,
		Name:      name,
		Goal:      goal,
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

func (v *SprintValidator) ValidateUpdateSprintStatus(
	ctx context.Context,
	req *appdto.UpdateSprintStatusRequest,
) (*UpdateSprintStatusPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	sprintID := strings.TrimSpace(req.SprintID)
	if sprintID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id is required"))
	}

	if _, err := uuid.Parse(sprintID); err != nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id must be a valid UUID"))
	}

	if !req.Status.IsValid() {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("invalid sprint status"))
	}

	sprint, err := v.sprintRepo.GetSprintByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	if sprint == nil {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("sprint not found"))
	}

	sprintToValidate := *sprint
	if err := sprintToValidate.ChangeStatus(req.Status); err != nil {
		return nil, err
	}

	if req.Status == enum.SprintStatusActive {
		sprints, err := v.sprintRepo.GetSprintsByGroupID(ctx, sprint.GroupID)
		if err != nil {
			return nil, err
		}

		for _, existing := range sprints {
			if existing == nil || existing.ID == sprint.ID {
				continue
			}

			if existing.Status == enum.SprintStatusActive {
				return nil, errorbase.New(errdict.ErrConflict, errorbase.WithDetail("another active sprint already exists in this group"))
			}
		}
	}

	return &UpdateSprintStatusPayload{
		SprintID: sprintID,
		Status:   req.Status,
	}, nil
}

func (v *SprintValidator) ValidateDeleteSprint(
	ctx context.Context,
	req *appdto.DeleteSprintRequest,
) (*DeleteSprintPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	sprintID := strings.TrimSpace(req.SprintID)
	if sprintID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id is required"))
	}

	if _, err := uuid.Parse(sprintID); err != nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id must be a valid UUID"))
	}

	sprint, err := v.sprintRepo.GetSprintByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	if sprint == nil {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("sprint not found"))
	}

	if sprint.Status != enum.SprintStatusDraft {
		return nil, errorbase.New(errdict.ErrUnprocessable, errorbase.WithDetail("only draft sprint can be deleted"))
	}

	return &DeleteSprintPayload{SprintID: sprintID}, nil
}

func (v *SprintValidator) ValidateExportSprint(
	ctx context.Context,
	req *appdto.ExportSprintRequest,
) (*ExportSprintPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	sprintID := strings.TrimSpace(req.SprintID)
	if sprintID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id is required"))
	}

	if _, err := uuid.Parse(sprintID); err != nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id must be a valid UUID"))
	}

	groupID := strings.TrimSpace(utils.GetGroupIDFromContext(ctx))
	if groupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required in context"))
	}

	if _, err := uuid.Parse(groupID); err != nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id in context must be a valid UUID"))
	}

	groupExists, err := v.groupRepo.CheckGroupExists(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if !groupExists {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group not found"))
	}

	sprint, err := v.sprintRepo.GetSprintByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	if sprint == nil {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("sprint not found"))
	}

	if sprint.GroupID != groupID {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("sprint not found in group"))
	}

	return &ExportSprintPayload{
		SprintID: sprintID,
		GroupID:  groupID,
	}, nil
}

func (v *SprintValidator) ValidateGenerateSprint(
	ctx context.Context,
	req *appdto.GenerateSprintRequest,
) (*GenerateSprintPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	groupID := strings.TrimSpace(utils.GetGroupIDFromContext(ctx))
	if groupID == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id is required in context"))
	}

	if _, err := uuid.Parse(groupID); err != nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("group id in context must be a valid UUID"))
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint name is required"))
	}

	if len(name) > 255 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint name must not exceed 255 characters"))
	}

	goal := strings.TrimSpace(req.Goal)
	if goal == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint goal is required"))
	}

	if len(goal) > 1000 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint goal must not exceed 1000 characters"))
	}

	startDate := strings.TrimSpace(req.StartDate)
	if startDate == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("start date is required"))
	}

	endDate := strings.TrimSpace(req.EndDate)
	if endDate == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("end date is required"))
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("start date must follow YYYY-MM-DD format"))
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("end date must follow YYYY-MM-DD format"))
	}

	if !start.Before(end) {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("end date must be after start date"))
	}

	isOverlap, err := v.sprintRepo.IsSprintOverlap(ctx, groupID, normalizeDateUTC(start), normalizeDateUTC(end))
	if err != nil {
		return nil, errorbase.New(errdict.ErrInternal, errorbase.WithDetail("failed to validate sprint date range overlap"))
	}

	if isOverlap {
		return nil, errorbase.New(errdict.ErrConflict, errorbase.WithDetail("sprint date range overlaps an existing sprint"))
	}

	var additionalContext *string
	if req.AdditionalContext != nil {
		value := strings.TrimSpace(*req.AdditionalContext)
		if len(value) > 2000 {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("additional context must not exceed 2000 characters"))
		}

		if value != "" {
			additionalContext = &value
		}
	}

	files := make([]appdto.AISprintGenerationFile, 0, len(req.Files))
	for _, file := range req.Files {
		objectKey := strings.TrimSpace(file.ObjectKey)
		if objectKey == "" {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("file object key is required"))
		}

		if file.Size <= 0 {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("file size must be greater than 0"))
		}

		files = append(files, appdto.AISprintGenerationFile{
			ObjectKey: objectKey,
			Size:      file.Size,
		})
	}

	return &GenerateSprintPayload{
		GroupID:           groupID,
		Name:              name,
		Goal:              goal,
		StartDate:         startDate,
		EndDate:           endDate,
		AdditionalContext: additionalContext,
		Files:             files,
	}, nil
}

func isDateRangeOverlapInclusive(startA, endA, startB, endB time.Time) bool {
	return !startA.After(endB) && !startB.After(endA)
}

func normalizeDateUTC(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
