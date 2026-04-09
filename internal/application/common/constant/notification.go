package appconstant

// Event Types for the Notification System
// These constants are used to identify the type of event being processed
const (
	// Group Events (group.*)
	// Triggered by organizational and membership changes
	EventTypeGroupCreated      = "GROUP_CREATED"
	EventTypeGroupUpdated      = "GROUP_UPDATED"
	EventTypeGroupDeleted      = "GROUP_DELETED"
	EventTypeMemberJoined      = "MEMBER_JOINED"
	EventTypeMemberRemoved     = "MEMBER_REMOVED"
	EventTypeLeaveGroup        = "LEAVE_GROUP"
	EventTypeMemberRoleUpdated = "MEMBER_ROLE_UPDATED"
	EventTypeInviteCreated     = "INVITE_CREATED"
	EventTypeInviteAccepted    = "INVITE_ACCEPTED"
	EventTypeInviteError       = "INVITE_ERROR"

	// Sprint Events (sprint.*)
	// Triggered by time-management and progress-tracking changes
	EventTypeSprintCreated   = "SPRINT_CREATED"
	EventTypeSprintActivated = "SPRINT_ACTIVATED"
	EventTypeSprintCompleted = "Sprint_COMPLETED"
	EventTypeSprintCancelled = "SPRINT_CANCELLED"
	// Additional sprint events
	EventTypeSprintUpdated              = "SPRINT_UPDATED"
	EventTypeSprintDeleted              = "SPRINT_DELETED"
	EventTypeSprintGenerationSuccessful = "SPRINT_GENERATION_SUCCESSFUL"
	EventTypeSprintGenerationFailed     = "SPRINT_GENERATION_FAILED"

	// Work Events (work.*)
	// Triggered by task/ticket lifecycle changes
	EventTypeWorkCreated       = "WORK_CREATED"
	EventTypeWorkAssigned      = "WORK_ASSIGNED"
	EventTypeWorkStatusChanged = "WORK_STATUS_CHANGED"
	EventTypeWorkCommented     = "WORK_COMMENTED"
	// Additional work events
	EventTypeWorkUpdated = "WORK_UPDATED"
	EventTypeWorkDeleted = "WORK_DELETED"

	// Checklist item events
	EventTypeChecklistItemCreated = "CHECKLIST_ITEM_CREATED"
	EventTypeChecklistItemUpdated = "CHECKLIST_ITEM_UPDATED"
	EventTypeChecklistItemDeleted = "CHECKLIST_ITEM_DELETED"

	// Comment events (can reuse WORK_COMMENTED for some flows, but explicit ones help clarity)
	EventTypeCommentCreated = "COMMENT_CREATED"
	EventTypeCommentUpdated = "COMMENT_UPDATED"
	EventTypeCommentDeleted = "COMMENT_DELETED"
)

// Correlation Types
// Used to identify the base entity associated with the notification
const (
	CorrelationTypeGroup     int32 = 31
	CorrelationTypeSprint    int32 = 32
	CorrelationTypeWork      int32 = 33
	CorrelationTypeChecklist int32 = 34
	CorrelationTypeComment   int32 = 35
)

// GetDisplayTitle returns the Vietnamese UI title for a specific event type
// This title is used for Push Notification headers and Email subjects
func GetDisplayTitle(eventType string) string {
	switch eventType {
	case EventTypeGroupCreated:
		return "Tạo nhóm thành công"
	case EventTypeGroupUpdated:
		return "Cập nhật nhóm"
	case EventTypeGroupDeleted:
		return "Nhóm đã bị xóa"
	case EventTypeMemberJoined:
		return "Thành viên mới"
	case EventTypeMemberRemoved:
		return "Thay đổi nhân sự"
	case EventTypeLeaveGroup:
		return "Thay đổi nhân sự"
	case EventTypeMemberRoleUpdated:
		return "Thay đổi vai trò thành viên"
	case EventTypeInviteCreated:
		return "Bạn đã được mời tham gia nhóm"
	case EventTypeInviteAccepted:
		return "Lời mời đã được chấp nhận"
	case EventTypeInviteError:
		return "Lỗi lời mời"
	case EventTypeSprintCreated:
		return "Sprint mới đã sẵn sàng"
	case EventTypeSprintActivated:
		return "Sprint đã bắt đầu"
	case EventTypeSprintCompleted:
		return "Kết thúc Sprint"
	case EventTypeSprintCancelled:
		return "Hủy bỏ Sprint"
	case EventTypeSprintUpdated:
		return "Cập nhật Sprint"
	case EventTypeSprintDeleted:
		return "Sprint đã bị xóa"
	case EventTypeSprintGenerationSuccessful:
		return "Sprint đã được tạo bằng AI"
	case EventTypeSprintGenerationFailed:
		return "AI tạo Sprint thất bại"
	case EventTypeWorkCreated:
		return "Công việc mới"
	case EventTypeWorkAssigned:
		return "Giao việc cho bạn"
	case EventTypeWorkStatusChanged:
		return "Cập nhật trạng thái"
	case EventTypeWorkCommented:
		return "Bình luận mới"
	case EventTypeWorkUpdated:
		return "Cập nhật công việc"
	case EventTypeWorkDeleted:
		return "Công việc đã bị xóa"
	case EventTypeChecklistItemCreated:
		return "Checklist item mới"
	case EventTypeChecklistItemUpdated:
		return "Checklist item được cập nhật"
	case EventTypeChecklistItemDeleted:
		return "Checklist item đã bị xóa"
	case EventTypeCommentCreated:
		return "Bình luận mới"
	case EventTypeCommentUpdated:
		return "Bình luận đã được cập nhật"
	case EventTypeCommentDeleted:
		return "Bình luận đã bị xóa"
	default:
		return "Thông báo hệ thống"
	}
}

// IsRequireEmail determines if the event should trigger an email dispatch by default
// Based on the criticality of the event for the user
func IsRequireEmail(eventType string) bool {
	switch eventType {
	case EventTypeMemberRemoved,
		EventTypeSprintActivated,
		EventTypeSprintCompleted,
		EventTypeWorkAssigned:
		// These events are critical and require external notification
		return true
	default:
		// Other events are handled via In-App Push/WebSocket only
		return false
	}
}
