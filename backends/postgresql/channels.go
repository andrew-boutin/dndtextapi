package postgresql

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/andrew-boutin/dndtextapi/channels"
	log "github.com/sirupsen/logrus"
)

const (
	channelsTable      = "channels"
	channelsUsersTable = "channels_users"
)

var channelColumns = []string{
	"name",
	"description",
	"id",
	"ownerid",
	"createdon",
	"lastupdated",
	"isprivate",
	"dmid",
}

func (backend PostgresqlBackend) GetChannel(id int) (*channels.Channel, error) {
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
		return nil, err
	}

	// TODO: Maybe partial users instead of full users
	users, err := backend.GetUsersInChannel(id)

	if err != nil {
		return nil, err
	}

	channel.Users = users
	return channel, err
}

// TODO: Be able to filter on things such as private/not private
func (backend PostgresqlBackend) GetChannels() (*channels.ChannelCollection, error) {
	sql, args, err := PSQLBuilder().
		Select(channelColumns...).
		From(channelsTable).
		ToSql()

	if err != nil {
		return nil, err
	}

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

		outChannels = append(outChannels, channel)
	}

	return &outChannels, nil
}

// TODO: Don't require description, default isprivate to false
func (backend PostgresqlBackend) CreateChannel(c *channels.Channel) (*channels.Channel, error) {
	sql, args, err := PSQLBuilder().
		Insert(channelsTable).
		Columns("name", "description", "ownerid", "isprivate", "dmid").
		Values(c.Name, c.Description, c.OwnerID, c.IsPrivate, c.DMID).
		Suffix("RETURNING id, name, description, ownerid, isprivate, dmid, createdon, lastupdated"). // TODO: Use channelColumns...
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

	// TODO: Need to add Users before returning..?
	err = backend.AddUsersToChannel(newChannel.ID, c.Users.GetUserIDs())

	if err != nil {
		log.WithError(err).Error("Issue adding user channel mappings.")
		return nil, err
	}

	return newChannel, nil
}

func (backend PostgresqlBackend) DeleteChannel(id int) error {
	err := backend.RemoveUsersFromChannel(id)

	if err != nil {
		return err
	}

	sql, args, err := PSQLBuilder().
		Delete(channelsTable).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return err
	}

	// TODO: Check result?
	_, err = backend.db.Exec(sql, args...)

	return err
}

func (backend PostgresqlBackend) UpdateChannel(id int, c *channels.Channel) (*channels.Channel, error) {
	setMap := map[string]interface{}{
		"name":        c.Name,
		"description": c.Description,
		"ownerid":     c.OwnerID,
		"isprivate":   c.IsPrivate,
		"dmid":        c.DMID,
	}
	sql, args, err := PSQLBuilder().
		Update(channelsTable).
		SetMap(setMap).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, name, description, ownerid, isprivate, dmid, createdon, lastupdated"). // TODO: Use channelColumns...
		ToSql()

	if err != nil {
		return nil, err
	}

	updatedChannel := &channels.Channel{}
	err = backend.db.QueryRowx(sql, args...).StructScan(updatedChannel)
	if err != nil {
		return nil, err
	}

	err = backend.RemoveUsersFromChannel(id)

	if err != nil {
		return nil, err
	}

	// TODO: Need to add Users before returning..?
	err = backend.AddUsersToChannel(id, c.Users.GetUserIDs())

	if err != nil {
		return nil, err
	}

	return updatedChannel, nil
}
