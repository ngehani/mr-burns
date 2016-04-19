package controller

type TaskState int

const (
	TASK_STATE_WAITING TaskState = 1
	TASK_STATE_RUNNING TaskState = 2
	TASK_STATE_DONE TaskState = 3
)

type Task struct {
	NextRuntimeMillisecond int64
	State                  TaskState
	Data                   interface{}
}
