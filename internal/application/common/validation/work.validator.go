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

type WorkValidator struct {
	workRepo   irepository.WorkRepository
	sprintRepo irepository.SprintRepository
	groupRepo  irepository.GroupRepository
	userRepo   irepository.UserRepository
}

type GetWorkPayload struct {
	GroupID string
	WorkID  string
}

type ListWorksPayload struct {
	GroupID  string
	SprintID *string
}

type UpdateWorkPayload struct {
	GroupID string
	WorkID  string
	Request *appdto.UpdateWorkRequest
}

type DeleteWorkPayload struct {
	GroupID string
	WorkID  string
}

type CreateChecklistItemPayload struct {
	GroupID string
	WorkID  string
	Item    *entity.ChecklistItem
}

type UpdateChecklistItemPayload struct {
	GroupID string
	WorkID  string
	ItemID  string
	Request *appdto.UpdateChecklistItemRequest
}

type DeleteChecklistItemPayload struct {
	GroupID string
	WorkID  string
	ItemID  string
}

type CreateCommentPayload struct {
	GroupID string
	WorkID  string
	Comment *entity.Comment
}

type UpdateCommentPayload struct {
	GroupID   string
	WorkID    string
	CommentID string
	Request   *appdto.UpdateCommentRequest
}

type DeleteCommentPayload struct {
	GroupID   string
	WorkID    string
	CommentID string
}

func NewWorkValidator(
	workRepo irepository.WorkRepository,
	sprintRepo irepository.SprintRepository,
	groupRepo irepository.GroupRepository,
	userRepo irepository.UserRepository,
) *WorkValidator {
	return &WorkValidator{
		workRepo:   workRepo,
		sprintRepo: sprintRepo,
		groupRepo:  groupRepo,
		userRepo:   userRepo,
	}
}

func (v *WorkValidator) ValidateCreateWork(
	ctx context.Context,
	req *appdto.CreateWorkRequest,
) (*entity.Work, errorbase.AppError) {
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

	creatorID := strings.TrimSpace(utils.GetUserIDFromOutgoingContext(ctx))
	if creatorID == "" {
		return nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("missing user context"))
	}

	if _, err := uuid.Parse(creatorID); err != nil {
		return nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("user id in context is invalid"))
	}

	groupExists, err := v.groupRepo.CheckGroupExists(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if !groupExists {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("group not found"))
	}

	user, err := v.userRepo.GetUserByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errorbase.New(errdict.ErrUnauthorized, errorbase.WithDetail("user not found"))
	}

	if user.Status != enum.UserStatusActive {
		return nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("inactive user is not allowed to create work"))
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("work name is required"))
	}

	if len(name) > 500 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("work name must be between 1 and 500 characters"))
	}

	var description *string
	if req.Description != nil {
		descriptionValue := strings.TrimSpace(*req.Description)
		if descriptionValue == "" {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("description must not be empty when provided"))
		}

		description = &descriptionValue
	}

	sprintID := ""
	if req.SprintID != nil {
		rawSprintID := strings.TrimSpace(*req.SprintID)
		if rawSprintID != "" {
			if _, err := uuid.Parse(rawSprintID); err != nil {
				return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id must be a valid UUID"))
			}

			sprint, err := v.sprintRepo.GetSprintByID(ctx, rawSprintID)
			if err != nil {
				return nil, err
			}

			if sprint == nil {
				return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("sprint not found"))
			}

			if sprint.GroupID != groupID {
				return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint does not belong to current group"))
			}

			if sprint.Status == enum.SprintStatusCompleted {
				return nil, errorbase.New(errdict.ErrUnprocessable, errorbase.WithDetail("cannot create work in completed sprint"))
			}

			if sprint.Status == enum.SprintStatusCancelled {
				return nil, errorbase.New(errdict.ErrUnprocessable, errorbase.WithDetail("cannot create work in cancelled sprint"))
			}

			works, err := v.workRepo.GetWorksBySprint(ctx, groupID, utils.Ptr(rawSprintID))
			if err != nil {
				return nil, err
			}

			if len(works) >= 250 {
				return nil, errorbase.New(errdict.ErrConflict, errorbase.WithDetail("sprint has reached maximum 250 works"))
			}

			sprintID = rawSprintID
		}
	}

	work, err := entity.NewWork(
		uuid.NewString(),
		groupID,
		sprintID,
		name,
		description,
		creatorID,
		"",
		nil,
		nil,
		nil,
		nil,
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}

	return work, nil
}

