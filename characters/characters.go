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
	ID          int       `json:"ID" db:"id"`
	ChannelID   int       `json:"ChannelID" db:"channel_id"`
	UserID      int       `json:"UserID" db:"user_id"`
	Name        string    `json:"Name" db:"name"`
	Description string    `json:"Description" db:"description"`
	CreatedOn   time.Time `json:"CreatedOn" db:"created_on"`
	LastUpdated time.Time `json:"LastUpdated" db:"last_updated"`
}

// CharacterCollection is a collection of Characters
type CharacterCollection []*Character
