package subscription

import (
	"context"
	"sync"

	"backend/internal/graph/model"
)

// EventType represents the type of subscription event
type EventType string

const (
	PostAddedEvent    EventType = "POST_ADDED"
	PostUpdatedEvent  EventType = "POST_UPDATED"
	CommentAddedEvent EventType = "COMMENT_ADDED"
)

// Event represents a subscription event
type Event struct {
	Type    EventType
	PostID  string
	Post    *model.Post
	Comment *model.Comment
}

// Subscriber represents a subscription channel
type Subscriber struct {
	ID      string
	Channel chan *Event
	Filter  func(*Event) bool
}

// Manager handles GraphQL subscriptions
type Manager struct {
	subscribers map[string]*Subscriber
	mutex       sync.RWMutex
}

// NewManager creates a new subscription manager
func NewManager() *Manager {
	return &Manager{
		subscribers: make(map[string]*Subscriber),
	}
}

// Subscribe adds a new subscriber
func (m *Manager) Subscribe(ctx context.Context, id string, filter func(*Event) bool) <-chan *Event {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ch := make(chan *Event, 10) // Buffer to prevent blocking
	subscriber := &Subscriber{
		ID:      id,
		Channel: ch,
		Filter:  filter,
	}

	m.subscribers[id] = subscriber

	// Clean up when context is cancelled
	go func() {
		<-ctx.Done()
		m.Unsubscribe(id)
	}()

	return ch
}

// Unsubscribe removes a subscriber
func (m *Manager) Unsubscribe(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if subscriber, exists := m.subscribers[id]; exists {
		close(subscriber.Channel)
		delete(m.subscribers, id)
	}
}

// Publish sends an event to all matching subscribers
func (m *Manager) Publish(event *Event) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, subscriber := range m.subscribers {
		// Apply filter if provided
		if subscriber.Filter != nil && !subscriber.Filter(event) {
			continue
		}

		// Non-blocking send
		select {
		case subscriber.Channel <- event:
		default:
			// Channel is full, skip this subscriber
			// In production, you might want to log this or handle it differently
		}
	}
}

// PublishPostAdded publishes a post added event
func (m *Manager) PublishPostAdded(post *model.Post) {
	m.Publish(&Event{
		Type: PostAddedEvent,
		Post: post,
	})
}

// PublishPostUpdated publishes a post updated event
func (m *Manager) PublishPostUpdated(post *model.Post) {
	m.Publish(&Event{
		Type:   PostUpdatedEvent,
		PostID: post.ID.String(),
		Post:   post,
	})
}

// PublishCommentAdded publishes a comment added event
func (m *Manager) PublishCommentAdded(comment *model.Comment) {
	m.Publish(&Event{
		Type:    CommentAddedEvent,
		PostID:  comment.PostID.String(),
		Comment: comment,
	})
}

// GetSubscriberCount returns the number of active subscribers
func (m *Manager) GetSubscriberCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.subscribers)
}