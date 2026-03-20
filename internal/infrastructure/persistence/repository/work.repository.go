package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	appdto "team_service/internal/application/common/dto"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/internal/infrastructure/share/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type WorkRepository struct {
	q *database.Queries
}

func NewWorkRepository(
	q *database.Queries,
) *WorkRepository {
	return &WorkRepository{
		q: q,
	}
}

func (r *WorkRepository) CreateWork(
	ctx context.Context,
	work *entity.Work,
) (*entity.Work, errorbase.AppError) {
	id, appErr := parseUUID(work.ID, "work id")
	if appErr != nil {
		return nil, appErr
	}

	groupID, appErr := parseUUID(work.GroupID, "group id")
	if appErr != nil {
		return nil, appErr
	}

	creatorID, appErr := parseUUID(work.CreatorID, "creator id")
	if appErr != nil {
		return nil, appErr
	}

	sprintID, appErr := parseOptionalUUIDString(work.SprintID, "sprint id")
	if appErr != nil {
		return nil, appErr
	}

	created, err := r.q.CreateWork(ctx, database.CreateWorkParams{
		ID:          id,
		GroupID:     groupID,
		SprintID:    sprintID,
		Name:        work.Name,
		Description: toNullableText(work.Description),
		CreatorID:   creatorID,
	})
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to create work id=%s", work.ID)),
		)
	}

	return mapDBWorkToEntity(created), nil
}

func (r *WorkRepository) UpdateWork(
	ctx context.Context,
	req *appdto.UpdateWorkRequest,
) (*appdto.WorkResponse, errorbase.AppError) {
	id, appErr := parseUUID(req.WorkID, "work id")
	if appErr != nil {
		return nil, appErr
	}

	sprintID, appErr := parseOptionalUUIDPtr(req.SprintID, "sprint id")
	if appErr != nil {
		return nil, appErr
	}

	assigneeID, appErr := parseOptionalUUIDPtr(req.AssigneeID, "assignee id")
	if appErr != nil {
		return nil, appErr
	}

	updated, err := r.q.UpdateWork(ctx, database.UpdateWorkParams{
		Name:        toNullableText(req.Name),
		Description: toNullableText(req.Description),
		SprintID:    sprintID,
		AssigneeID:  assigneeID,
		Status:      toNullableWorkStatus(req.Status),
		StoryPoint:  toNullableInt4(req.StoryPoint),
		DueDate:     toNullableDate(req.DueDate),
		Priority:    toNullableWorkPriority(req.Priority),
		ID:          id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorbase.New(
				errdict.ErrNotFound,
				errorbase.WithDetail(fmt.Sprintf("work not found id=%s", req.WorkID)),
			)
		}

		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to update work id=%s", req.WorkID)),
		)
	}

	resp := &appdto.WorkResponse{
		ID:          updated.ID.String(),
		Name:        utils.SafeString(fromAnyText(updated.Name)),
		Description: fromAnyText(updated.Description),
		SprintID:    fromAnyUUID(updated.SprintID),
		AssigneeID:  fromAnyUUID(updated.AssigneeID),
		Status:      enum.WorkStatus(utils.SafeString(fromAnyText(updated.Status))),
		StoryPoint:  fromAnyInt32(updated.StoryPoint),
		DueDate:     fromAnyDate(updated.DueDate),
		Priority:    enum.WorkPriority(utils.SafeString(fromAnyText(updated.Priority))),
		UpdatedAt:   updated.UpdatedAt.Time,
		Version:     req.Version,
	}

	return resp, nil
}

func (r *WorkRepository) DeleteWork(
	ctx context.Context,
	workID string,
) (*appdto.DeleteWorkResponse, errorbase.AppError) {
	id, appErr := parseUUID(workID, "work id")
	if appErr != nil {
		return nil, appErr
	}

	deleted, err := r.q.DeleteWork(ctx, id)
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to delete work id=%s", workID)),
		)
	}

	if !deleted.Success {
		return nil, errorbase.New(
			errdict.ErrNotFound,
			errorbase.WithDetail(fmt.Sprintf("work not found id=%s", workID)),
		)
	}

	return &appdto.DeleteWorkResponse{Success: true}, nil
}

