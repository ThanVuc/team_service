package repository

import (
	"context"
	"errors"
	"fmt"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/internal/infrastructure/share/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SprintRepository struct {
	q *database.Queries
}

func NewSprintRepository(
	q *database.Queries,
) *SprintRepository {
	return &SprintRepository{
		q: q,
	}
}

func (r *SprintRepository) CreateSprint(
	ctx context.Context,
	sprint *entity.Sprint,
) (*entity.Sprint, errorbase.AppError) {

	// -------- UUID mapping --------
	var id pgtype.UUID
	if err := id.Scan(sprint.ID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse sprint id"),
		)
	}

	var groupID pgtype.UUID
	if err := groupID.Scan(sprint.GroupID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	// -------- Goal (nullable text) --------
	var goal pgtype.Text
	if sprint.Goal != nil {
		goal = pgtype.Text{
			String: *sprint.Goal,
			Valid:  true,
		}
	}

	// -------- Date mapping --------
	var startDate pgtype.Date
	if err := startDate.Scan(sprint.StartDate); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse start date"),
		)
	}

	var endDate pgtype.Date
	if err := endDate.Scan(sprint.EndDate); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse end date"),
		)
	}

	// -------- Insert --------
	dbSprint, err := r.q.CreateSprint(ctx, database.CreateSprintParams{
		ID:        id,
		GroupID:   groupID,
		Name:      sprint.Name,
		Goal:      goal,
		StartDate: startDate,
		EndDate:   endDate,
	})

	if err != nil {
		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf(
				"failed to create sprint name=%s group=%s",
				sprint.Name,
				sprint.GroupID,
			)),
		)
	}

	// -------- Mapping to entity --------

	return &entity.Sprint{
		ID:               dbSprint.ID.String(),
		GroupID:          dbSprint.GroupID.String(),
		Name:             dbSprint.Name,
		Goal:             utils.Ptr(dbSprint.Goal.String),
		StartDate:        dbSprint.StartDate.Time,
		EndDate:          dbSprint.EndDate.Time,
		Status:           enum.SprintStatus(dbSprint.Status),
		VelocityWork:     utils.Ptr(int32(dbSprint.VelocityWork.Int32)),
		VelocityEstimate: utils.Ptr(float64(dbSprint.VelocityEstimate.Float64)),
		WorkDeleted:      utils.Ptr(dbSprint.WorkDeleted.Int32),
		CreatedAt:        dbSprint.CreatedAt.Time,
		UpdatedAt:        dbSprint.UpdatedAt.Time,
	}, nil
}

func (r *SprintRepository) IsSprintOverlap(
	ctx context.Context,
	groupID string,
	startDate,
	endDate time.Time,
) (bool, errorbase.AppError) {

	var gid pgtype.UUID
	if err := gid.Scan(groupID); err != nil {
		return false, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	var start pgtype.Date
	if err := start.Scan(startDate); err != nil {
		return false, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse start date"),
		)
	}

	var end pgtype.Date
	if err := end.Scan(endDate); err != nil {
		return false, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse end date"),
		)
	}

	isOverlap, err := r.q.IsSprintOverlap(ctx, database.IsSprintOverlapParams{
		GroupID: gid,
		Column2: start,
		Column3: end,
	})

	if err != nil {
		return false, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to check sprint overlap group=%s", groupID)),
		)
	}

	return isOverlap, nil
}

