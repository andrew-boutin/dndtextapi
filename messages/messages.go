package messages

import (
	"fmt"
	"time"
)

// ErrMessageNotFound is the error to use when the Message is not found.
var ErrMessageNotFound = fmt.Errorf("message not found")

// Message contains message data
type Message struct {
	ID          int
	Content     string
	UserID      int
	ChannelID   int
	IsStory     bool
	CreatedOn   time.Time
	LastUpdated time.Time
}

// MessageCollection is a collection of messages
type MessageCollection []Message
