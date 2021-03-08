package pubsub

import (
	"context"
	"fmt"
	"time"
)

// Topic string
type Topic string

// Event interface
type Event interface {
	SetTopic(topic Topic) // Topic setter

	ID() string           // ID getter
	Topic() Topic         // Topic getter
	Data() interface{}    // Data getter
	CreatedAt() time.Time // CreatedAt getter

	String() string
}

// PubSub interface for publish/subscribe pattern
type PubSub interface {
	// Publish an event to topic
	Publish(ctx context.Context, topic Topic, event Event) error

	// Subscribe returns a channel to receive events from specified topic and an unsubscribe method
	Subscribe(ctx context.Context, topic Topic) (c <-chan Event, unsubscribe func())
}

/////////////////////////////////////////////////
/////// Event code //////////////////////////////
/////////////////////////////////////////////////

type event struct {
	id        string
	topic     Topic
	data      interface{}
	createdAt time.Time
}

// NewEvent creates a new event with specified data
func NewEvent(data interface{}) Event {
	now := time.Now().UTC()
	return &event{
		id:        fmt.Sprintf("%d", now.Nanosecond()),
		data:      data,
		createdAt: now,
	}
}

func (e *event) SetTopic(topic Topic) { e.topic = topic }

func (e *event) ID() string           { return e.id }
func (e *event) Topic() Topic         { return e.topic }
func (e *event) Data() interface{}    { return e.data }
func (e *event) CreatedAt() time.Time { return e.createdAt }

func (e *event) String() string { return fmt.Sprintf("Event %s", e.topic) }
