// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package postgresql

import (
	"fmt"
	"time"

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
	"lastlogin",
	"createdon",
	"lastupdated",
}

func init() {
	// Add the User table name in front of the columms to avoid ambigious references.
	for i, col := range userColumns {
		userColumns[i] = fmt.Sprintf("%s.%s", usersTable, col)
	}
}

// GetAllUsers retrieves all Users from the database - including their User.IsAdmin flag.
func (backend Backend) GetAllUsers() (users.UserCollection, error) {
	sql, args, err := PSQLBuilder().
		Select(userColumns...).
		From(usersTable).
		ToSql()
	if err != nil {
		return nil, err
	}

	return backend.runMultiUsersQuery(sql, args)
}

// runMultiUsersQuery runs the given query with arguments to retrieve a User
// collection.
func (backend Backend) runMultiUsersQuery(sql string, args []interface{}) (users.UserCollection, error) {
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

		usersCollection = append(usersCollection, &user)
	}

	return usersCollection, nil
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
		Suffix("RETURNING id, username, email, bio, isadmin, isbanned, lastlogin, createdon, lastupdated"). // TODO: Use messageColumns...
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
		Suffix("RETURNING id, username, email, bio, isadmin, isbanned, lastlogin, createdon, lastupdated"). // TODO: Use userColumns...
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

// UpdateUserLastLogin sets the passed in User's last login time to now.
func (backend Backend) UpdateUserLastLogin(u *users.User) (*users.User, error) {
	setMap := map[string]interface{}{
		"lastlogin": time.Now(),
	}
	sql, args, err := PSQLBuilder().
		Update(usersTable).
		SetMap(setMap).
		Where(sq.Eq{"id": u.ID}).
		Suffix("RETURNING id, username, email, bio, isadmin, isbanned, lastlogin, createdon, lastupdated"). // TODO: Use messageColumns...
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
