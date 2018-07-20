// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

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
	"isadmin",
	"isbanned",
	"createdon",
	"lastupdated",
}

func init() {
	// Add the User table name in front of the columms to avoid ambigious references.
	for i, col := range userColumns {
		userColumns[i] = fmt.Sprintf("%s.%s", usersTable, col)
	}
}

// GetUsersInChannel retrieves all Users that are in the Channel matching
// the provided ID.
func (backend Backend) GetUsersInChannel(id int) (*users.UserCollection, error) {
	sql, args, err := PSQLBuilder().
		Select("id", "username", "email", "bio", "isadmin", "isbanned", "users.createdon", "users.lastupdated").
		From(usersTable).
		Join(fmt.Sprintf("%s ON %s.id = %s.userid", channelsUsersTable, usersTable, channelsUsersTable)). // TODO: Is this an inner join?
		Where(sq.Eq{"channelid": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	return backend.runMultiUsersQuery(sql, args)
}

// GetAllUsers retrieves all Users from the database - including their User.IsAdmin flag.
func (backend Backend) GetAllUsers() (*users.UserCollection, error) {
	selects := []string{"id", "username", "email", "bio", "isadmin", "isbanned", "users.createdon", "users.lastupdated"}
	sql, args, err := PSQLBuilder().
		Select(selects...).
		From(usersTable).
		ToSql()
	if err != nil {
		return nil, err
	}

	return backend.runMultiUsersQuery(sql, args)
}

// runMultiUsersQuery runs the given query with arguments to retrieve a User
// collection.
func (backend Backend) runMultiUsersQuery(sql string, args []interface{}) (*users.UserCollection, error) {
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

	return &usersCollection, nil
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

// UpdateUser updates the given User with the given User data.
func (backend Backend) UpdateUser(id int, u *users.User) (*users.User, error) {
	setMap := map[string]interface{}{
		"username": u.Username,
		"bio":      u.Bio,
	}
	sql, args, err := PSQLBuilder().
		Update(usersTable).
		SetMap(setMap).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, username, email, bio, isadmin, isbanned, createdon, lastupdated"). // TODO: Use messageColumns...
		ToSql()
	if err != nil {
		return nil, err
	}

	updatedUser := &users.User{}
	err = backend.db.QueryRowx(sql, args...).StructScan(updatedUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// DeleteUser removes a User from the Users table.
func (backend Backend) DeleteUser(userID int) error {
	sql, args, err := PSQLBuilder().
		Delete(usersTable).
		Where(sq.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	// TODO: Check value
	_, err = backend.db.Exec(sql, args...)
	return err
}

// GetUserByEmail retrieves a User by using the given email address.
// TODO: Combine with GetUserByID and take in query params of some sort
func (backend Backend) GetUserByEmail(email string) (*users.User, error) {
	sql, args, err := PSQLBuilder().
		Select(userColumns...).
		From(usersTable).
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, err
	}

	return backend.runSingleUserQuery(sql, args)
}

// GetUserByID retrieves a User by using the given id.
func (backend Backend) GetUserByID(id int) (*users.User, error) {
	sql, args, err := PSQLBuilder().
		Select(userColumns...).
		From(usersTable).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	return backend.runSingleUserQuery(sql, args)
}

// runSingleUserQuery executes the given query with the given arguments to
// retrieve a single User.
func (backend Backend) runSingleUserQuery(sql string, args []interface{}) (*users.User, error) {
	user := &users.User{}
	err := backend.db.Get(user, sql, args...)
	if err != nil {
		if err == sqlP.ErrNoRows {
			return nil, users.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// CreateUser creates a new User in the database using the provided data.
func (backend Backend) CreateUser(gu *users.GoogleUser) (*users.User, error) {
	sql, args, err := PSQLBuilder().
		Insert(usersTable).
		Columns("username", "email").
		Values(gu.Email, gu.Email).
		Suffix("RETURNING id, username, email, bio, isadmin, isbanned, createdon, lastupdated"). // TODO: Use userColumns...
		ToSql()
	if err != nil {
		log.WithError(err).Error("Issue building create user sql.")
		return nil, err
	}

	newUser := &users.User{}
	err = backend.db.QueryRowx(sql, args...).StructScan(newUser)
	if err != nil {
		log.WithError(err).Error("Issue running create user sql.")
		return nil, err
	}

	return newUser, nil
}
