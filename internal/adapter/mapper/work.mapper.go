package adapermapper

import (
	"strconv"
	"strings"
	appdto "team_service/internal/application/common/dto"
	"team_service/internal/domain/enum"
	"team_service/proto/common"
	"team_service/proto/team_service"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToCreateWorkDTO(req *team_service.CreateWorkRequest) *appdto.CreateWorkRequest {
	if req == nil {
		return nil
	}

	return &appdto.CreateWorkRequest{
		SprintID:    req.SprintId,
		Name:        req.Name,
		Description: req.Description,
	}
}

func ToGetWorkDTO(req *common.IDRequest) *appdto.GetWorkRequest {
	if req == nil {
		return nil
	}

	return &appdto.GetWorkRequest{
		WorkID: req.Id,
	}
}

func ToListWorksDTO(req *team_service.ListWorksRequest) *appdto.ListWorksRequest {
	if req == nil {
		return nil
	}

	return &appdto.ListWorksRequest{
		SprintID:   req.SprintId,
		AssigneeID: req.AssigneeId,
	}
}

func ToUpdateWorkDTO(req *team_service.UpdateWorkRequest) *appdto.UpdateWorkRequest {
	if req == nil {
		return nil
	}

	return &appdto.UpdateWorkRequest{
		WorkID:      req.Id,
		Name:        req.Name,
		Description: req.Description,
		Status:      optionalWorkStatus(req.Status),
		SprintID:    req.SprintId,
		AssigneeID:  req.AssigneeId,
		StoryPoint:  parseOptionalInt32(req.StoryPoint),
		Priority:    optionalWorkPriority(req.Priority),
		DueDate:     optionalDateToTime(req.DueDate),
		Version:     req.Version,
	}
}

func ToDeleteWorkDTO(req *common.IDRequest) *appdto.DeleteWorkRequest {
	if req == nil {
		return nil
	}

	return &appdto.DeleteWorkRequest{
		WorkID: req.Id,
	}
}

func ToCreateChecklistItemDTO(req *team_service.CreateChecklistItemRequest) *appdto.CreateChecklistItemRequest {
	if req == nil {
		return nil
	}

	return &appdto.CreateChecklistItemRequest{
		WorkID: req.WorkId,
		Name:   req.Name,
	}
}

func ToUpdateChecklistItemDTO(req *team_service.UpdateChecklistItemRequest) *appdto.UpdateChecklistItemRequest {
	if req == nil {
		return nil
	}

	return &appdto.UpdateChecklistItemRequest{
		ItemID:      req.Id,
		Name:        req.Name,
		IsCompleted: req.IsCompleted,
	}
}

func ToDeleteChecklistItemDTO(req *common.IDRequest) *appdto.DeleteChecklistItemRequest {
	if req == nil {
		return nil
	}

	return &appdto.DeleteChecklistItemRequest{
		ItemID: req.Id,
	}
}

func ToCreateCommentDTO(req *team_service.CreateCommentRequest) *appdto.CreateCommentRequest {
	if req == nil {
		return nil
	}

	return &appdto.CreateCommentRequest{
		WorkID:  req.WorkId,
		Content: req.Content,
	}
}

func ToUpdateCommentDTO(req *team_service.UpdateCommentRequest) *appdto.UpdateCommentRequest {
	if req == nil {
		return nil
	}

	return &appdto.UpdateCommentRequest{
		CommentID: req.Id,
		Content:   req.Content,
	}
}

func ToDeleteCommentDTO(req *common.IDRequest) *appdto.DeleteCommentRequest {
	if req == nil {
		return nil
	}

	return &appdto.DeleteCommentRequest{
		CommentID: req.Id,
	}
}

func ToCreateWorkGrpcResponse(resp *appdto.BaseResponse[appdto.WorkResponse]) *team_service.CreateWorkResponse {
	if resp == nil {
		return &team_service.CreateWorkResponse{
			Work:  nil,
			Error: ToProtoError(nil),
		}
	}

	return &team_service.CreateWorkResponse{
		Work:  ToWorkMessage(resp.Data),
		Error: ToProtoError(resp.Error),
	}
}

func ToGetWorkGrpcResponse(resp *appdto.BaseResponse[appdto.WorkResponse]) *team_service.GetWorkResponse {
	if resp == nil {
		return &team_service.GetWorkResponse{
			Work:  nil,
			Error: ToProtoError(nil),
		}
	}

	return &team_service.GetWorkResponse{
		Work:  ToWorkMessage(resp.Data),
		Error: ToProtoError(resp.Error),
	}
}

func ToListWorksGrpcResponse(resp *appdto.BaseResponse[appdto.ListWorksResponse]) *team_service.ListWorksResponse {
	if resp == nil {
		return &team_service.ListWorksResponse{
			Works: nil,
			Error: ToProtoError(nil),
		}
	}

	works := make([]*team_service.WorkMessage, 0)
	if resp.Data != nil {
		works = make([]*team_service.WorkMessage, 0, len(resp.Data.Works))
		for _, work := range resp.Data.Works {
			workCopy := work
			works = append(works, ToWorkMessage(&workCopy))
		}
	}

	return &team_service.ListWorksResponse{
		Works: works,
		Error: ToProtoError(resp.Error),
	}
}

func ToUpdateWorkGrpcResponse(resp *appdto.BaseResponse[appdto.WorkResponse]) *team_service.UpdateWorkResponse {
	if resp == nil {
		return &team_service.UpdateWorkResponse{
			Work:  nil,
			Error: ToProtoError(nil),
		}
	}

	return &team_service.UpdateWorkResponse{
		Work:  ToWorkMessage(resp.Data),
		Error: ToProtoError(resp.Error),
	}
}

func ToDeleteWorkGrpcResponse(resp *appdto.BaseResponse[appdto.DeleteWorkResponse]) *team_service.DeleteWorkResponse {
	if resp == nil {
		return &team_service.DeleteWorkResponse{
			Success: false,
			Error:   ToProtoError(nil),
		}
	}

	success := false
	if resp.Data != nil {
		success = resp.Data.Success
	}

	return &team_service.DeleteWorkResponse{
		Success: success,
		Error:   ToProtoError(resp.Error),
	}
}

func ToCreateChecklistItemGrpcResponse(resp *appdto.BaseResponse[appdto.ChecklistItemResponse]) *team_service.CreateChecklistItemResponse {
	if resp == nil {
		return &team_service.CreateChecklistItemResponse{
			Item:  nil,
			Error: ToProtoError(nil),
		}
	}

	return &team_service.CreateChecklistItemResponse{
		Item:  ToChecklistItemMessage(resp.Data),
		Error: ToProtoError(resp.Error),
	}
}

func ToUpdateChecklistItemGrpcResponse(resp *appdto.BaseResponse[appdto.ChecklistItemResponse]) *team_service.UpdateChecklistItemResponse {
	if resp == nil {
		return &team_service.UpdateChecklistItemResponse{
			Checklist: nil,
			Error:     ToProtoError(nil),
		}
	}

	return &team_service.UpdateChecklistItemResponse{
		Checklist: ToChecklistItemMessage(resp.Data),
		Error:     ToProtoError(resp.Error),
	}
}

func ToDeleteChecklistItemGrpcResponse(resp *appdto.BaseResponse[appdto.ChecklistItemResponse]) *team_service.DeleteChecklistItemResponse {
	if resp == nil {
		return &team_service.DeleteChecklistItemResponse{
			Checklist: nil,
			Error:     ToProtoError(nil),
		}
	}

	return &team_service.DeleteChecklistItemResponse{
		Checklist: ToChecklistItemMessage(resp.Data),
		Error:     ToProtoError(resp.Error),
	}
}

func ToCreateCommentGrpcResponse(resp *appdto.BaseResponse[appdto.CommentListResponse]) *team_service.CreateCommentResponse {
	if resp == nil {
		return &team_service.CreateCommentResponse{
			Comment: nil,
			Error:   ToProtoError(nil),
		}
	}

	return &team_service.CreateCommentResponse{
		Comment: ToCommentListMessage(resp.Data),
		Error:   ToProtoError(resp.Error),
	}
}

func ToUpdateCommentGrpcResponse(resp *appdto.BaseResponse[appdto.CommentListResponse]) *team_service.UpdateCommentResponse {
	if resp == nil {
		return &team_service.UpdateCommentResponse{
			Comment: nil,
			Error:   ToProtoError(nil),
		}
	}

	return &team_service.UpdateCommentResponse{
		Comment: ToCommentListMessage(resp.Data),
		Error:   ToProtoError(resp.Error),
	}
}

func ToDeleteCommentGrpcResponse(resp *appdto.BaseResponse[appdto.CommentListResponse]) *team_service.DeleteCommentResponse {
	if resp == nil {
		return &team_service.DeleteCommentResponse{
			Comment: nil,
			Error:   ToProtoError(nil),
		}
	}

	return &team_service.DeleteCommentResponse{
		Comment: ToCommentListMessage(resp.Data),
		Error:   ToProtoError(resp.Error),
	}
}

func ToWorkMessage(work *appdto.WorkResponse) *team_service.WorkMessage {
	if work == nil {
		return nil
	}

	var description string
	if work.Description != nil {
		description = *work.Description
	}

	var createdAt *timestamppb.Timestamp
	if !work.CreatedAt.IsZero() {
		createdAt = timestamppb.New(work.CreatedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if !work.UpdatedAt.IsZero() {
		updatedAt = timestamppb.New(work.UpdatedAt)
	}

	return &team_service.WorkMessage{
		Id:          work.ID,
		Name:        work.Name,
		Description: description,
		Status:      MapWorkStatus(work.Status),
		Sprint:      ToSimpleSprintMessage(work.Sprint),
		Assignee:    ToSimpleUserMessage(work.Assignee),
		StoryPoint:  safeInt32(work.StoryPoint),
		DueDate:     optionalTimeToDate(work.DueDate),
		CheckList:   ToChecklistMessage(work.CheckList),
		Comments:    ToCommentListMessage(work.Comments),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Version:     work.Version,
	}
}

func ToChecklistItemMessage(item *appdto.ChecklistItemResponse) *team_service.ChecklistItemMessage {
	if item == nil {
		return nil
	}

	return &team_service.ChecklistItemMessage{
		Id:          item.ID,
		Name:        item.Name,
		IsCompleted: item.IsCompleted,
	}
}

func ToChecklistMessage(checklist *appdto.ChecklistSummaryResponse) *team_service.ChecklistMessage {
	if checklist == nil {
		return nil
	}

	items := make([]*team_service.ChecklistItemMessage, 0, len(checklist.Items))
	for _, item := range checklist.Items {
		itemCopy := item
		items = append(items, ToChecklistItemMessage(&itemCopy))
	}

	return &team_service.ChecklistMessage{
		Total:     checklist.Total,
		Completed: checklist.Completed,
		Items:     items,
	}
}

func ToCommentMessage(comment *appdto.CommentResponse) *team_service.CommentMessage {
	if comment == nil {
		return nil
	}

	var createdAt *timestamppb.Timestamp
	if !comment.CreatedAt.IsZero() {
		createdAt = timestamppb.New(comment.CreatedAt)
	}

	return &team_service.CommentMessage{
		Id:        comment.ID,
		Content:   comment.Content,
		Creator:   toSimpleUserFromComment(comment.Creator),
		CreatedAt: createdAt,
	}
}

func ToCommentListMessage(commentList *appdto.CommentListResponse) *team_service.CommentListMessage {
	if commentList == nil {
		return nil
	}

	comments := make([]*team_service.CommentMessage, 0, len(commentList.Comments))
	for _, comment := range commentList.Comments {
		commentCopy := comment
		comments = append(comments, ToCommentMessage(&commentCopy))
	}

	return &team_service.CommentListMessage{
		Total:    commentList.Total,
		Comments: comments,
	}
}

func ToSimpleSprintMessage(sprint *appdto.SimpleSprintDTO) *team_service.SimpleSprintMessage {
	if sprint == nil {
		return nil
	}

	return &team_service.SimpleSprintMessage{
		Id:   sprint.ID,
		Name: sprint.Name,
	}
}

func ToSimpleUserMessage(user *appdto.SimpleUserDTO) *team_service.SimpleUserMessage {
	if user == nil {
		return nil
	}

	return &team_service.SimpleUserMessage{
		Id:     user.ID,
		Email:  user.Email,
		Avatar: user.Avatar,
	}
}

func toSimpleUserFromComment(user appdto.UserSummaryDTO) *team_service.SimpleUserMessage {
	return &team_service.SimpleUserMessage{
		Id:     user.ID,
		Email:  user.Email,
		Avatar: user.Avatar,
	}
}

func optionalWorkPriority(priority *team_service.WorkPriority) *enum.WorkPriority {
	if priority == nil {
		return nil
	}

	mapped := MapProtoWorkPriority(*priority)
	if mapped == "" {
		return nil
	}

	return &mapped
}

func optionalWorkStatus(status *team_service.WorkStatus) *enum.WorkStatus {
	if status == nil {
		return nil
	}

	mapped := MapProtoWorkStatus(*status)
	if mapped == "" {
		return nil
	}

	return &mapped
}

func parseOptionalInt32(raw *string) *int32 {
	if raw == nil {
		return nil
	}

	text := strings.TrimSpace(*raw)
	if text == "" {
		return nil
	}

	parsed, err := strconv.ParseInt(text, 10, 32)
	if err != nil {
		return nil
	}

	v := int32(parsed)
	return &v
}

func safeInt32(v *int32) int32 {
	if v == nil {
		return 0
	}

	return *v
}

func optionalTimeToDate(t *time.Time) *team_service.Date {
	if t == nil || t.IsZero() {
		return nil
	}

	return FromTimeToDate(*t)
}
