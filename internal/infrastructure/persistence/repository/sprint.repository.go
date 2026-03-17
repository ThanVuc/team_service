package repository

import (
	"context"
	"fmt"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/internal/infrastructure/share/utils"
	"time"

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

func (r *SprintRepository) GetSprintsByGroupID(
	ctx context.Context,
	sprintID string,
) ([]database.Sprint, errorbase.AppError) {
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

	return sprints, nil
}
