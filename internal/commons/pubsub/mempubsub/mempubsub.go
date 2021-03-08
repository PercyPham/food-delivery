package mempubsub

import (
	"context"
	"food-delivery/internal/commons/pubsub"
	"sync"
)

type mempubsub struct {
	eventQueue    chan pubsub.Event
	topicChansMap map[pubsub.Topic][]chan pubsub.Event
	mu            sync.RWMutex
}

// New instance of PubSub implementation
func New() pubsub.PubSub {
	ps := &mempubsub{
		eventQueue:    make(chan pubsub.Event, 10000),
		topicChansMap: make(map[pubsub.Topic][]chan pubsub.Event),
	}
	ps.run()
	return ps
}

func (ps *mempubsub) run() {
	go func() {
		for {
			e := <-ps.eventQueue
			topic := e.Topic()
			if chans, ok := ps.topicChansMap[topic]; ok {
				for i := range chans {
					go func(c chan pubsub.Event) {
						c <- e
					}(chans[i])
				}
			}
		}
	}()
}

func (ps *mempubsub) Publish(ctx context.Context, topic pubsub.Topic, event pubsub.Event) error {
	event.SetTopic(topic)
	go func() {
		ps.eventQueue <- event
	}()
	return nil
}

func (ps *mempubsub) Subscribe(ctx context.Context, topic pubsub.Topic) (<-chan pubsub.Event, func()) {
	c := make(chan pubsub.Event)

	ps.mu.Lock()

	if chans, ok := ps.topicChansMap[topic]; ok {
		chans = append(ps.topicChansMap[topic], c)
		ps.topicChansMap[topic] = chans
	} else {
		ps.topicChansMap[topic] = []chan pubsub.Event{c}
	}

	ps.mu.Unlock()

	unsubscribe := func() {
		if chans, ok := ps.topicChansMap[topic]; ok {
			for i := range chans {
				if chans[i] == c {
					chans = append(chans[:i], chans[i+1:]...) // remove chan at index i

					ps.mu.Lock()
					ps.topicChansMap[topic] = chans
					ps.mu.Unlock()
					break
				}
			}
		}
	}

	return c, unsubscribe
}
