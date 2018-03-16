package model

// Config : represents application configuration
type Config struct {
	Manager ManagerInfo
}

// ManagerInfo : config fo the manager section
type ManagerInfo struct {
	ExecutionQueueSize int64 `toml:"execution_queue_size"`
	TaskQueueSize      int64 `toml:"task_queue_size"`
}
