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
	usersTable     = "users"
	usersReturning = "RETURNING id, username, email, bio, is_admin, is_banned, last_login, created_on, last_updated"
)

var userColumns = []string{
	"id",
	"username",
	"email",
	"bio",
	"is_admin",
	"is_banned",
	"last_login",
	"created_on",
	"last_updated",
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
		log.WithError(err).Error("Failed to create get all users sql.")
		return nil, err
	}

	return backend.runMultiUsersQuery(sql, args)
}

// runMultiUsersQuery runs the given query with arguments to retrieve a User
// collection.
func (backend Backend) runMultiUsersQuery(sql string, args []interface{}) (users.UserCollection, error) {
	rows, err := backend.db.Queryx(sql, args...)
	if err != nil {
		log.WithError(err).Error("Failed to execute multi user query.")
		return nil, err
	}

	usersCollection := make(users.UserCollection, 0)
	for rows.Next() {
		var user users.User
		err = rows.StructScan(&user)
		if err != nil {
			log.WithError(err).Error("Failed to load user result from multi user query.")
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

	updatedUser := &users.User{}
	err := backend.updateSingle(id, usersTable, usersReturning, setMap, updatedUser)
	if err != nil {
		if err == sqlP.ErrNoRows {
			return nil, users.ErrUserNotFound
		}
		log.WithError(err).Error("Issue with query for update user.")
		return nil, err
	}

	return updatedUser, nil
}

// DeleteUser removes a User from the Users table.
func (backend Backend) DeleteUser(userID int) error {
	wasDeleted, err := backend.deleteSingle(userID, usersTable)
	if err != nil {
		log.WithError(err).Error("Failed to execute delete user query.")
	} else if !wasDeleted {
		return users.ErrUserNotFound
	}
	return err
}

// GetUserByID retrieves a User by using the given id.
func (backend Backend) GetUserByID(id int) (*users.User, error) {
	user := &users.User{}
	err := backend.getSingle(id, usersTable, userColumns, user)
	if err != nil {
		if err == sqlP.ErrNoRows {
			return nil, users.ErrUserNotFound
		}
		log.WithError(err).Error("Query issue for get user.")
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a User by using the given email address.
func (backend Backend) GetUserByEmail(email string) (*users.User, error) {
	sql, args, err := PSQLBuilder().
		Select(userColumns...).
		From(usersTable).
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build get user by email query.")
		return nil, err
	}

	user := &users.User{}
	err = backend.db.Get(user, sql, args...)
	if err != nil {
		if err == sqlP.ErrNoRows {
			return nil, users.ErrUserNotFound
		}
		log.WithError(err).Error("Issue executing get user by email query.")
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
		Suffix(usersReturning).
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
		"last_login": time.Now(),
	}
	sql, args, err := PSQLBuilder().
		Update(usersTable).
		SetMap(setMap).
		Where(sq.Eq{"id": u.ID}).
		Suffix(usersReturning).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build update user last login query.")
		return nil, err
	}

	updatedUser := &users.User{}
	err = backend.db.QueryRowx(sql, args...).StructScan(updatedUser)
	if err != nil {
		log.WithError(err).Error("Failed to execute update user last login query.")
		return nil, err
	}

	return updatedUser, nil
}
