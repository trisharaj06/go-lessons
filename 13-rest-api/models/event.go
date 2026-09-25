package models

import (
	"rest-example/db"
	"time"
)

type Event struct {
	ID          int64
	Name        string `binding:"required"`
	Description string
	Location    string
	DateTime    time.Time
	UserId      int
}

var events = []Event{}

func (e Event) Save() error {
	query := `INSERT INTO events(name, description, location, dateTime, user_id)
	VALUES(?, ?, ?, ?, ?)
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserId)

	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	e.ID = id
	return err

}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserId)

		if err != nil {
			return nil, err
		}
		events = append(events, event)

	}
	return events, nil
}
