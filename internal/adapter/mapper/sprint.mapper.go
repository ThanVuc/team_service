package adapermapper

import (
	appconstant "team_service/internal/application/common/constant"
	appdto "team_service/internal/application/common/dto"
	"team_service/proto/common"
	"team_service/proto/team_service"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToCreateSprintDTO(req *team_service.CreateSprintRequest) *appdto.CreateSprintRequest {
	if req == nil {
		return nil
	}

	return &appdto.CreateSprintRequest{
		GroupID:   req.GroupId,
		Name:      req.Name,
		Goal:      req.Goal,
		StartDate: FromDateToTime(req.StartDate),
		EndDate:   FromDateToTime(req.EndDate),
	}
}

func ToGetSprintDTO(req *common.IDRequest) *appdto.GetSprintRequest {
	if req == nil {
		return nil
	}

	return &appdto.GetSprintRequest{
		SprintID: req.Id,
	}
}

func ToListSprintsDTO(req *team_service.ListSprintsRequest) *appdto.ListSprintsRequest {
	if req == nil {
		return nil
	}

	return &appdto.ListSprintsRequest{
		GroupID: req.GroupId,
	}
}

func ToGetSimpleSprintsDTO(req *common.IDRequest) *appdto.ListSprintsRequest {
	if req == nil {
		return nil
	}

	return &appdto.ListSprintsRequest{
		GroupID: req.Id,
	}
}

func ToUpdateSprintDTO(req *team_service.UpdateSprintRequest) *appdto.UpdateSprintRequest {
	if req == nil {
		return nil
	}

	return &appdto.UpdateSprintRequest{
		SprintID:  req.Id,
		Name:      req.Name,
		Goal:      req.Goal,
		StartDate: optionalDateToTime(req.StartDate),
		EndDate:   optionalDateToTime(req.EndDate),
	}
}

func ToUpdateSprintStatusDTO(req *team_service.UpdateSprintStatusRequest) *appdto.UpdateSprintStatusRequest {
	if req == nil {
		return nil
	}

	return &appdto.UpdateSprintStatusRequest{
		SprintID: req.Id,
		Status:   MapProtoSprintStatus(req.Status),
	}
}

func ToDeleteSprintDTO(req *common.IDRequest) *appdto.DeleteSprintRequest {
	if req == nil {
		return nil
	}

	return &appdto.DeleteSprintRequest{
		SprintID: req.Id,
	}
}

func ToExportSprintDTO(req *common.IDRequest) *appdto.ExportSprintRequest {
	if req == nil {
		return nil
	}

	return &appdto.ExportSprintRequest{
		SprintID: req.Id,
	}
}

func ToGenerateSprintDTO(req *team_service.AISprintGenerationRequest) *appdto.GenerateSprintRequest {
	if req == nil {
		return nil
	}

	files := make([]appdto.AISprintGenerationFile, 0, len(req.Files))
	for _, file := range req.Files {
		if file == nil {
			continue
		}

		files = append(files, appdto.AISprintGenerationFile{
			ObjectKey: file.ObjectKey,
			Size:      file.Size,
		})
	}

	return &appdto.GenerateSprintRequest{
		Name:              req.Name,
		Goal:              req.Goal,
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		AdditionalContext: req.AdditionalContext,
		Files:             files,
	}
}

func ToCreateSprintGrpcResponse(resp *appdto.BaseResponse[appdto.SprintResponse]) *team_service.CreateSprintResponse {
	if resp == nil {
		return &team_service.CreateSprintResponse{
			Sprint: nil,
			Error:  ToProtoError(nil),
		}
	}

	return &team_service.CreateSprintResponse{
		Sprint: ToSprintMessage(resp.Data),
		Error:  ToProtoError(resp.Error),
	}
}

func ToGetSprintGrpcResponse(resp *appdto.BaseResponse[appdto.SprintResponse]) *team_service.GetSprintResponse {
	if resp == nil {
		return &team_service.GetSprintResponse{
			Sprint: nil,
			Error:  ToProtoError(nil),
		}
	}

	return &team_service.GetSprintResponse{
		Sprint: ToSprintMessage(resp.Data),
		Error:  ToProtoError(resp.Error),
	}
}

func ToListSprintsGrpcResponse(resp *appdto.BaseResponse[appdto.ListSprintsResponse]) *team_service.ListSprintsResponse {
	if resp == nil {
		return &team_service.ListSprintsResponse{
			Sprints: nil,
			Total:   0,
			Error:   ToProtoError(nil),
		}
	}

	messages := make([]*team_service.SprintMessage, 0)
	if resp.Data != nil {
		messages = make([]*team_service.SprintMessage, 0, len(resp.Data.Sprints))
		for _, sprint := range resp.Data.Sprints {
			sprintCopy := sprint
			messages = append(messages, ToSprintMessage(&sprintCopy))
		}
	}

	total := int32(0)
	if resp.Data != nil {
		total = resp.Data.Total
	}

	return &team_service.ListSprintsResponse{
		Sprints: messages,
		Total:   total,
		Error:   ToProtoError(resp.Error),
	}
}

func ToGetSimpleSprintsGrpcResponse(resp *appdto.BaseResponse[[]appdto.SimpleSprintResponse]) *team_service.GetSimpleSprintsResponse {
	if resp == nil {
		return &team_service.GetSimpleSprintsResponse{
			Sprints: nil,
			Error:   ToProtoError(nil),
		}
	}

	messages := make([]*team_service.SimpleSprintMessage, 0)
	if resp.Data != nil {
		messages = make([]*team_service.SimpleSprintMessage, 0, len(*resp.Data))
		for _, sprint := range *resp.Data {
			sprintCopy := sprint
			status := MapSprintStatus(sprintCopy.Status)
			messages = append(messages, &team_service.SimpleSprintMessage{
				Id:     sprintCopy.ID,
				Name:   sprintCopy.Name,
				Status: &status,
			})
		}
	}

	return &team_service.GetSimpleSprintsResponse{
		Sprints: messages,
		Error:   ToProtoError(resp.Error),
	}
}

func ToUpdateSprintGrpcResponse(resp *appdto.BaseResponse[appdto.SprintResponse]) *team_service.UpdateSprintResponse {
	if resp == nil {
		return &team_service.UpdateSprintResponse{
			Sprint: nil,
			Error:  ToProtoError(nil),
		}
	}

	return &team_service.UpdateSprintResponse{
		Sprint: ToSprintMessage(resp.Data),
		Error:  ToProtoError(resp.Error),
	}
}

func ToUpdateSprintStatusGrpcResponse(resp *appdto.BaseResponse[appdto.UpdateSprintStatusResponse]) *team_service.UpdateSprintStatusResponse {
	if resp == nil {
		return &team_service.UpdateSprintStatusResponse{
			Id:     "",
			Status: team_service.SprintStatus_SPRINT_STATUS_UNSPECIFIED,
			Error:  ToProtoError(nil),
		}
	}

	var id string
	status := team_service.SprintStatus_SPRINT_STATUS_UNSPECIFIED
	if resp.Data != nil {
		id = resp.Data.SprintID
		status = MapSprintStatus(resp.Data.Status)
	}

	return &team_service.UpdateSprintStatusResponse{
		Id:     id,
		Status: status,
		Error:  ToProtoError(resp.Error),
	}
}

func ToDeleteSprintGrpcResponse(resp *appdto.BaseResponse[appdto.DeleteSprintResponse]) *team_service.DeleteSprintResponse {
	if resp == nil {
		return &team_service.DeleteSprintResponse{
			Success: false,
			Error:   ToProtoError(nil),
		}
	}

	success := false
	if resp.Data != nil {
		success = resp.Data.Success
	}

	return &team_service.DeleteSprintResponse{
		Success: success,
		Error:   ToProtoError(resp.Error),
	}
}

func ToExportSprintGrpcResponse(resp *appdto.BaseResponse[appdto.ExportSprintResponse]) *team_service.ExportSprintResponse {
	if resp == nil {
		return &team_service.ExportSprintResponse{
			File:        nil,
			Filename:    "",
			ContentType: appconstant.SprintExportContentType,
			Error:       ToProtoError(nil),
		}
	}

	file := make([]byte, 0)
	filename := ""
	contentType := appconstant.SprintExportContentType

	if resp.Data != nil {
		file = resp.Data.File
		filename = resp.Data.FileName
		if resp.Data.ContentType != "" {
			contentType = resp.Data.ContentType
		}
	}

	return &team_service.ExportSprintResponse{
		File:        file,
		Filename:    filename,
		ContentType: contentType,
		Error:       ToProtoError(resp.Error),
	}
}

func ToGenerateSprintGrpcResponse(resp *appdto.BaseResponse[appdto.GenerateSprintResponse]) *team_service.AISprintGenerationResponse {
	if resp == nil {
		return &team_service.AISprintGenerationResponse{
			Message: "",
			Error:   ToProtoError(nil),
		}
	}

	message := ""
	if resp.Data != nil {
		message = resp.Data.Message
	}

	return &team_service.AISprintGenerationResponse{
		Message: message,
		Error:   ToProtoError(resp.Error),
	}
}

func ToSprintMessage(sprint *appdto.SprintResponse) *team_service.SprintMessage {
	if sprint == nil {
		return nil
	}

	goal := ""
	if sprint.Goal != nil {
		goal = *sprint.Goal
	}

	var createdAt *timestamppb.Timestamp
	if !sprint.CreatedAt.IsZero() {
		createdAt = timestamppb.New(sprint.CreatedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if !sprint.UpdatedAt.IsZero() {
		updatedAt = timestamppb.New(sprint.UpdatedAt)
	}

	return &team_service.SprintMessage{
		Id:              sprint.ID,
		GroupId:         sprint.GroupID,
		Name:            sprint.Name,
		Goal:            goal,
		Status:          MapSprintStatus(sprint.Status),
		StartDate:       FromTimeToDate(sprint.StartDate),
		EndDate:         FromTimeToDate(sprint.EndDate),
		TotalWork:       sprint.TotalWork,
		CompletedWork:   sprint.CompletedWork,
		ProgressPercent: sprint.ProgressPercent,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

func optionalDateToTime(date *team_service.Date) *time.Time {
	if date == nil {
		return nil
	}

	t := FromDateToTime(date)
	return &t
}

func ToDeleteDraftSprintDTO(req *common.IDRequest) *appdto.DeleteSprintRequest {
	if req == nil {
		return nil
	}

	return &appdto.DeleteSprintRequest{
		SprintID: req.Id,
	}
}

func ToDeleteDraftSprintGrpcResponse(resp *appdto.BaseResponse[appdto.DeleteSprintResponse]) *team_service.DeleteDraftSprintsResponse {
	if resp == nil {
		return &team_service.DeleteDraftSprintsResponse{
			Success: false,
			Error:   ToProtoError(nil),
		}
	}

	success := false
	if resp.Data != nil {
		success = resp.Data.Success
	}

	return &team_service.DeleteDraftSprintsResponse{
		Success: success,
		Error:   ToProtoError(resp.Error),
	}
}
