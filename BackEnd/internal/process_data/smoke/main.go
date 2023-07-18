package main

import (
	"database/sql"
	"fmt"
	"time"
)

func main() {
	// Conexión a la base de datos
	db, err := sql.Open("postgres", "postgres://admin:admin@localhost/draft?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Insertar datos en una tabla
	sqlStatement := `
        INSERT INTO temp1 (temp, time)
        VALUES ($1, $2)
        RETURNING id`
	id := 0
	temp := 25.15
	t := time.Now()
	formattedTime := t.Format("02/01/2006 15:04:05")
	err = db.QueryRow(sqlStatement, temp, formattedTime).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)
}
