package appmapper

import (
	appdto "team_service/internal/application/common/dto"
	"team_service/internal/domain/entity"
)

func ToSprintResponse(sprint *entity.Sprint) *appdto.SprintResponse {
	if sprint == nil {
		return nil
	}

	response := &appdto.SprintResponse{
		ID:            sprint.ID,
		GroupID:       sprint.GroupID,
		Name:          sprint.Name,
		Goal:          sprint.Goal,
		Status:        sprint.Status,
		StartDate:     sprint.StartDate,
		EndDate:       sprint.EndDate,
		CompletedWork: 0,
		CreatedAt:     sprint.CreatedAt,
		UpdatedAt:     sprint.UpdatedAt,
	}

	if sprint.VelocityWork != nil {
		response.TotalWork = *sprint.VelocityWork
	}

	if sprint.VelocityEstimate != nil {
		response.ProgressPercent = float32(*sprint.VelocityEstimate)
	}

	return response
}
