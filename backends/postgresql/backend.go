// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package postgresql

import (
	"fmt"
	"io/ioutil"
	"strings"

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
func MakePostgresqlBackend(user, password, dbname string) Backend {
	dbinfo := fmt.Sprintf("host=db user=%s password=%s dbname=%s sslmode=disable",
		user, password, dbname)
	db, err := sqlx.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}

	err = RunHealthCheck(db)
	if err != nil {
		panic(err)
	}

	rows, err := db.Queryx("SELECT * FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';")
	if err != nil {
		panic(err)
	}

	if !rows.Next() {
		initSchema(db)
	}

	return Backend{db: db}
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
		log.Fatalf("DB health check failed: %s", err)
	}
	return err
}

// initSchema initializes the schema for the Postgresql database.
func initSchema(db *sqlx.DB) {
	// The function defined in functionsFilePath has multiple semi colons that
	// are part of a single statement so it gets processed differently
	// than the statements in the schema file
	file, err := ioutil.ReadFile(functionsFilePath)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(string(file))
	if err != nil {
		panic(err)
	}

	file, err = ioutil.ReadFile(schemaFilePath)
	if err != nil {
		panic(err)
	}

	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		_, err = db.Exec(request)
		if err != nil {
			panic(err)
		}
	}
}
