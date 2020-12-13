package chat

import (
	"sync"

	"github.com/kvnxiao/pictorio/events"
)

type History interface {
	Append(event events.ChatEvent)
	GetAll() []events.ChatEvent
	Clear()
}

type Chat struct {
	mu       sync.RWMutex
	messages []events.ChatEvent
}

func NewChatHistory() History {
	return &Chat{
		messages: nil,
	}
}

func (c *Chat) Append(event events.ChatEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = append(c.messages, event)
}

func (c *Chat) GetAll() []events.ChatEvent {
	c.mu.RLock()
	defer c.mu.RUnlock()

	msgs := make([]events.ChatEvent, len(c.messages))
	copy(msgs, c.messages)

	return msgs
}

func (c *Chat) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = nil
}
