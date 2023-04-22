package main

import (
	"fmt"
	"time"

	query "github.com/ltovarm/Geoffrey_App/BackEnd/internal/query/queries"
)

func main() {

	my_db := query.NewDb()
	if err := my_db.Connect_to_database("postgres", "mysecretpassword", "temp"); err != nil {
		fmt.Printf("Error connecting to db: %v\n", err)
		return
	}
	sqlTable := "temperature"
	temp := 25.25
	t := time.Now()
	formattedTime := t.Format("2006-01-02 15:04:05")
	var parameters [2]string
	parameters = [2]string{"temperature", "date"}

	values := [2]string{fmt.Sprintf("%f", temp), formattedTime}
	my_db.Send_data(sqlTable, parameters[:], values[:])
}
