package enum

type SprintStatus string

const (
	SprintStatusDraft     SprintStatus = "draft"
	SprintStatusActive    SprintStatus = "active"
	SprintStatusCompleted SprintStatus = "completed"
	SprintStatusCancelled SprintStatus = "cancelled"
)

func (s SprintStatus) IsValid() bool {
	switch s {
	case SprintStatusDraft,
		SprintStatusActive,
		SprintStatusCompleted,
		SprintStatusCancelled:
		return true
	}
	return false
}
