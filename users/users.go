// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package users

import (
	"fmt"
	"time"
)

// ErrUserNotFound is the error to use when the User is not found.
var ErrUserNotFound = fmt.Errorf("user not found")

// User holds User info
type User struct {
	ID          int       `json:"ID" db:"id"`
	Username    string    `json:"Username" db:"username"`
	Email       string    `json:"Email" db:"email"`
	Bio         string    `json:"Bio" db:"bio"`
	IsAdmin     bool      `json:"IsAdmin" db:"is_admin"`
	IsBanned    bool      `json:"IsBanned" db:"is_banned"`
	LastLogin   time.Time `json:"LastLogin" db:"last_login"`
	CreatedOn   time.Time `json:"CreatedOn" db:"created_on"`
	LastUpdated time.Time `json:"LastUpdated" db:"last_updated"`
}

// UserCollection is a slice of Users.
type UserCollection []*User

// GetUserIDs retrieves a slice of just the User IDs in the collection
// of Users.
func (uc UserCollection) GetUserIDs() []int {
	ids := make([]int, len(uc))
	for i, user := range uc {
		ids[i] = user.ID
	}
	return ids
}

// GoogleUser has all of the fields that we expect to come back from querying Google for User data.
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}
