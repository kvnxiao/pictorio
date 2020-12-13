package chat

import (
	"sync"

	"github.com/kvnxiao/pictorio/events"
)

type Chat struct {
	mu       sync.Mutex
	messages []events.ChatEvent
}

func NewChatHistory() *Chat {
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
	c.mu.Lock()
	defer c.mu.Unlock()

	msgs := make([]events.ChatEvent, len(c.messages))
	copy(msgs, c.messages)

	return msgs
}

func (c *Chat) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = nil
}
