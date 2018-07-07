package channels

import "time"

type Channel struct {
	Name        string    `db:"name"`
	Description string    `db:"description"`
	ID          int       `db:"id"`
	OwnerID     int       `db:"ownerid"`
	IsPrivate   bool      `db:"isprivate"`
	CreatedOn   time.Time `db:"createdon"`
	LastUpdated time.Time `db:"lastupdated"` // TODO: handle this
	DMID        int       `db:"dmid"`
}

// TODO: Should this be a collection of pointers..?
type ChannelCollection []Channel
