// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package postgresql

import (
	sqlP "database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/andrew-boutin/dndtextapi/channels"
	log "github.com/sirupsen/logrus"
)

const (
	channelsTable     = "channels"
	channelsReturning = "RETURNING id, name, description, topic, owner_id, is_private, dm_id, created_on, last_updated"
)

var channelColumns = []string{
	"id",
	"name",
	"description",
	"topic",
	"owner_id",
	"is_private",
	"dm_id",
	"created_on",
	"last_updated",
}

func init() {
	// Add the Channel table name in front of the columms to avoid ambigious references.
	for i, col := range channelColumns {
		channelColumns[i] = fmt.Sprintf("%s.%s", channelsTable, col)
	}
}

// GetChannel retrieves the channel corresponding to the given id.
func (backend Backend) GetChannel(id int) (*channels.Channel, error) {
	sql, args, err := PSQLBuilder().
		Select(channelColumns...).
		From(channelsTable).
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	channel := &channels.Channel{}
	err = backend.db.Get(channel, sql, args...)
	if err != nil {
		if err == sqlP.ErrNoRows {
			return nil, channels.ErrChannelNotFound
		}
		return nil, err
	}

	return channel, err
}

// GetChannelsOwnedByUser retrieves all of the Channels where the provided User ID
// is the owner of the Channel.
func (backend Backend) GetChannelsOwnedByUser(userID int) (channels.ChannelCollection, error) {
	sql, args, err := PSQLBuilder().
		Select(channelColumns...).
		From(channelsTable).
		Where(sq.Eq{"owner_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	return backend.runMultiChannelQuery(sql, args)
}

// GetAllChannels returns a list of all Channels if the isPrivate flag is nil. If the flag is set then only
// private Channels are returned. If the flag is not set then only public Channels are returned.
func (backend Backend) GetAllChannels(isPrivate *bool) (channels.ChannelCollection, error) {
	builder := PSQLBuilder().
		Select(channelColumns...).
		From(channelsTable)

	if isPrivate != nil {
		builder = builder.Where(sq.Eq{"is_private": *isPrivate})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	return backend.runMultiChannelQuery(sql, args)
}

// GetChannelsUserHasCharacterIn finds all of the Channels that the given User has at least one Character in.
func (backend Backend) GetChannelsUserHasCharacterIn(userID int, isPrivate *bool) (channels.ChannelCollection, error) {
	builder := PSQLBuilder().
		Select(channelColumns...).
		Distinct().
		From(channelsTable).
		Join(fmt.Sprintf("%s ON %s.%s = %s.%s", charactersTable, charactersTable, "channel_id", channelsTable, "id")).
		Where(sq.Eq{"user_id": userID})

	if isPrivate != nil {
		builder = builder.Where(sq.Eq{"is_private": *isPrivate})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build query to find channels that the user has at least one character in.")
		return nil, err
	}

	return backend.runMultiChannelQuery(sql, args)
}

func (backend Backend) runMultiChannelQuery(sql string, args []interface{}) (channels.ChannelCollection, error) {
	rows, err := backend.db.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	outChannels := make(channels.ChannelCollection, 0)
	for rows.Next() {
		var channel channels.Channel
		err = rows.StructScan(&channel)
		if err != nil {
			return nil, err
		}

		outChannels = append(outChannels, &channel)
	}

	return outChannels, nil
}

// CreateChannel creates a new channel using the provided channel info
// and returns the result from the database.
func (backend Backend) CreateChannel(c *channels.Channel, userID int) (*channels.Channel, error) {
	// TODO: Don't require description, default isprivate to false
	sql, args, err := PSQLBuilder().
		Insert(channelsTable).
		Columns("name", "description", "topic", "owner_id", "is_private", "dm_id").
		Values(c.Name, c.Description, c.OwnerID, c.IsPrivate, c.DMID).
		Suffix(channelsReturning).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Issue building create channel sql.")
		return nil, err
	}

	newChannel := &channels.Channel{}
	err = backend.db.QueryRowx(sql, args...).StructScan(newChannel)
	if err != nil {
		log.WithError(err).Error("Issue running create channel sql.")
		return nil, err
	}

	return newChannel, nil
}

// DeleteChannel deletes the channel that corresponds to the given ID.
func (backend Backend) DeleteChannel(id int) error {
	sql, args, err := PSQLBuilder().
		Delete(channelsTable).
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
		return channels.ErrChannelNotFound
	}

	return err
}

// UpdateChannel updates the channel matching the given ID using the data
// provided in the input channel. Returns the channel data from the database.
func (backend Backend) UpdateChannel(id int, c *channels.Channel) (*channels.Channel, error) {
	setMap := map[string]interface{}{
		"name":        c.Name,
		"description": c.Description,
		"topic":       c.Topic,
		"owner_id":    c.OwnerID,
		"is_private":  c.IsPrivate,
		"dm_id":       c.DMID,
	}
	sql, args, err := PSQLBuilder().
		Update(channelsTable).
		SetMap(setMap).
		Where(sq.Eq{"id": id}).
		Suffix(channelsReturning).
		ToSql()
	if err != nil {
		return nil, err
	}

	updatedChannel := &channels.Channel{}
	err = backend.db.QueryRowx(sql, args...).StructScan(updatedChannel)
	if err != nil {
		if err == sqlP.ErrNoRows {
			return nil, channels.ErrChannelNotFound
		}
		return nil, err
	}

	return updatedChannel, nil
}
