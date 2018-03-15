package model

// Config : represents application configuration
type Config struct {
	Manager ManagerInfo
}

// ManagerInfo : config fo the manager section
type ManagerInfo struct {
	ExecutionQueueSize int `toml:"execution_queue_size"`
	TaskQueueSize      int `toml:"task_queue_size"`
}
