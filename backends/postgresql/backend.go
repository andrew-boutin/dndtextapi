package postgresql

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const schemaFilePath = "./backends/postgresql/schema.sql"

type PostgresqlBackend struct {
	db *sqlx.DB
}

// TODO: Should use transactions
func GetPostgresqlBackend(user, password, dbname string) PostgresqlBackend {
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

	return PostgresqlBackend{db: db}
}

func PSQLBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func RunHealthCheck(db *sqlx.DB) error {
	err := db.Ping()
	if err != nil {
		log.Fatalf("DB health check failed: %s", err)
	}
	return err
}

func initSchema(db *sqlx.DB) {
	file, err := ioutil.ReadFile(schemaFilePath)

	if err != nil {
		panic(err)
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		_, err := db.Exec(request)
		if err != nil {
			panic(err)
		}
	}
}