func (r *WorkRepository) GetWorksBySprint(
	ctx context.Context,
	groupID string,
	sprintID *string,
) ([]appdto.WorkResponse, errorbase.AppError) {
	groupUUID, appErr := parseUUID(groupID, "group id")
	if appErr != nil {
		return nil, appErr
	}

	pgSprintID, err := utils.StringPtrToPgUUID(sprintID)
	if err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to parse sprint id=%s", utils.SafeString(sprintID))),
		)
	}

	rows, err := r.q.GetWorksBySprint(ctx, database.GetWorksBySprintParams{
		GroupID:  groupUUID,
		SprintID: pgSprintID,
	})
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to get works by sprint id=%s group id=%s", utils.SafeString(sprintID), groupID)),
		)
	}

	works := make([]appdto.WorkResponse, 0, len(rows))
	for _, row := range rows {
		works = append(works, mapGetWorksBySprintRowToDTO(row))
	}

	return works, nil
}

func (r *WorkRepository) GetWorkAggregation(
	ctx context.Context,
	workID string,
) (*appdto.WorkResponse, errorbase.AppError) {
	id, appErr := parseUUID(workID, "work id")
	if appErr != nil {
		return nil, appErr
	}

	row, err := r.q.GetWork(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorbase.New(
				errdict.ErrNotFound,
				errorbase.WithDetail(fmt.Sprintf("work not found id=%s", workID)),
			)
		}

		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to get work id=%s", workID)),
		)
	}

	resp := mapGetWorkRowToDTO(row)

	checklistRows, err := r.q.GetCheckListByWorkId(ctx, id)
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to get checklist by work id=%s", workID)),
		)
	}

	checklistItems := make([]appdto.ChecklistItemResponse, 0, len(checklistRows))
	var completedChecklist int32
	for _, item := range checklistRows {
		checklistItems = append(checklistItems, appdto.ChecklistItemResponse{
			ID:          item.ID.String(),
			WorkID:      item.WorkID.String(),
			Name:        item.Name,
			IsCompleted: item.IsCompleted,
			CreatedAt:   item.CreatedAt.Time,
			UpdatedAt:   item.UpdatedAt.Time,
		})

		if item.IsCompleted {
			completedChecklist++
		}
	}

	resp.CheckList = &appdto.ChecklistSummaryResponse{
		Total:     int32(len(checklistItems)),
		Completed: completedChecklist,
		Items:     checklistItems,
	}

	commentRows, err := r.q.GetCommentsByWorkId(ctx, id)
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to get comments by work id=%s", workID)),
		)
	}

	commentItems := make([]appdto.CommentResponse, 0, len(commentRows))
	for _, comment := range commentRows {
		commentItems = append(commentItems, appdto.CommentResponse{
			ID:      comment.ID.String(),
			Content: comment.Content,
			Creator: appdto.UserSummaryDTO{
				ID:     comment.CreatorID.String(),
				Name:   comment.CreatorEmail,
				Email:  comment.CreatorEmail,
				Avatar: nullableTextToPtr(comment.CreatorAvatarUrl),
			},
			CreatedAt: comment.CreatedAt.Time,
			UpdatedAt: comment.UpdatedAt.Time,
		})
	}

	resp.Comments = &appdto.CommentListResponse{
		Total:    int32(len(commentItems)),
		Comments: commentItems,
	}

	return &resp, nil
}

func (r *WorkRepository) GetChecklistItemMeta(
	ctx context.Context,
	itemID string,
) (*appdto.ChecklistItemMeta, errorbase.AppError) {
	id, appErr := parseUUID(itemID, "checklist item id")
	if appErr != nil {
		return nil, appErr
	}

	row, err := r.q.GetChecklistItemMeta(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorbase.New(
				errdict.ErrNotFound,
				errorbase.WithDetail(fmt.Sprintf("checklist item not found id=%s", itemID)),
			)
		}

		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to get checklist item id=%s", itemID)),
		)
	}

	return &appdto.ChecklistItemMeta{
		ID:     row.ID.String(),
		WorkID: row.WorkID.String(),
	}, nil
}

func (r *WorkRepository) GetCommentMeta(
	ctx context.Context,
	commentID string,
) (*appdto.CommentMeta, errorbase.AppError) {
	id, appErr := parseUUID(commentID, "comment id")
	if appErr != nil {
		return nil, appErr
	}

	row, err := r.q.GetCommentMeta(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorbase.New(
				errdict.ErrNotFound,
				errorbase.WithDetail(fmt.Sprintf("comment not found id=%s", commentID)),
			)
		}

		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to get comment id=%s", commentID)),
		)
	}

	return &appdto.CommentMeta{
		ID:        row.ID.String(),
		WorkID:    row.WorkID.String(),
		CreatorID: row.CreatorID.String(),
	}, nil
}