func (r *SprintRepository) CancelSprint(ctx context.Context, sprintID string) errorbase.AppError {
	var sprintUUID pgtype.UUID
	if err := sprintUUID.Scan(sprintID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	err := r.q.CancelActiveSprintsByGroupID(ctx, sprintUUID)
	if err != nil {
		return errorbase.Wrap(err, errdict.ErrInternal)
	}

	return nil
}

func (r *SprintRepository) DeleteDraftSprint(ctx context.Context, sprintID string) errorbase.AppError {
	var sprintUUID pgtype.UUID
	if err := sprintUUID.Scan(sprintID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	err := r.q.DeleteDraftSprintsByGroupID(ctx, sprintUUID)
	if err != nil {
		return errorbase.Wrap(err, errdict.ErrInternal)
	}

	return nil
}

func (r *SprintRepository) DeleteSprint(ctx context.Context, sprintID string) errorbase.AppError {
	var sprintUUID pgtype.UUID
	if err := sprintUUID.Scan(sprintID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse sprint id"),
		)
	}

	rowsAffected, err := r.q.DeleteSprint(ctx, sprintUUID)
	if err != nil {
		return errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to delete sprint id=%s", sprintID)),
		)
	}

	if rowsAffected == 0 {
		return errorbase.New(
			errdict.ErrUnprocessable,
			errorbase.WithDetail("sprint not found or not deletable (must be draft)"),
		)
	}

	return nil
}

func (r *SprintRepository) GetSprintByID(
	ctx context.Context,
	sprintID string,
) (*entity.Sprint, errorbase.AppError) {
	var sprintUUID pgtype.UUID
	if err := sprintUUID.Scan(sprintID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse sprint id"),
		)
	}

	dbSprint, err := r.q.GetSprintByID(ctx, sprintUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorbase.New(
				errdict.ErrNotFound,
				errorbase.WithDetail(fmt.Sprintf("sprint not found id=%s", sprintID)),
			)
		}

		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to get sprint id=%s", sprintID)),
		)
	}

	var goal *string
	if dbSprint.Goal.Valid {
		goal = utils.Ptr(dbSprint.Goal.String)
	}

	return &entity.Sprint{
		ID:               dbSprint.ID.String(),
		GroupID:          dbSprint.GroupID.String(),
		Name:             dbSprint.Name,
		Goal:             goal,
		Status:           enum.SprintStatus(dbSprint.Status),
		StartDate:        dbSprint.StartDate.Time,
		EndDate:          dbSprint.EndDate.Time,
		CreatedAt:        dbSprint.CreatedAt.Time,
		UpdatedAt:        dbSprint.UpdatedAt.Time,
		VelocityWork:     utils.Ptr(int32(dbSprint.VelocityWork.Int32)),
		VelocityEstimate: utils.Ptr(float64(dbSprint.VelocityEstimate.Float64)),
		WorkDeleted:      utils.Ptr(int32(dbSprint.WorkDeleted.Int32)),
	}, nil
}

func (r *SprintRepository) UpdateSprint(
	ctx context.Context,
	sprintID string,
	name, goal *string,
	startDate, endDate *time.Time,
) (*entity.Sprint, errorbase.AppError) {
	var sprintUUID pgtype.UUID
	if err := sprintUUID.Scan(sprintID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse sprint id"),
		)
	}

	var dbName pgtype.Text
	if name != nil {
		dbName = pgtype.Text{String: *name, Valid: true}
	}

	var dbGoal pgtype.Text
	if goal != nil {
		dbGoal = pgtype.Text{String: *goal, Valid: true}
	}

	var dbStartDate pgtype.Date
	if startDate != nil {
		if err := dbStartDate.Scan(*startDate); err != nil {
			return nil, errorbase.New(
				errdict.ErrInternal,
				errorbase.WithDetail("failed to parse start date"),
			)
		}
	}

	var dbEndDate pgtype.Date
	if endDate != nil {
		if err := dbEndDate.Scan(*endDate); err != nil {
			return nil, errorbase.New(
				errdict.ErrInternal,
				errorbase.WithDetail("failed to parse end date"),
			)
		}
	}

	updatedSprint, err := r.q.UpdateSprint(ctx, database.UpdateSprintParams{
		Name:      dbName,
		Goal:      dbGoal,
		StartDate: dbStartDate,
		EndDate:   dbEndDate,
		ID:        sprintUUID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorbase.New(
				errdict.ErrUnprocessable,
				errorbase.WithDetail("sprint not found or not editable"),
			)
		}

		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to update sprint id=%s", sprintID)),
		)
	}

	return &entity.Sprint{
		ID:               updatedSprint.ID.String(),
		GroupID:          updatedSprint.GroupID.String(),
		Name:             updatedSprint.Name,
		Goal:             utils.Ptr(updatedSprint.Goal.String),
		Status:           enum.SprintStatus(updatedSprint.Status),
		StartDate:        updatedSprint.StartDate.Time,
		EndDate:          updatedSprint.EndDate.Time,
		CreatedAt:        updatedSprint.CreatedAt.Time,
		UpdatedAt:        updatedSprint.UpdatedAt.Time,
		VelocityWork:     utils.Ptr(int32(updatedSprint.VelocityWork.Int32)),
		VelocityEstimate: utils.Ptr(float64(updatedSprint.VelocityEstimate.Float64)),
		WorkDeleted:      utils.Ptr(int32(updatedSprint.WorkDeleted.Int32)),
	}, nil
}

