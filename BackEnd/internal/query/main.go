package main

import (
	"time"

	"github.com/ltovarm/Geoffrey_App/tree/GA-5-Query/BackEnd/internal/query/query"
)

func main() {

	my_db := query.NewDb()
	my_db.connect_to_database("postgres", "mysecretpassword", "temp")
	sqlTable := "temperature"
	temp := 25.25
	t := time.Now()
	formattedTime := t.Format("2006-01-02 15:04:05")
	var parameters [2]string
	parameters = [2]string{"temperature", "date"}
	var values [2]interface{}
	values[0] = temp
	values[1] = formattedTime
	my_db.send_data(sqlTable, parameters, values)
}
