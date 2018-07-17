package users

import (
	"time"
)

// User holds User info
type User struct {
	Username    string
	Email       string
	Bio         string
	ID          int
	CreatedOn   time.Time
	LastUpdated time.Time
}

// UserCollection is a slice of Users.
type UserCollection []User

// GetUserIDs retrieves a slice of just the User IDs in the collection
// of Users.
func (uc UserCollection) GetUserIDs() []int {
	ids := make([]int, len(uc))
	for i, user := range uc {
		ids[i] = user.ID
	}
	return ids
}
