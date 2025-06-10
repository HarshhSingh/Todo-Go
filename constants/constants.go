package constants

const (
	TaskStatusPending    = "pending"
	TaskStatusCompleted  = "completed"
	TaskStatusTodo       = "to-do"
	TaskStatusInProgress = "in-progress"
)

var AllTaskStatuses = []string{
	TaskStatusPending,
	TaskStatusCompleted,
	TaskStatusTodo,
	TaskStatusInProgress,
}
