package enum

type WorkPriority string

const (
	WorkPriorityLow    WorkPriority = "low"
	WorkPriorityMedium WorkPriority = "medium"
	WorkPriorityHigh   WorkPriority = "high"
)

func (wp WorkPriority) IsValid() bool {
	switch wp {
	case WorkPriorityLow,
		WorkPriorityMedium,
		WorkPriorityHigh:
		return true
	}
	return false
}
