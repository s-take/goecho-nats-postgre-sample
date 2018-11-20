package event

import (
	"time"
)

// Message interface
type Message interface {
	Key() string
}

// TaskPublishedMessage struct
type TaskPublishedMessage struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

// Key ...
func (m *TaskPublishedMessage) Key() string {
	return "task.published"
}
