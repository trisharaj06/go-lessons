package db

import (
	"database/sql"
	_"modernc.org/sqlite"
)

var DB *sql.DB
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "data.db")

	if err != nil{
		panic("could not connect to the database!")
	}

	DB.SetMaxOpenConns(5)
	createTable()
}

func createTable(){
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT,
		dateTime DATETIME,
		user_id INTEGER
	)
	`
	_, err := DB.Exec(createEventsTable)

	if err != nil{
		panic("unable to create Event Table")
	}
}