func (v *WorkValidator) ValidateGetWork(
	ctx context.Context,
	req *appdto.GetWorkRequest,
	actor *appdto.UserWithPermission,
) (*GetWorkPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	workID, err := validateRequiredUUID(req.WorkID, "work id")
	if err != nil {
		return nil, err
	}

	return &GetWorkPayload{
		GroupID: actor.GroupId,
		WorkID:  workID,
	}, nil
}

func (v *WorkValidator) ValidateListWorks(
	ctx context.Context,
	req *appdto.ListWorksRequest,
	actor *appdto.UserWithPermission,
) (*ListWorksPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	return &ListWorksPayload{
		GroupID:  actor.GroupId,
		SprintID: req.SprintID,
	}, nil
}

func (v *WorkValidator) ValidateUpdateWork(
	ctx context.Context,
	req *appdto.UpdateWorkRequest,
	actor *appdto.UserWithPermission,
) (*UpdateWorkPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	workID, err := validateRequiredUUID(req.WorkID, "work id")
	if err != nil {
		return nil, err
	}

	if req.Version < 0 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("version must be greater or equal than 0"))
	}

	providedFields := countProvidedWorkUpdateFields(req)
	if providedFields == 0 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("at least one field is required"))
	}

	if providedFields > 1 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("only one field can be updated per request"))
	}

	work, err := v.getWorkInGroup(ctx, workID, actor.GroupId)
	if err != nil {
		return nil, err
	}

	if err := v.validateMutableByCurrentSprint(ctx, work); err != nil {
		return nil, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("work name must not be empty when provided"))
		}

		if len(name) > 500 {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("work name must be between 1 and 500 characters"))
		}

		req.Name = &name
	}

	if req.Description != nil {
		description := strings.TrimSpace(*req.Description)
		if description == "" {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("description must not be empty when provided"))
		}

		if len(description) > 5000 {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("description must not exceed 5000 characters"))
		}

		req.Description = &description
	}

	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("invalid work status"))
		}
	}

	if req.SprintID != nil {
		sprintID := strings.TrimSpace(*req.SprintID)
		if sprintID == "" {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint id must not be empty when provided"))
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

		if sprint.GroupID != actor.GroupId {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("sprint does not belong to current group"))
		}

		if sprint.Status == enum.SprintStatusCompleted || sprint.Status == enum.SprintStatusCancelled {
			return nil, errorbase.New(errdict.ErrUnprocessable, errorbase.WithDetail("target sprint must be draft or active"))
		}

		currentSprintID := ""
		if work.SprintID != nil {
			currentSprintID = strings.TrimSpace(*work.SprintID)
		}

		if sprintID != currentSprintID {
			worksInSprint, err := v.workRepo.GetWorksBySprint(ctx, actor.GroupId, utils.Ptr(sprintID))
			if err != nil {
				return nil, err
			}

			if len(worksInSprint) >= 250 {
				return nil, errorbase.New(errdict.ErrConflict, errorbase.WithDetail("sprint has reached maximum 250 works"))
			}
		}

		req.SprintID = &sprintID
	}

	if req.AssigneeID != nil {
		assigneeID := strings.TrimSpace(*req.AssigneeID)
		if assigneeID == "" {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("assignee id must not be empty when provided"))
		}

		if _, err := uuid.Parse(assigneeID); err != nil {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("assignee id must be a valid UUID"))
		}

		assignee, err := v.userRepo.GetUserWithPermissionByID(ctx, actor.GroupId, assigneeID)
		if err != nil {
			return nil, err
		}

		if assignee == nil {
			return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("assignee not found in current group"))
		}

		if assignee.Status != enum.UserStatusActive {
			return nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("inactive user cannot be assigned"))
		}

		req.AssigneeID = &assigneeID
	}

	if req.StoryPoint != nil {
		if *req.StoryPoint <= 0 || *req.StoryPoint > 100 {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("story point must be between 1 and 100"))
		}

		if !isFibonacci(*req.StoryPoint) {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("story point must be a fibonacci number"))
		}
	}

	if req.Priority != nil {
		if !req.Priority.IsValid() {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("invalid work priority"))
		}
	}

	if req.DueDate != nil {
		if req.DueDate.IsZero() {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("due date is invalid"))
		}

		normalized := normalizeDateUTC(*req.DueDate)
		req.DueDate = &normalized
	}

	req.WorkID = workID

	return &UpdateWorkPayload{
		GroupID: actor.GroupId,
		WorkID:  workID,
		Request: req,
	}, nil
}

