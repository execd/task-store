package model

type Config struct {
	Manager ManagerInfo
}

type ManagerInfo struct {
	ExecutionQueueSize int `toml:"execution_queue_size"`
	TaskQueueSize      int `toml:"task_queue_size"`
}
