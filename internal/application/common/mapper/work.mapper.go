package appmapper

import (
	appdto "team_service/internal/application/common/dto"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
)

func ToWorkResponse(work *entity.Work) *appdto.WorkResponse {
	if work == nil {
		return nil
	}

	var sprintID *string
	if work.SprintID != "" {
		sprintID = &work.SprintID
	}

	var assigneeID *string
	if work.AssigneeID != "" {
		assigneeID = &work.AssigneeID
	}

	priority := enum.WorkPriority("")
	if work.Priority != nil {
		priority = *work.Priority
	}

	return &appdto.WorkResponse{
		ID:            work.ID,
		GroupID:       work.GroupID,
		SprintID:      sprintID,
		Name:          work.Name,
		Description:   work.Description,
		Status:        work.Status,
		Priority:      priority,
		AssigneeID:    assigneeID,
		CreatorID:     work.CreatorID,
		EstimateHours: work.EstimateHours,
		StoryPoint:    work.StoryPoint,
		DueDate:       work.DueDate,
		CheckList: &appdto.ChecklistSummaryResponse{
			Total:     0,
			Completed: 0,
			Items:     make([]appdto.ChecklistItemResponse, 0),
		},
		Comments: &appdto.CommentListResponse{
			Total:    0,
			Comments: make([]appdto.CommentResponse, 0),
		},
		CreatedAt: work.CreatedAt,
		UpdatedAt: work.UpdatedAt,
	}
}
