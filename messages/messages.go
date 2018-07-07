package messages

import (
	"time"
)

type Message struct {
	ID        int
	Content   string
	UserID    int
	CreatedOn time.Time
}

type MessageCollection []Message
