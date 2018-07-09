package channels

import (
	"time"

	"github.com/andrew-boutin/dndtextapi/users"
)

type Channel struct {
	Name        string               `db:"name"` // TODO: Verify these tags are utilized in StructScan
	Description string               `db:"description"`
	ID          int                  `db:"id"`
	OwnerID     int                  `db:"ownerid"`
	IsPrivate   bool                 `db:"isprivate"`
	CreatedOn   time.Time            `db:"createdon"`
	LastUpdated time.Time            `db:"lastupdated"` // TODO: handle this
	DMID        int                  `db:"dmid"`
	Users       users.UserCollection `json:"users,omitempty"`
}

// TODO: Should this be a collection of pointers..?
type ChannelCollection []Channel
