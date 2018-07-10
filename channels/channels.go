package channels

import (
	"fmt"
	"time"

	"github.com/andrew-boutin/dndtextapi/users"
)

// ErrChannelNotFound is the error to use when the Channel is not found.
var ErrChannelNotFound = fmt.Errorf("channel not found")

// Channel contains all of the data for a channel
type Channel struct {
	Name        string               `db:"name"` // TODO: Verify these tags are utilized in StructScan
	Description string               `db:"description"`
	ID          int                  `db:"id"`
	OwnerID     int                  `db:"ownerid"`
	IsPrivate   bool                 `db:"isprivate"`
	CreatedOn   time.Time            `db:"createdon"`
	LastUpdated time.Time            `db:"lastupdated"`
	DMID        int                  `db:"dmid"`
	Users       users.UserCollection `json:"users,omitempty"`
}

// ChannelCollection is a collection of channels
// TODO: Should this be a collection of pointers..?
type ChannelCollection []Channel
