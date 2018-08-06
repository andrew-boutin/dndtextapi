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
	ID          int       `json:"ID" db:"id"`
	Content     string    `json:"Content" db:"content"`
	CharacterID int       `json:"CharacterID" db:"character_id"`
	ChannelID   int       `json:"ChannelID" db:"channel_id"`
	IsStory     bool      `json:"IsStory" db:"is_story"`
	CreatedOn   time.Time `json:"CreatedOn" db:"created_on"`
	LastUpdated time.Time `json:"LastUpdated" db:"last_updated"`
}

// MessageCollection is a collection of messages
type MessageCollection []*Message
