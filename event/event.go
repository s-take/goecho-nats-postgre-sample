package event

import "github.com/s-take/goecho-nats-postgre-sample/schema"

// EventStore ...
type EventStore interface {
	Close()
	PublishTask(task schema.Task) error
	OnTaskPublished(f func(TaskPublishedMessage)) error
}

var impl EventStore

// SetEventStore ...
func SetEventStore(es EventStore) {
	impl = es
}

// Close ...
func Close() {
	impl.Close()
}

// PublishTask ...
func PublishTask(task schema.Task) error {
	return impl.PublishTask(task)
}

// OnTaskPublished ...
func OnTaskPublished(f func(TaskPublishedMessage)) error {
	return impl.OnTaskPublished(f)
}
