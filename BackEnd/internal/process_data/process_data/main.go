package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to data base
	db, err := sql.Open("postgres", "postgres://admin:admin@localhost/draft?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// get data from table
	rows, err := db.Query("SELECT id, temp, date FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// print data
	for rows.Next() {
		var id int
		var temp float32
		var date time.Time
		err := rows.Scan(&id, &temp, &date)
		if err != nil {
			panic(err)
		}
		fmt.Printf("id: %d, name: %f, date: %v\n", id, temp, date.Format("02/01/2006 15:04:05"))
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}
