package postgresql

import (
	"fmt"

	sqlP "database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/andrew-boutin/dndtextapi/users"
	log "github.com/sirupsen/logrus"
)

const (
	usersTable = "users"
)

var userColumns = []string{
	"id",
	"username",
	"email",
	"bio",
	"createdon",
	"lastupdated",
}

// GetUsersInChannel retrieves all Users that are in the Channel matching
// the provided ID.
func (backend Backend) GetUsersInChannel(id int) (users.UserCollection, error) {
	sql, args, err := PSQLBuilder().
		Select("id", "username", "email", "bio", "users.createdon", "users.lastupdated").
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

// RemoveAllUsersFromChannel removes all Users from the Channel that matches
// the given ID.
func (backend Backend) RemoveAllUsersFromChannel(id int) error {
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

// IsUserInChannel determines if the given User is a member of the given Channel.
func (backend Backend) IsUserInChannel(userID, channelID int) (bool, error) {
	sql, args, err := PSQLBuilder().
		Select("1").
		From(channelsUsersTable).
		Where(sq.Eq{"channelid": channelID}).
		Where(sq.Eq{"userid": userID}).
		ToSql()
	if err != nil {
		return false, err
	}

	// TODO: Move out into own function
	// Inspiration https://snippets.aktagon.com/snippets/756-checking-if-a-row-exists-in-go-database-sql-and-sqlx-
	sql = fmt.Sprintf("select exists(%s)", sql)

	var exists bool
	err = backend.db.QueryRow(sql, args...).Scan(&exists)
	if err != nil && err != sqlP.ErrNoRows {
		return false, err
	}

	return exists, nil
}

// AddUserToChannel adds the given User to the given Channel
func (backend Backend) AddUserToChannel(userID, channelID int) error {
	sql, args, err := PSQLBuilder().
		Insert(channelsUsersTable).
		Columns("userid", "channelid").
		Values(userID, channelID).
		ToSql()
	if err != nil {
		return err
	}

	// TODO: Check value
	_, err = backend.db.Exec(sql, args...)
	return err
}
