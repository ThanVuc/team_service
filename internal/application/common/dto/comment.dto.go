package appdto

import "time"

type CreateCommentRequest struct {
	WorkID  string
	Content string
}

type UpdateCommentRequest struct {
	CommentID string
	Content   string
}

type DeleteCommentRequest struct {
	CommentID string
}

type CommentResponse struct {
	ID      string
	Content string

	Creator UserSummaryDTO

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentMeta struct {
	ID        string
	WorkID    string
	CreatorID string
}

type CommentListResponse struct {
	Total    int32
	Comments []CommentResponse
}
