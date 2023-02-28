package main

import (
	"query"
)

func main() {

	my_db := query.Database{}

	my_db.ConnectToDatabase("postgres", "mysecretpassword", "temp")
	sqlTable = "temperature"
	temp = 25.25
	t := time.Now()
	formattedTime := t.Format("2006-01-02 15:04:05")
	parameters := ["temperature ", "date"]
	values := [temp, t]
	my_db.send_data(sqlTable, parameters, values)
}