func (r *WorkRepository) CreateChecklistItem(
	ctx context.Context,
	item *entity.ChecklistItem,
) (*appdto.ChecklistItemResponse, errorbase.AppError) {
	id, appErr := parseUUID(item.ID, "checklist item id")
	if appErr != nil {
		return nil, appErr
	}

	workID, appErr := parseUUID(item.WorkID, "work id")
	if appErr != nil {
		return nil, appErr
	}

	created, err := r.q.CreateChecklistItem(ctx, database.CreateChecklistItemParams{
		ID:     id,
		WorkID: workID,
		Name:   item.Name,
	})
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to create checklist item id=%s", item.ID)),
		)
	}

	return &appdto.ChecklistItemResponse{
		ID:          created.ID.String(),
		WorkID:      created.WorkID.String(),
		Name:        created.Name,
		IsCompleted: created.IsCompleted,
		CreatedAt:   created.CreatedAt.Time,
		UpdatedAt:   created.UpdatedAt.Time,
	}, nil
}

func (r *WorkRepository) UpdateChecklistItem(
	ctx context.Context,
	req *appdto.UpdateChecklistItemRequest,
) (*appdto.ChecklistItemResponse, errorbase.AppError) {
	id, appErr := parseUUID(req.ItemID, "checklist item id")
	if appErr != nil {
		return nil, appErr
	}

	updated, err := r.q.UpdateChecklistItem(ctx, database.UpdateChecklistItemParams{
		Name:        toNullableText(req.Name),
		IsCompleted: toNullableBool(req.IsCompleted),
		ID:          id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorbase.New(
				errdict.ErrNotFound,
				errorbase.WithDetail(fmt.Sprintf("checklist item not found id=%s", req.ItemID)),
			)
		}

		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to update checklist item id=%s", req.ItemID)),
		)
	}

	name := fromAnyText(updated.Name)
	isCompleted := fromAnyBool(updated.IsCompleted)

	return &appdto.ChecklistItemResponse{
		ID:          updated.ID.String(),
		Name:        utils.SafeString(name),
		IsCompleted: boolValue(isCompleted),
		UpdatedAt:   updated.UpdatedAt.Time,
	}, nil
}

func (r *WorkRepository) DeleteChecklistItem(
	ctx context.Context,
	itemID string,
) (*appdto.ChecklistItemResponse, errorbase.AppError) {
	id, appErr := parseUUID(itemID, "checklist item id")
	if appErr != nil {
		return nil, appErr
	}

	deleted, err := r.q.DeleteChecklistItem(ctx, id)
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to delete checklist item id=%s", itemID)),
		)
	}

	if !deleted.Success {
		return nil, errorbase.New(
			errdict.ErrNotFound,
			errorbase.WithDetail(fmt.Sprintf("checklist item not found id=%s", itemID)),
		)
	}

	return &appdto.ChecklistItemResponse{ID: deleted.ID.String()}, nil
}

func (r *WorkRepository) CreateComment(
	ctx context.Context,
	comment *entity.Comment,
) (*appdto.CommentListResponse, errorbase.AppError) {
	id, appErr := parseUUID(comment.ID, "comment id")
	if appErr != nil {
		return nil, appErr
	}

	workID, appErr := parseUUID(comment.WorkID, "work id")
	if appErr != nil {
		return nil, appErr
	}

	creatorID, appErr := parseUUID(comment.CreatorID, "creator id")
	if appErr != nil {
		return nil, appErr
	}

	created, err := r.q.CreateComment(ctx, database.CreateCommentParams{
		ID:        id,
		WorkID:    workID,
		CreatorID: creatorID,
		Content:   comment.Content,
	})
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to create comment id=%s", comment.ID)),
		)
	}

	commentResp := appdto.CommentResponse{
		ID:      created.ID.String(),
		Content: created.Content,
		Creator: appdto.UserSummaryDTO{
			ID:     created.CreatorID.String(),
			Name:   created.CreatorEmail,
			Email:  created.CreatorEmail,
			Avatar: nullableTextToPtr(created.CreatorAvatarUrl),
		},
		CreatedAt: created.CreatedAt.Time,
		UpdatedAt: created.UpdatedAt.Time,
	}

	return &appdto.CommentListResponse{
		Total:    1,
		Comments: []appdto.CommentResponse{commentResp},
	}, nil
}

