// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package channels

import (
	"fmt"
	"time"
)

// ErrChannelNotFound is the error to use when the Channel is not found.
var ErrChannelNotFound = fmt.Errorf("channel not found")

// Channel contains all of the data for a channel
type Channel struct {
	Name        string    `json:"Name" db:"name"`
	Description string    `json:"Description" db:"description"`
	Topic       string    `json:"Topic" db:"topic"`
	ID          int       `json:"ID" db:"id"`
	OwnerID     int       `json:"OwnerID" db:"owner_id"`
	IsPrivate   bool      `json:"IsPrivate" db:"is_private"`
	CreatedOn   time.Time `json:"CreatedOn" db:"created_on"`
	LastUpdated time.Time `json:"LastUpdated" db:"last_updated"`
	DMID        int       `json:"DMID" db:"dm_id"`
}

// ChannelCollection is a collection of channels
type ChannelCollection []*Channel
