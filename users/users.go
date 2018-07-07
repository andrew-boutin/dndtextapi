package users

import (
	"time"
)

type User struct {
	Name        string
	Bio         string
	ID          int
	CreatedOn   time.Time
	LastUpdated time.Time // TODO: handle this
}

type UserCollection []User
