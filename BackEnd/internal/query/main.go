package main

import (
	"encoding/json"
	"fmt"

	query "github.com/ltovarm/Geoffrey_App/BackEnd/internal/query/queries"
)

func main() {

	my_db := query.NewDb()
	if err := my_db.ConnectToDatabaseFromEnvVar(); err != nil {
		fmt.Printf("Error connecting to db: %v\n", err)
		return
	}
	sqlTable := "temperatures"

	data, err := my_db.GetLastData(sqlTable)
	if err != nil {
		fmt.Printf("Error Get LastData to db: %v\n", err)
		return
	}

	my_db.CloseDatabase()

	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error al serializar objeto a JSON: ", err)
		return
	}
	fmt.Println(string(jsonString))

}
