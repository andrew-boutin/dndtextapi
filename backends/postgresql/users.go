package postgresql

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/andrew-boutin/dndtextapi/users"
	log "github.com/sirupsen/logrus"
)

const (
	usersTable = "users"
)

var userColumns = []string{
	"id",
	"name",
	"bio",
	"createdon",
	"lastupdated",
}

// GetUsersInChannel retrieves all Users that are in the Channel matching
// the provided ID.
func (backend Backend) GetUsersInChannel(id int) (users.UserCollection, error) {
	sql, args, err := PSQLBuilder().
		Select("id", "name", "bio", "users.createdon", "users.lastupdated").
		From(usersTable).
		Join(fmt.Sprintf("%s ON %s.id = %s.userid", channelsUsersTable, usersTable, channelsUsersTable)). // TODO: Is this an inner join?
		Where(sq.Eq{"channelid": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := backend.db.Queryx(sql, args...)

	if err != nil {
		return nil, err
	}

	usersCollection := make(users.UserCollection, 0)
	for rows.Next() {
		var user users.User
		err = rows.StructScan(&user)

		if err != nil {
			return nil, err
		}

		usersCollection = append(usersCollection, user)
	}

	return usersCollection, nil
}

// AddUsersToChannel adds the given Users matching the input User IDs to the Channel
// that matches the input Channel ID.
func (backend Backend) AddUsersToChannel(channelID int, userIDs []int) error {
	// Short circuit if there are no Users to add
	if len(userIDs) <= 0 {
		return nil
	}

	// TODO: unique constraint on channelid/userid
	builder := PSQLBuilder().
		Insert(channelsUsersTable).
		Columns("channelid", "userid")

	for _, userID := range userIDs {
		builder = builder.Values(channelID, userID)
	}

	sql, args, err := builder.ToSql()

	if err != nil {
		log.WithError(err).Error("Issue building add channel user mapping sql.")
		return err
	}

	// TODO: Check value
	_, err = backend.db.Exec(sql, args...)
	return err
}

// RemoveUsersFromChannel removes all Users from the Channel that matches
// the given ID.
func (backend Backend) RemoveUsersFromChannel(id int) error {
	sql, args, err := PSQLBuilder().
		Delete(channelsUsersTable).
		Where(sq.Eq{"channelid": id}).
		ToSql()

	if err != nil {
		return err
	}

	// TODO: Check result?
	_, err = backend.db.Exec(sql, args...)
	return err
}
