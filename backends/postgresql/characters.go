// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package postgresql

import (
	sqlP "database/sql"
	"fmt"

	"github.com/andrew-boutin/dndtextapi/characters"
	log "github.com/sirupsen/logrus"

	sq "github.com/Masterminds/squirrel"
)

const (
	charactersTable = "characters"

	// TODO: Figure out how to use characterColumns... instead - maybe init func w/ string join
	charactersReturning = "RETURNING id, user_id, channel_id, name, description, created_on, last_updated"
)

var characterColumns = []string{
	"id",
	"user_id",
	"channel_id",
	"name",
	"description",
	"created_on",
	"last_updated",
}

func init() {
	// Add the Character table name in front of the columms to avoid ambigious references.
	for i, col := range characterColumns {
		characterColumns[i] = fmt.Sprintf("%s.%s", charactersTable, col)
	}
}

// DoesUserHaveCharacterInChannel determines if the given User has a Character in the given
// Channel.
func (backend Backend) DoesUserHaveCharacterInChannel(userID, channelID int) (bool, error) {
	sql, args, err := PSQLBuilder().
		Select("1").
		From(charactersTable).
		Where(sq.Eq{"channel_id": channelID}).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build query for does user have character in channel.")
		return false, err
	}

	// TODO: Move out into own function
	// Inspiration https://snippets.aktagon.com/snippets/756-checking-if-a-row-exists-in-go-database-sql-and-sqlx-
	sql = fmt.Sprintf("select exists(%s)", sql)

	var exists bool
	err = backend.db.QueryRow(sql, args...).Scan(&exists)
	if err != nil && err != sqlP.ErrNoRows {
		log.WithError(err).Error("Issue executing query for does user have character in channel.")
		return false, err
	}

	return exists, nil
}

// GetCharactersInChannel retrieves all of the Characters in the given Channel.
func (backend Backend) GetCharactersInChannel(channelID int) (characters.CharacterCollection, error) {
	sql, args, err := PSQLBuilder().
		Select(characterColumns...).
		From(charactersTable).
		Where(sq.Eq{"channel_id": channelID}).ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build query for get characters in channel.")
		return nil, err
	}

	rows, err := backend.db.Queryx(sql, args...)
	if err != nil {
		log.WithError(err).Error("Issue executing query for get characters in channel.")
		return nil, err
	}

	outChars := make(characters.CharacterCollection, 0)
	for rows.Next() {
		var char characters.Character
		err = rows.StructScan(&char)
		if err != nil {
			log.WithError(err).Error("Issue loading character for get characters in channel.")
			return nil, err
		}

		outChars = append(outChars, &char)
	}

	return outChars, nil
}

// CreateCharacter creates a new Character in the given Channel with the given data.
func (backend Backend) CreateCharacter(c *characters.Character) (*characters.Character, error) {
	kvs := map[string]interface{}{
		"user_id":    c.UserID,
		"channel_id": c.ChannelID,
	}

	newChar := &characters.Character{}
	err := backend.createSingle(charactersTable, charactersReturning, kvs, newChar)
	if err != nil {
		log.WithError(err).Error("Issue with create character sql.")
		return nil, err
	}

	return newChar, nil
}

// UpdateCharacter updates the Character matching the input ID using the data from
// the input Character.
func (backend Backend) UpdateCharacter(id int, c *characters.Character) (*characters.Character, error) {
	setMap := map[string]interface{}{
		"name":        c.Name,
		"description": c.Description,
	}

	updatedCharacter := &characters.Character{}
	wasFound, err := backend.updateSingle(id, channelsTable, channelsReturning, setMap, updatedCharacter)
	if err != nil {
		log.WithError(err).Error("Issue with query for update character.")
		return nil, err
	} else if !wasFound {
		return nil, characters.ErrCharacterNotFound
	}

	return updatedCharacter, nil
}

// DeleteCharacter deletes the Character matching the input ID.
func (backend Backend) DeleteCharacter(characterID int) error {
	wasFound, err := backend.deleteSingle(characterID, charactersTable)
	if err != nil {
		log.WithError(err).Error("Failed to execute delete user query.")
	} else if !wasFound {
		return characters.ErrCharacterNotFound
	}
	return err
}

// GetCharacter retrieves a single Character by ID.
func (backend Backend) GetCharacter(id int) (*characters.Character, error) {
	char := &characters.Character{}
	wasFound, err := backend.getSingle(id, channelsTable, characterColumns, char)
	if err != nil {
		log.WithError(err).Error("Query issue for get character.")
		return nil, err
	} else if !wasFound {
		return nil, characters.ErrCharacterNotFound
	}

	return char, nil
}
