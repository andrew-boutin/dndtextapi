// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package postgresql

import (
	"fmt"

	"github.com/Masterminds/squirrel"

	sq "github.com/Masterminds/squirrel"
	"github.com/andrew-boutin/dndtextapi/messages"
	log "github.com/sirupsen/logrus"
)

const (
	messagesTable     = "messages"
	messagesReturning = "RETURNING id, character_id, channel_id, content, is_story, created_on, last_updated"
)

var messageColumns = []string{
	"id",
	"character_id",
	"channel_id",
	"content",
	"is_story",
	"created_on",
	"last_updated",
}

func init() {
	// Add the Message table name in front of the columms to avoid ambigious references.
	for i, col := range messageColumns {
		messageColumns[i] = fmt.Sprintf("%s.%s", messagesTable, col)
	}
}

// GetMessagesInChannel retrieves all of the Messages in the database
// for the given Channel by ID. If onlyStory is nil then both msgType are returned.
// If onlyStory is set then only story Messages are returned. Otherwise only meta
// Messages are retrieved.
func (backend Backend) GetMessagesInChannel(channelID int, onlyStory *bool) (messages.MessageCollection, error) {
	builder := PSQLBuilder().
		Select(messageColumns...).
		From(messagesTable).
		Where(sq.Eq{"channel_id": channelID})

	if onlyStory != nil {
		builder = builder.Where(sq.Eq{"is_story": *onlyStory})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		log.WithError(err).Error("Failed tobuild get messages in channel query.")
		return nil, err
	}

	rows, err := backend.db.Queryx(sql, args...)
	if err != nil {
		log.WithError(err).Error("Failed to execute get messages in channel query.")
		return nil, err
	}

	outMessages := make(messages.MessageCollection, 0)
	for rows.Next() {
		var message messages.Message
		err = rows.StructScan(&message)
		if err != nil {
			log.WithError(err).Error("Failed to load message from get messages in channel query.")
			return nil, err
		}

		outMessages = append(outMessages, &message)
	}

	return outMessages, nil
}

// CreateMessage creates a new Message in the database using the provided data.
func (backend Backend) CreateMessage(m *messages.Message) (*messages.Message, error) {
	kvs := map[string]interface{}{
		"character_id": m.CharacterID,
		"channel_id":   m.ChannelID,
		"content":      m.Content,
	}

	newMessage := &messages.Message{}
	err := backend.createSingle(messagesTable, messagesReturning, kvs, newMessage)
	if err != nil {
		log.WithError(err).Error("Issue with create message sql.")
		return nil, err
	}

	return newMessage, nil
}

// GetMessage retrieves the Message from the database that matches the
// given ID.
func (backend Backend) GetMessage(id int) (*messages.Message, error) {
	message := &messages.Message{}
	wasFound, err := backend.getSingle(id, messagesTable, messageColumns, message)
	if err != nil {
		return nil, err
	} else if !wasFound {
		return nil, messages.ErrMessageNotFound
	}

	return message, nil
}

// DeleteMessagesFromChannel deletes all of the Messages in the database
// that have their Channel match the given Channel ID.
func (backend Backend) DeleteMessagesFromChannel(channelID int) error {
	sql, args, err := PSQLBuilder().
		Delete(messagesTable).
		Where(sq.Eq{"channel_id": channelID}).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build delete messages from channel query.")
		return err
	}

	_, err = backend.db.Exec(sql, args...)
	if err != nil {
		log.WithError(err).Error("Failed to execute delete messages from channel query.")
	}
	return err
}

// DeleteMessage deletes the Message in the database that matches
// the given ID.
func (backend Backend) DeleteMessage(id int) error {
	wasFound, err := backend.deleteSingle(id, messagesTable)
	if err != nil {
		log.WithError(err).Error("Failed to execute delete message query.")
	} else if !wasFound {
		return messages.ErrMessageNotFound
	}
	return err
}

// UpdateMessage updates the Message in the database matching the input ID
// with the data from the given Message.
func (backend Backend) UpdateMessage(id int, m *messages.Message) (*messages.Message, error) {
	setMap := map[string]interface{}{
		"content": m.Content,
	}

	updatedMessage := &messages.Message{}
	wasFound, err := backend.updateSingle(id, messagesTable, messagesReturning, setMap, updatedMessage)
	if err != nil {
		log.WithError(err).Error("Issue with query for update message.")
		return nil, err
	} else if !wasFound {
		return nil, messages.ErrMessageNotFound
	}

	return updatedMessage, nil
}

// DeleteMessagesFromUser deletes all of the messages that were from the input
// User. This means that the Messages are from a Character that is the User's.
func (backend Backend) DeleteMessagesFromUser(userID int) error {
	findMessagesQuery := fmt.Sprintf("SELECT messages.id FROM messages INNER JOIN "+
		"characters characters.id ON messages.character_id WHERE characters.user_id = %d", userID)

	sql, args, err := PSQLBuilder().
		Delete(messagesTable).
		Where(fmt.Sprintf("id IN (%s)", findMessagesQuery)).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build delete messages from user sql.")
		return err
	}

	_, err = backend.db.Exec(sql, args...)
	if err != nil {
		log.WithError(err).Error("Failed to execute delete messages from user query.")
	}
	return err
}

// DeleteMessagesFromCharacter deletes all of the messages that match the input
// Character ID.
func (backend Backend) DeleteMessagesFromCharacter(characterID int) error {
	sql, args, err := PSQLBuilder().
		Delete(messagesTable).
		Where(squirrel.Eq{"character_id": characterID}).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build delete messages from character sql.")
		return err
	}

	_, err = backend.db.Exec(sql, args...)
	if err != nil {
		log.WithError(err).Error("Failed to execute delete messages from character query.")
	}
	return err
}
