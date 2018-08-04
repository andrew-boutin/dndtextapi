// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package messages

import (
	"fmt"
	"time"
)

// ErrMessageNotFound is the error to use when the Message is not found.
var ErrMessageNotFound = fmt.Errorf("message not found")

// Message contains message data
type Message struct {
	ID          int       `db:"id"`
	Content     string    `db:"content"`
	CharacterID int       `db:"character_id"`
	ChannelID   int       `db:"channel_id"`
	IsStory     bool      `db:"is_story"`
	CreatedOn   time.Time `db:"created_on"`
	LastUpdated time.Time `db:"last_updated"`
}

// MessageCollection is a collection of messages
type MessageCollection []*Message
