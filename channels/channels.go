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
	Name        string
	Description string
	Topic       string
	ID          int
	OwnerID     int
	IsPrivate   bool
	CreatedOn   time.Time
	LastUpdated time.Time
	DMID        int
}

// ChannelCollection is a collection of channels
type ChannelCollection []*Channel
