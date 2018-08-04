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
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Topic       string    `db:"topic"`
	ID          int       `db:"id"`
	OwnerID     int       `db:"owner_id"`
	IsPrivate   bool      `db:"is_private"`
	CreatedOn   time.Time `db:"created_on"`
	LastUpdated time.Time `db:"last_updated"`
	DMID        int       `db:"dm_id"`
}

// ChannelCollection is a collection of channels
type ChannelCollection []*Channel
