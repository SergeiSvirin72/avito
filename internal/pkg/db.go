package pkg

import (
	"database/sql"
	"log"
)

func NewDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	return db
}