func (r *SprintRepository) UpdateSprintStatus(
	ctx context.Context,
	sprintID string,
	status enum.SprintStatus,
) (*entity.Sprint, errorbase.AppError) {
	var sprintUUID pgtype.UUID
	if err := sprintUUID.Scan(sprintID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse sprint id"),
		)
	}

	updatedSprint, err := r.q.UpdateSprintStatus(ctx, database.UpdateSprintStatusParams{
		Status: status.String(),
		ID:     sprintUUID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorbase.New(
				errdict.ErrNotFound,
				errorbase.WithDetail(fmt.Sprintf("sprint not found id=%s", sprintID)),
			)
		}

		return nil, errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to update sprint status id=%s", sprintID)),
		)
	}

	return &entity.Sprint{
		ID:               updatedSprint.ID.String(),
		GroupID:          updatedSprint.GroupID.String(),
		Name:             updatedSprint.Name,
		Goal:             utils.Ptr(updatedSprint.Goal.String),
		Status:           enum.SprintStatus(updatedSprint.Status),
		StartDate:        updatedSprint.StartDate.Time,
		EndDate:          updatedSprint.EndDate.Time,
		CreatedAt:        updatedSprint.CreatedAt.Time,
		UpdatedAt:        updatedSprint.UpdatedAt.Time,
		VelocityWork:     utils.Ptr(int32(updatedSprint.VelocityWork.Int32)),
		VelocityEstimate: utils.Ptr(float64(updatedSprint.VelocityEstimate.Float64)),
		WorkDeleted:      utils.Ptr(int32(updatedSprint.WorkDeleted.Int32)),
	}, nil
}

func (r *SprintRepository) GetSprintsByGroupID(
	ctx context.Context,
	sprintID string,
) ([]*entity.Sprint, errorbase.AppError) {
	var sprintUUID pgtype.UUID
	if err := sprintUUID.Scan(sprintID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	sprints, err := r.q.GetSprintsByGroupID(ctx, sprintUUID)
	if err != nil {
		return nil, errorbase.Wrap(err, errdict.ErrInternal)
	}

	var result []*entity.Sprint
	for _, sprint := range sprints {
		result = append(result, &entity.Sprint{
			ID:               sprint.ID.String(),
			GroupID:          sprint.GroupID.String(),
			Status:           enum.SprintStatus(sprint.Status),
			Name:             sprint.Name,
			Goal:             &sprint.Goal.String,
			StartDate:        sprint.StartDate.Time,
			EndDate:          sprint.EndDate.Time,
			CreatedAt:        sprint.CreatedAt.Time,
			UpdatedAt:        sprint.UpdatedAt.Time,
			VelocityWork:     utils.Ptr(int32(sprint.VelocityWork.Int32)),
			VelocityEstimate: utils.Ptr(float64(sprint.VelocityEstimate.Float64)),
			WorkDeleted:      utils.Ptr(int32(sprint.WorkDeleted.Int32)),
		})
	}

	return result, nil
}
