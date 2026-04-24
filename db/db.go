package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func InitDb() *sql.DB {
	dsn := os.Getenv("DB_PATH")
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}
