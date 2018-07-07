package postgresql

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/andrew-boutin/dndtextapi/channels"
)

const channelsTable = "channels"

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
		return nil, err
	}

	newChannel := &channels.Channel{}
	err = backend.db.QueryRowx(sql, args...).StructScan(newChannel)
	if err != nil {
		return nil, err
	}
	return newChannel, nil
}

// TODO: Delete, Update
