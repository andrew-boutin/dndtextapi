package characters

import (
	"fmt"
	"time"
)

// ErrCharacterNotFound is the error to use when the Character is not found.
var ErrCharacterNotFound = fmt.Errorf("character not found")

// Character holds all of the information that makes up a Character.
type Character struct {
	ID          int
	ChannelID   int
	UserID      int
	Name        string
	Description string
	CreatedOn   time.Time
	LastUpdated time.Time
}

// CharacterCollection is a collection of Characters
type CharacterCollection []*Character