func (r *WorkRepository) UpdateComment(
	ctx context.Context,
	req *appdto.UpdateCommentRequest,
) (*appdto.CommentListResponse, errorbase.AppError) {
	id, appErr := parseUUID(req.CommentID, "comment id")
	if appErr != nil {
		return nil, appErr
	}

	updated, err := r.q.UpdateComment(ctx, database.UpdateCommentParams{
		Content: req.Content,
		ID:      id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorbase.New(
				errdict.ErrNotFound,
				errorbase.WithDetail(fmt.Sprintf("comment not found id=%s", req.CommentID)),
			)
		}

		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to update comment id=%s", req.CommentID)),
		)
	}

	creator, err := r.q.GetUserByID(ctx, updated.CreatorID)
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to get comment creator id=%s", updated.CreatorID.String())),
		)
	}

	commentResp := appdto.CommentResponse{
		ID:      updated.ID.String(),
		Content: updated.Content,
		Creator: appdto.UserSummaryDTO{
			ID:     creator.ID.String(),
			Name:   creator.Email,
			Email:  creator.Email,
			Avatar: nullableTextToPtr(creator.AvatarUrl),
		},
		CreatedAt: updated.CreatedAt.Time,
		UpdatedAt: updated.UpdatedAt.Time,
	}

	return &appdto.CommentListResponse{
		Total:    1,
		Comments: []appdto.CommentResponse{commentResp},
	}, nil
}

func (r *WorkRepository) DeleteComment(
	ctx context.Context,
	commentID string,
) (*appdto.CommentListResponse, errorbase.AppError) {
	id, appErr := parseUUID(commentID, "comment id")
	if appErr != nil {
		return nil, appErr
	}

	deleted, err := r.q.DeleteComment(ctx, id)
	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to delete comment id=%s", commentID)),
		)
	}

	if !deleted.Success {
		return nil, errorbase.New(
			errdict.ErrNotFound,
			errorbase.WithDetail(fmt.Sprintf("comment not found id=%s", commentID)),
		)
	}

	return &appdto.CommentListResponse{
		Total:    0,
		Comments: make([]appdto.CommentResponse, 0),
	}, nil
}

func parseUUID(raw string, field string) (pgtype.UUID, errorbase.AppError) {
	u, err := utils.ToUUID(raw)
	if err != nil || !u.Valid {
		return pgtype.UUID{}, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to parse %s", field)),
		)
	}

	return u, nil
}

func parseOptionalUUIDString(raw string, field string) (pgtype.UUID, errorbase.AppError) {
	if raw == "" {
		return pgtype.UUID{}, nil
	}

	return parseUUID(raw, field)
}

func parseOptionalUUIDPtr(raw *string, field string) (pgtype.UUID, errorbase.AppError) {
	if raw == nil {
		return pgtype.UUID{}, nil
	}

	if *raw == "" {
		return pgtype.UUID{}, nil
	}

	return parseUUID(*raw, field)
}

func toNullableText(v *string) pgtype.Text {
	if v == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{String: *v, Valid: true}
}

func toNullableBool(v *bool) pgtype.Bool {
	if v == nil {
		return pgtype.Bool{}
	}

	return pgtype.Bool{Bool: *v, Valid: true}
}

func toNullableInt4(v *int32) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{}
	}

	return pgtype.Int4{Int32: *v, Valid: true}
}

func toNullableDate(v *time.Time) pgtype.Date {
	if v == nil {
		return pgtype.Date{}
	}

	return pgtype.Date{Time: *v, Valid: true}
}

func toNullableWorkStatus(v *enum.WorkStatus) pgtype.Text {
	if v == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{String: string(*v), Valid: true}
}

func toNullableWorkPriority(v *enum.WorkPriority) pgtype.Text {
	if v == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{String: string(*v), Valid: true}
}

func nullableTextToPtr(v pgtype.Text) *string {
	if !v.Valid {
		return nil
	}

	return utils.Ptr(v.String)
}

func nullableUUIDToPtr(v pgtype.UUID) *string {
	if !v.Valid {
		return nil
	}

	return utils.Ptr(v.String())
}

func nullableDateToPtr(v pgtype.Date) *time.Time {
	if !v.Valid {
		return nil
	}

	return utils.Ptr(v.Time)
}

func nullableInt4ToPtr(v pgtype.Int4) *int32 {
	if !v.Valid {
		return nil
	}

	return utils.Ptr(v.Int32)
}

func nullableFloat8ToPtr(v pgtype.Float8) *float64 {
	if !v.Valid {
		return nil
	}

	return utils.Ptr(v.Float64)
}

