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

func (backend PostgresqlBackend) GetUsersInChannel(id int) (users.UserCollection, error) {
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

// TODO: unique constraint on channelid/userid
func (backend PostgresqlBackend) AddUsersToChannel(channelID int, userIDs []int) error {
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

func (backend PostgresqlBackend) RemoveUsersFromChannel(id int) error {
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
