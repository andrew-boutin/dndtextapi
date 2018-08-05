// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package postgresql

import (
	sqlP "database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/andrew-boutin/dndtextapi/channels"
	log "github.com/sirupsen/logrus"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const (
	schemaFilePath    = "./backends/postgresql/schema.sql"
	functionsFilePath = "./backends/postgresql/functions.sql"
)

// Backend contains all of the data specific to a Postgres backend
type Backend struct {
	db *sqlx.DB
}

// MakePostgresqlBackend creates a Postgresql backend with connection to the
// actual DB, verifies the connection, and initializes the schema if it
// isn't already populated.
func MakePostgresqlBackend(user, password, dbname string) (b Backend, err error) {
	dbinfo := fmt.Sprintf("host=db user=%s password=%s dbname=%s sslmode=disable",
		user, password, dbname)
	db, err := sqlx.Open("postgres", dbinfo)
	if err != nil {
		log.WithError(err).Error("Failed to connect to database.")
		return b, err
	}

	err = RunHealthCheck(db)
	if err != nil {
		log.WithError(err).Error("Error during db health check.")
		return b, err
	}

	rows, err := db.Queryx("SELECT * FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';")
	if err != nil {
		log.WithError(err).Error("Failed to check if db was already populated.")
		return b, err
	}

	if !rows.Next() {
		err = initSchema(db)
		if err != nil {
			log.WithError(err).Error("Failed to initialize db schema.")
			return b, err
		}
	}

	return Backend{db: db}, nil
}

// PSQLBuilder retruns a squirrel SQL builder that uses
// placeholders in the format that Postgresql expects.
func PSQLBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

// RunHealthCheck returns an error if there is an issue
// connecting to the database.
func RunHealthCheck(db *sqlx.DB) error {
	err := db.Ping()
	if err != nil {
		log.WithError(err).Error("DB health check failed.")
	}
	return err
}

// initSchema initializes the schema for the Postgresql database.
func initSchema(db *sqlx.DB) error {
	// The function defined in functionsFilePath has multiple semi colons that
	// are part of a single statement so it gets processed differently
	// than the statements in the schema file
	file, err := ioutil.ReadFile(functionsFilePath)
	if err != nil {
		log.WithError(err).Error("Failed to read in sql file defining postresql functions.")
		return err
	}

	_, err = db.Exec(string(file))
	if err != nil {
		log.WithError(err).Error("Failed execute postresql functions.")
		return err
	}

	file, err = ioutil.ReadFile(schemaFilePath)
	if err != nil {
		log.WithError(err).Error("Failed to read in sql file defining postresql schema.")
		return err
	}

	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		_, err = db.Exec(request)
		if err != nil {
			log.WithError(err).Error("Failed to execute sql to init schema.")
			return err
		}
	}

	return nil
}

// getSingle retrieves a single row from the given table matching the given id. The columns
// retrieved are also determined from the input. The resulting record is loaded into the
// input object, not returned, so a pointer should be passed in for persistance. If no record
// exists then the return flag will be false.
func (backend Backend) getSingle(id int, tableName string, cols []string, obj interface{}) (wasFound bool, err error) {
	sql, args, err := PSQLBuilder().
		Select(cols...).
		From(tableName).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Issue building query for get single.")
		return false, err
	}

	// This could be an error from the thing not existing which is not unexpected
	err = backend.db.Get(obj, sql, args...)
	if err != nil {
		if err == sqlP.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}

// updateSingle updates a single row from the given table matching the given id. The input map
// determines what columns should be updated to what values. The returning parameter dictates
// what fields to retrieve after the update. The resulting record is loaded into the
// input object, not returned, so a pointer should be passed in for persistance. If no record
// exists then the return flag will be false.
func (backend Backend) updateSingle(id int, tableName, returning string, setMap map[string]interface{}, obj interface{}) (wasFound bool, err error) {
	sql, args, err := PSQLBuilder().
		Update(tableName).
		SetMap(setMap).
		Where(sq.Eq{"id": id}).
		Suffix(returning).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build query for update single.")
		return false, err
	}

	// This could be an error from the thing not existing which is not unexpected
	err = backend.db.QueryRowx(sql, args...).StructScan(obj)
	if err != nil {
		if err == sqlP.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}

// deleteSingle deletes a single row from the given table matching the given id. If no record
// exists then the return flag will be false.
func (backend Backend) deleteSingle(id int, tableName string) (wasFound bool, err error) {
	sql, args, err := PSQLBuilder().
		Delete(tableName).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build query for delete single.")
		return false, err
	}

	result, err := backend.db.Exec(sql, args...)
	if err != nil {
		log.WithError(err).Error("Failed to execute delete single query.")
	}

	numRowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithError(err).Error("Failed to determine how many rows were affected by delete query.")
		return false, err
	}

	if numRowsAffected <= 0 {
		return false, channels.ErrChannelNotFound
	}

	return true, err
}

// createSingle creates a single row in the given table. The input map determines what columns
// should be set to what values. The returning parameter dictates what fields to retrieve after
// the create. The resulting record is loaded into the input object, not returned, so a pointer
// should be passed in for persistance.
func (backend Backend) createSingle(tableName, returning string, kvs map[string]interface{}, obj interface{}) (err error) {
	cols := make([]string, len(kvs))
	vals := make([]interface{}, len(kvs))
	i := 0
	for k, v := range kvs {
		cols[i] = k
		vals[i] = v
		i++
	}

	sql, args, err := PSQLBuilder().
		Insert(tableName).
		Columns(cols...).
		Values(vals...).
		Suffix(returning).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Issue building create single sql.")
		return err
	}

	err = backend.db.QueryRowx(sql, args...).StructScan(obj)
	if err != nil {
		log.WithError(err).Error("Issue running create single sql.")
		return err
	}

	return
}
