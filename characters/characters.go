// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package characters

import (
	"fmt"
	"time"
)

// ErrCharacterNotFound is the error to use when the Character is not found.
var ErrCharacterNotFound = fmt.Errorf("character not found")

// Character holds all of the information that makes up a Character.
type Character struct {
	ID          int       `db:"id"`
	ChannelID   int       `db:"channel_id"`
	UserID      int       `db:"user_id"`
	Name        string    `db:"id"`
	Description string    `db:"id"`
	CreatedOn   time.Time `db:"created_on"`
	LastUpdated time.Time `db:"last_updated"`
}

// CharacterCollection is a collection of Characters
type CharacterCollection []*Character
