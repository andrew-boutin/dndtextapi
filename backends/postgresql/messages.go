package postgresql

import (
	sqlP "database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/andrew-boutin/dndtextapi/messages"
	log "github.com/sirupsen/logrus"
)

const messagesTable = "messages"

var messageColumns = []string{
	"id",
	"userid",
	"channelid",
	"content",
	"isstory",
	"createdon",
	"lastupdated",
}

// GetMessagesInChannel retrieves all of the Messages in the database
// for the given Channel by ID. If onlyStory is nil then both msgType are returned.
// If onlyStory is set then only story Messages are returned. Otherwise only meta
// Messages are retrieved.
func (backend Backend) GetMessagesInChannel(channelID int, onlyStory *bool) (*messages.MessageCollection, error) {
	builder := PSQLBuilder().
		Select(messageColumns...).
		From(messagesTable).
		Where(sq.Eq{"channelid": channelID})

	if onlyStory != nil {
		builder = builder.Where(sq.Eq{"isstory": *onlyStory})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := backend.db.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	outMessages := make(messages.MessageCollection, 0)
	for rows.Next() {
		var message messages.Message
		err = rows.StructScan(&message)

		if err != nil {
			return nil, err
		}

		outMessages = append(outMessages, message)
	}

	return &outMessages, nil
}

// CreateMessage creates a new Message in the database using the provided data.
func (backend Backend) CreateMessage(m *messages.Message) (*messages.Message, error) {
	sql, args, err := PSQLBuilder().
		Insert(messagesTable).
		Columns("userid", "channelid", "content").
		Values(m.UserID, m.ChannelID, m.Content).
		Suffix("RETURNING id, userid, channelid, content, type, createdon, lastupdated"). // TODO: Use messageColumns...
		ToSql()
	if err != nil {
		log.WithError(err).Error("Issue building create message sql.")
		return nil, err
	}

	newMessage := &messages.Message{}
	err = backend.db.QueryRowx(sql, args...).StructScan(newMessage)
	if err != nil {
		log.WithError(err).Error("Issue running create message sql.")
		return nil, err
	}

	return newMessage, nil
}

// GetMessage retrieves the Message from the database that matches the
// given ID.
func (backend Backend) GetMessage(id int) (*messages.Message, error) {
	sql, args, err := PSQLBuilder().
		Select(messageColumns...).
		From(messagesTable).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	message := &messages.Message{}
	err = backend.db.Get(message, sql, args...)
	if err != nil {
		if err == sqlP.ErrNoRows {
			return nil, messages.ErrMessageNotFound
		}
		return nil, err
	}

	return message, nil
}

// DeleteMessagesFromChannel deletes all of the Messages in the database
// that have their Channel match the given Channel ID.
func (backend Backend) DeleteMessagesFromChannel(channelID int) error {
	sql, args, err := PSQLBuilder().
		Delete(messagesTable).
		Where(sq.Eq{"channelid": channelID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = backend.db.Exec(sql, args...)
	return err
}

// DeleteMessage deletes the Message in the database that matches
// the given ID.
func (backend Backend) DeleteMessage(id int) error {
	sql, args, err := PSQLBuilder().
		Delete(messagesTable).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := backend.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	numRowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numRowsAffected <= 0 {
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
	sql, args, err := PSQLBuilder().
		Update(messagesTable).
		SetMap(setMap).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, userid, channelid, content, type, createdon, lastupdated"). // TODO: Use messageColumns...
		ToSql()
	if err != nil {
		return nil, err
	}

	updatedMessage := &messages.Message{}
	err = backend.db.QueryRowx(sql, args...).StructScan(updatedMessage)
	if err != nil {
		if err == sqlP.ErrNoRows {
			return nil, messages.ErrMessageNotFound
		}
		return nil, err
	}

	return updatedMessage, nil
}
