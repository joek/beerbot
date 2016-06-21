package ws

// BotCommand command input for robot
type BotCommand struct {
	Motor *Motor `json:"motor,omitempty"`
	Event string `json:"command,omitempty"`
}

// Motor command
type Motor struct {
	Left  float32 `json:"left"`
	Right float32 `json:"right"`
}
