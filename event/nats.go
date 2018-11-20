package event

import (
	"bytes"
	"encoding/gob"

	"github.com/nats-io/go-nats-streaming"
	"github.com/s-take/goecho-nats-postgre-sample/schema"
)

// StanEventStore ...
type StanEventStore struct {
	sc                        stan.Conn
	taskPublishedSubscription stan.Subscription
}

// NewStan ...
func NewStan(url string, cluster string, client string) (*StanEventStore, error) {
	sc, err := stan.Connect(cluster, client, stan.NatsURL(url))
	if err != nil {
		return nil, err
	}
	return &StanEventStore{sc: sc}, nil
}

// OnTaskPublished ...
func (e *StanEventStore) OnTaskPublished(f func(TaskPublishedMessage)) (err error) {
	m := TaskPublishedMessage{}
	qgroup := "my-queue"
	durable := "my-durable"
	inflight := 1
	//aw, _ := time.ParseDuration("60s")
	e.taskPublishedSubscription, err = e.sc.QueueSubscribe(m.Key(), qgroup, func(msg *stan.Msg) {
		e.readMessage(msg.Data, &m)
		f(m)
	}, stan.DurableName(durable), stan.MaxInflight(inflight))
	//}, stan.DurableName(durable), stan.MaxInflight(inflight), stan.SetManualAckMode(), stan.AckWait(aw))
	return
}

// Close ...
func (e *StanEventStore) Close() {
	if e.sc != nil {
		e.sc.Close()
	}
	if e.taskPublishedSubscription != nil {
		e.taskPublishedSubscription.Unsubscribe()
	}
}

// PublishTask ...
func (e *StanEventStore) PublishTask(task schema.Task) error {
	m := TaskPublishedMessage{task.ID, task.Name, task.CreatedAt}
	data, err := e.writeMessage(&m)
	if err != nil {
		return err
	}
	return e.sc.Publish(m.Key(), data)
}

// writeMessage ...
func (e *StanEventStore) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// readMessage ...
func (e *StanEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}