func mapDBWorkToEntity(w database.Work) *entity.Work {
	var priority *enum.WorkPriority
	if w.Priority.Valid {
		p := enum.WorkPriority(w.Priority.String)
		priority = &p
	}

	return &entity.Work{
		ID:            w.ID.String(),
		GroupID:       w.GroupID.String(),
		SprintID:      utils.SafeString(nullableUUIDToPtr(w.SprintID)),
		Name:          w.Name,
		Description:   nullableTextToPtr(w.Description),
		Status:        enum.WorkStatus(w.Status),
		AssigneeID:    utils.SafeString(nullableUUIDToPtr(w.AssigneeID)),
		CreatorID:     w.CreatorID.String(),
		EstimateHours: nullableFloat8ToPtr(w.EstimateHours),
		StoryPoint:    nullableInt4ToPtr(w.StoryPoint),
		Priority:      priority,
		DueDate:       nullableDateToPtr(w.DueDate),
		CreatedAt:     w.CreatedAt.Time,
		UpdatedAt:     w.UpdatedAt.Time,
	}
}

func mapGetWorkRowToDTO(row database.GetWorkRow) appdto.WorkResponse {
	resp := appdto.WorkResponse{
		ID:            row.ID.String(),
		GroupID:       row.GroupID.String(),
		SprintID:      nullableUUIDToPtr(row.SprintID),
		Name:          row.Name,
		Description:   nullableTextToPtr(row.Description),
		Status:        enum.WorkStatus(row.Status),
		Priority:      enum.WorkPriority(utils.SafeString(nullableTextToPtr(row.Priority))),
		AssigneeID:    nullableUUIDToPtr(row.AssigneeID),
		CreatorID:     row.CreatorID.String(),
		EstimateHours: nullableFloat8ToPtr(row.EstimateHours),
		StoryPoint:    nullableInt4ToPtr(row.StoryPoint),
		DueDate:       nullableDateToPtr(row.DueDate),
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}

	if row.SprintID.Valid {
		resp.Sprint = &appdto.SimpleSprintDTO{
			ID:   row.SprintID.String(),
			Name: utils.SafeString(nullableTextToPtr(row.SprintName)),
		}
	}

	if row.AssigneeID.Valid {
		resp.Assignee = &appdto.SimpleUserDTO{
			ID:     row.AssigneeID.String(),
			Email:  utils.SafeString(nullableTextToPtr(row.AssigneeEmail)),
			Avatar: nullableTextToPtr(row.AssigneeAvatarUrl),
		}
	}

	return resp
}

func mapGetWorksBySprintRowToDTO(row database.GetWorksBySprintRow) appdto.WorkResponse {
	resp := appdto.WorkResponse{
		ID:            row.ID.String(),
		GroupID:       row.GroupID.String(),
		SprintID:      nullableUUIDToPtr(row.SprintID),
		Name:          row.Name,
		Description:   nullableTextToPtr(row.Description),
		Status:        enum.WorkStatus(row.Status),
		Priority:      enum.WorkPriority(utils.SafeString(nullableTextToPtr(row.Priority))),
		AssigneeID:    nullableUUIDToPtr(row.AssigneeID),
		CreatorID:     row.CreatorID.String(),
		EstimateHours: nullableFloat8ToPtr(row.EstimateHours),
		StoryPoint:    nullableInt4ToPtr(row.StoryPoint),
		DueDate:       nullableDateToPtr(row.DueDate),
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}

	if row.AssigneeID.Valid {
		resp.Assignee = &appdto.SimpleUserDTO{
			ID:     row.AssigneeID.String(),
			Email:  utils.SafeString(nullableTextToPtr(row.AssigneeEmail)),
			Avatar: nullableTextToPtr(row.AssigneeAvatarUrl),
		}
	}

	return resp
}

func fromAnyText(v interface{}) *string {
	if v == nil {
		return nil
	}

	var text pgtype.Text
	if err := text.Scan(v); err != nil || !text.Valid {
		return nil
	}

	return utils.Ptr(text.String)
}

func fromAnyUUID(v interface{}) *string {
	if v == nil {
		return nil
	}

	var id pgtype.UUID
	if err := id.Scan(v); err != nil || !id.Valid {
		return nil
	}

	return utils.Ptr(id.String())
}

func fromAnyBool(v interface{}) *bool {
	if v == nil {
		return nil
	}

	var b pgtype.Bool
	if err := b.Scan(v); err != nil || !b.Valid {
		return nil
	}

	return utils.Ptr(b.Bool)
}

func fromAnyInt32(v interface{}) *int32 {
	if v == nil {
		return nil
	}

	var i pgtype.Int4
	if err := i.Scan(v); err != nil || !i.Valid {
		return nil
	}

	return utils.Ptr(i.Int32)
}

func fromAnyDate(v interface{}) *time.Time {
	if v == nil {
		return nil
	}

	var d pgtype.Date
	if err := d.Scan(v); err != nil || !d.Valid {
		return nil
	}

	return utils.Ptr(d.Time)
}

func boolValue(v *bool) bool {
	if v == nil {
		return false
	}

	return *v
}
