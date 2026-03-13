package enum

type WorkStatus string

const (
	WorkStatusTodo       WorkStatus = "todo"
	WorkStatusInProgress WorkStatus = "inprogress"
	WorkStatusInReview   WorkStatus = "inreview"
	WorkStatusDone       WorkStatus = "done"
)

func (s WorkStatus) IsValid() bool {
	switch s {
	case WorkStatusTodo,
		WorkStatusInProgress,
		WorkStatusInReview,
		WorkStatusDone:
		return true
	}
	return false
}
