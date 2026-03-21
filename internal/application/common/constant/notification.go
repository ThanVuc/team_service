package appconstant

// Event Types for the Notification System
// These constants are used to identify the type of event being processed
const (
	// Group Events (group.*)
	// Triggered by organizational and membership changes
	EventTypeGroupCreated  = "GROUP_CREATED"
	EventTypeMemberJoined  = "MEMBER_JOINED"
	EventTypeMemberRemoved = "MEMBER_REMOVED"

	// Sprint Events (sprint.*)
	// Triggered by time-management and progress-tracking changes
	EventTypeSprintCreated   = "SPRINT_CREATED"
	EventTypeSprintActivated = "SPRINT_ACTIVATED"
	EventTypeSprintCompleted = "SPRINT_COMPLETED"
	EventTypeSprintCancelled = "SPRINT_CANCELLED"

	// Work Events (work.*)
	// Triggered by task/ticket lifecycle changes
	EventTypeWorkCreated       = "WORK_CREATED"
	EventTypeWorkAssigned      = "WORK_ASSIGNED"
	EventTypeWorkStatusChanged = "WORK_STATUS_CHANGED"
	EventTypeWorkCommented     = "WORK_COMMENTED"
)

// Correlation Types
// Used to identify the base entity associated with the notification
const (
	CorrelationTypeGroup  int32 = 31
	CorrelationTypeSprint int32 = 32
	CorrelationTypeWork   int32 = 33
)

// GetDisplayTitle returns the Vietnamese UI title for a specific event type
// This title is used for Push Notification headers and Email subjects
func GetDisplayTitle(eventType string) string {
	switch eventType {
	case EventTypeGroupCreated:
		return "Tạo nhóm thành công"
	case EventTypeMemberJoined:
		return "Thành viên mới"
	case EventTypeMemberRemoved:
		return "Thay đổi nhân sự"
	case EventTypeSprintCreated:
		return "Sprint mới đã sẵn sàng"
	case EventTypeSprintActivated:
		return "Sprint đã bắt đầu"
	case EventTypeSprintCompleted:
		return "Kết thúc Sprint"
	case EventTypeSprintCancelled:
		return "Hủy bỏ Sprint"
	case EventTypeWorkCreated:
		return "Công việc mới"
	case EventTypeWorkAssigned:
		return "Giao việc cho bạn"
	case EventTypeWorkStatusChanged:
		return "Cập nhật trạng thái"
	case EventTypeWorkCommented:
		return "Bình luận mới"
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
