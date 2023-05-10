package db

import (
	"database/sql"
	"os"
	"log"
	_ "rsc.io/sqlite"

)

var handle *sql.DB

func InitializeHandle() {
	var err error
	handle, err = sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	if err := handle.Ping(); err != nil {
		log.Fatal(err)
	}
}

func ApplySchema() {
	schema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = handle.Exec(string(schema))
	if err != nil {
		log.Fatal("apply schema: ", err)
	}
}