func (v *WorkValidator) ValidateDeleteWork(
	ctx context.Context,
	req *appdto.DeleteWorkRequest,
	actor *appdto.UserWithPermission,
) (*DeleteWorkPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	workID, err := validateRequiredUUID(req.WorkID, "work id")
	if err != nil {
		return nil, err
	}

	if _, err := v.getWorkInGroup(ctx, workID, actor.GroupId); err != nil {
		return nil, err
	}

	return &DeleteWorkPayload{
		GroupID: actor.GroupId,
		WorkID:  workID,
	}, nil
}

func (v *WorkValidator) ValidateCreateChecklistItem(
	ctx context.Context,
	req *appdto.CreateChecklistItemRequest,
	actor *appdto.UserWithPermission,
) (*CreateChecklistItemPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	workID, err := validateRequiredUUID(req.WorkID, "work id")
	if err != nil {
		return nil, err
	}

	work, err := v.getWorkInGroup(ctx, workID, actor.GroupId)
	if err != nil {
		return nil, err
	}

	if err := v.validateMutableByCurrentSprint(ctx, work); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("checklist name is required"))
	}

	if len(name) > 500 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("checklist name must be between 1 and 500 characters"))
	}

	item, err := entity.NewChecklistItem(
		uuid.NewString(),
		workID,
		name,
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}

	return &CreateChecklistItemPayload{
		GroupID: actor.GroupId,
		WorkID:  workID,
		Item:    item,
	}, nil
}

func (v *WorkValidator) ValidateUpdateChecklistItem(
	ctx context.Context,
	req *appdto.UpdateChecklistItemRequest,
	actor *appdto.UserWithPermission,
) (*UpdateChecklistItemPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	itemID, err := validateRequiredUUID(req.ItemID, "checklist item id")
	if err != nil {
		return nil, err
	}

	providedFields := 0
	if req.Name != nil {
		providedFields++
	}

	if req.IsCompleted != nil {
		providedFields++
	}

	if providedFields == 0 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("at least one field is required"))
	}

	itemMeta, err := v.workRepo.GetChecklistItemMeta(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if itemMeta == nil {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("checklist item not found"))
	}

	work, err := v.getWorkInGroup(ctx, itemMeta.WorkID, actor.GroupId)
	if err != nil {
		return nil, err
	}

	if err := v.validateMutableByCurrentSprint(ctx, work); err != nil {
		return nil, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("checklist name must not be empty when provided"))
		}

		if len(name) > 500 {
			return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("checklist name must be between 1 and 500 characters"))
		}

		req.Name = &name
	}

	req.ItemID = itemID

	return &UpdateChecklistItemPayload{
		GroupID: actor.GroupId,
		WorkID:  itemMeta.WorkID,
		ItemID:  itemID,
		Request: req,
	}, nil
}

func (v *WorkValidator) ValidateDeleteChecklistItem(
	ctx context.Context,
	req *appdto.DeleteChecklistItemRequest,
	actor *appdto.UserWithPermission,
) (*DeleteChecklistItemPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	itemID, err := validateRequiredUUID(req.ItemID, "checklist item id")
	if err != nil {
		return nil, err
	}

	itemMeta, err := v.workRepo.GetChecklistItemMeta(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if itemMeta == nil {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("checklist item not found"))
	}

	work, err := v.getWorkInGroup(ctx, itemMeta.WorkID, actor.GroupId)
	if err != nil {
		return nil, err
	}

	if err := v.validateMutableByCurrentSprint(ctx, work); err != nil {
		return nil, err
	}

	req.ItemID = itemID

	return &DeleteChecklistItemPayload{
		GroupID: actor.GroupId,
		WorkID:  itemMeta.WorkID,
		ItemID:  itemID,
	}, nil
}

func (v *WorkValidator) ValidateCreateComment(
	ctx context.Context,
	req *appdto.CreateCommentRequest,
	actor *appdto.UserWithPermission,
) (*CreateCommentPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	workID, err := validateRequiredUUID(req.WorkID, "work id")
	if err != nil {
		return nil, err
	}

	work, err := v.getWorkInGroup(ctx, workID, actor.GroupId)
	if err != nil {
		return nil, err
	}

	if err := v.validateMutableByCurrentSprint(ctx, work); err != nil {
		return nil, err
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("content is required"))
	}

	if len(content) > 5000 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("comment content must be between 1 and 5000 characters"))
	}

	comment, err := entity.NewComment(
		uuid.NewString(),
		workID,
		actor.ID,
		content,
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}

	return &CreateCommentPayload{
		GroupID: actor.GroupId,
		WorkID:  workID,
		Comment: comment,
	}, nil
}

