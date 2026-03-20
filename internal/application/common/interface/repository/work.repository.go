package irepository

import (
	"context"
	appdto "team_service/internal/application/common/dto"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/internal/domain/entity"
)

type WorkRepository interface {
	CreateWork(ctx context.Context, work *entity.Work) (*entity.Work, errorbase.AppError)
	UpdateWork(ctx context.Context, req *appdto.UpdateWorkRequest) (*appdto.WorkResponse, errorbase.AppError)
	DeleteWork(ctx context.Context, workID string) (*appdto.DeleteWorkResponse, errorbase.AppError)
	GetWorksBySprint(ctx context.Context, groupID string, sprintID *string) ([]appdto.WorkResponse, errorbase.AppError)
	GetWorkAggregation(ctx context.Context, workID string) (*appdto.WorkResponse, errorbase.AppError)
	GetChecklistItemMeta(ctx context.Context, itemID string) (*appdto.ChecklistItemMeta, errorbase.AppError)
	GetCommentMeta(ctx context.Context, commentID string) (*appdto.CommentMeta, errorbase.AppError)

	CreateChecklistItem(ctx context.Context, item *entity.ChecklistItem) (*appdto.ChecklistItemResponse, errorbase.AppError)
	UpdateChecklistItem(ctx context.Context, req *appdto.UpdateChecklistItemRequest) (*appdto.ChecklistItemResponse, errorbase.AppError)
	DeleteChecklistItem(ctx context.Context, itemID string) (*appdto.ChecklistItemResponse, errorbase.AppError)

	CreateComment(ctx context.Context, comment *entity.Comment) (*appdto.CommentListResponse, errorbase.AppError)
	UpdateComment(ctx context.Context, req *appdto.UpdateCommentRequest) (*appdto.CommentListResponse, errorbase.AppError)
	DeleteComment(ctx context.Context, commentID string) (*appdto.CommentListResponse, errorbase.AppError)
}
