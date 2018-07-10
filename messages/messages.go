package messages

import (
	"time"
)

// Message contains message data
type Message struct {
	ID          int
	Content     string
	UserID      int
	CreatedOn   time.Time
	LastUpdated time.Time
}

// MessageCollection is a collection of messages
type MessageCollection []Message
