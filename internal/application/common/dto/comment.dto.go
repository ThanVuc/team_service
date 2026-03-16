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
