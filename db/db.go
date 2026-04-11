package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDb() *sql.DB {
	db, err := sql.Open("sqlite", "file:app.db?_foreign_keys=on")
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}