func (v *WorkValidator) ValidateUpdateComment(
	ctx context.Context,
	req *appdto.UpdateCommentRequest,
	actor *appdto.UserWithPermission,
) (*UpdateCommentPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	commentID, err := validateRequiredUUID(req.CommentID, "comment id")
	if err != nil {
		return nil, err
	}

	commentMeta, err := v.workRepo.GetCommentMeta(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if commentMeta == nil {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("comment not found"))
	}

	work, err := v.getWorkInGroup(ctx, commentMeta.WorkID, actor.GroupId)
	if err != nil {
		return nil, err
	}

	if err := v.validateMutableByCurrentSprint(ctx, work); err != nil {
		return nil, err
	}

	if commentMeta.CreatorID != actor.ID {
		return nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("only comment creator can update comment"))
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("content is required"))
	}

	if len(content) > 5000 {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("comment content must be between 1 and 5000 characters"))
	}

	req.CommentID = commentID
	req.Content = content

	return &UpdateCommentPayload{
		GroupID:   actor.GroupId,
		WorkID:    commentMeta.WorkID,
		CommentID: commentID,
		Request:   req,
	}, nil
}

func (v *WorkValidator) ValidateDeleteComment(
	ctx context.Context,
	req *appdto.DeleteCommentRequest,
	actor *appdto.UserWithPermission,
) (*DeleteCommentPayload, errorbase.AppError) {
	if req == nil {
		return nil, errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail("request is required"))
	}

	commentID, err := validateRequiredUUID(req.CommentID, "comment id")
	if err != nil {
		return nil, err
	}

	commentMeta, err := v.workRepo.GetCommentMeta(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if commentMeta == nil {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("comment not found"))
	}

	work, err := v.getWorkInGroup(ctx, commentMeta.WorkID, actor.GroupId)
	if err != nil {
		return nil, err
	}

	if err := v.validateMutableByCurrentSprint(ctx, work); err != nil {
		return nil, err
	}

	if !actor.Role.HasPermission(enum.GroupRoleManager) && commentMeta.CreatorID != actor.ID {
		return nil, errorbase.New(errdict.ErrForbidden, errorbase.WithDetail("only creator, manager, or owner can delete comment"))
	}

	req.CommentID = commentID

	return &DeleteCommentPayload{
		GroupID:   actor.GroupId,
		WorkID:    commentMeta.WorkID,
		CommentID: commentID,
	}, nil
}

func (v *WorkValidator) getWorkInGroup(
	ctx context.Context,
	workID string,
	groupID string,
) (*appdto.WorkResponse, errorbase.AppError) {
	work, err := v.workRepo.GetWorkAggregation(ctx, workID)
	if err != nil {
		return nil, err
	}

	if work == nil || work.GroupID != groupID {
		return nil, errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("work not found"))
	}

	return work, nil
}

func (v *WorkValidator) validateMutableByCurrentSprint(
	ctx context.Context,
	work *appdto.WorkResponse,
) errorbase.AppError {
	if work == nil || work.SprintID == nil {
		return nil
	}

	currentSprintID := strings.TrimSpace(*work.SprintID)
	if currentSprintID == "" {
		return nil
	}

	sprint, err := v.sprintRepo.GetSprintByID(ctx, currentSprintID)
	if err != nil {
		return err
	}

	if sprint == nil {
		return errorbase.New(errdict.ErrNotFound, errorbase.WithDetail("sprint not found"))
	}

	if sprint.Status == enum.SprintStatusCompleted || sprint.Status == enum.SprintStatusCancelled {
		return errorbase.New(errdict.ErrUnprocessable, errorbase.WithDetail("cannot update work in completed or cancelled sprint"))
	}

	return nil
}

func validateRequiredUUID(raw string, field string) (string, errorbase.AppError) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail(field+" is required"))
	}

	if _, err := uuid.Parse(value); err != nil {
		return "", errorbase.New(errdict.ErrBadRequest, errorbase.WithDetail(field+" must be a valid UUID"))
	}

	return value, nil
}

func countProvidedWorkUpdateFields(req *appdto.UpdateWorkRequest) int {
	count := 0

	if req.Name != nil {
		count++
	}

	if req.Description != nil {
		count++
	}

	if req.Status != nil {
		count++
	}

	if req.SprintID != nil {
		count++
	}

	if req.AssigneeID != nil {
		count++
	}

	if req.StoryPoint != nil {
		count++
	}

	if req.Priority != nil {
		count++
	}

	if req.DueDate != nil {
		count++
	}

	return count
}

func isFibonacci(n int32) bool {
	if n <= 0 {
		return false
	}

	if n == 1 {
		return true
	}

	a, b := int32(1), int32(1)
	for b < n {
		a, b = b, a+b
	}

	return b == n
}
