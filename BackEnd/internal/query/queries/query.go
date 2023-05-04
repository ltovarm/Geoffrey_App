package query

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type JSONResponse struct {
	Data []map[string]interface{} `json:"data"`
}

type Database struct {
	DB         *sql.DB
	sqlStament *string
}

func NewDb() *Database {
	return &Database{nil, nil}
}

func (db *Database) ConnectToDatabaseFromEnvVar() error {

	sqlServerURL := os.Getenv("DATABASE_URL")
	// Connect to database
	var err error
	db.DB, err = sql.Open("postgres", sqlServerURL)
	if err != nil {
		panic(err)
	}

	log.Printf("Connection to database successfully.")

	return err
}

func (db *Database) ConnectToDatabase(user, pw, dbName, dbContainerName, port string) error {
	// Connect to database
	sqlServerURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pw, dbContainerName, port, dbName)
	var err error
	db.DB, err = sql.Open("postgres", sqlServerURL)
	if err != nil {
		panic(err)
	}

	log.Printf("Connection to database %s successfully completed.", dbName)

	return err
}

func (db *Database) CloseDatabase() {
	db.DB.Close()
}

func (db *Database) SendData(sqlTable string, parameters []string, values []string) (id int) {

	// Process parameters for staments in sql
	myFormattedParameters := strings.Join(parameters, ", ")

	// Process values for staments in sql
	myFormattedValues := strings.Join(values, ", ")

	sqlStatement := fmt.Sprintf(`
	INSERT INTO %s (%s)
	VALUES (%s)
	RETURNING id`, sqlTable, myFormattedParameters, myFormattedValues)

	err := db.DB.QueryRow(sqlStatement).Scan(&id)
	if err != nil {
		panic(err)
	}

	log.Println("New record ID is:", id)

	return id
}

func (db *Database) SendDataAsJSON(data map[string]interface{}, sqlTable string) error {

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error when serializing the Json: %s", err)
	}

	// Insert into table
	query := fmt.Sprintf("INSERT INTO %s (data) VALUES ($1)", sqlTable)
	log.Println(query)
	_, err = db.DB.Exec(query, jsonData)
	if err != nil {
		log.Fatalf("Error inserting into table: %s", err)
		return err
	}

	return nil
}

func (db *Database) GetXData(sqlTable string, numberOfData int) ([]map[string]interface{}, error) {

	// Query to get all rows from the table
	rows, err := db.DB.Query(fmt.Sprintf("SELECT data FROM %s ORDER BY id DESC LIMIT %d", sqlTable, numberOfData))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Array to store the result
	result := []map[string]interface{}{}

	// Loop through each row and scan its data into a map
	for rows.Next() {
		var jsonData []byte
		data := make(map[string]interface{})

		err := rows.Scan(&jsonData)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	log.Println("Query successfully finished")

	return result, nil
}

func (db *Database) GetLastData(sqlTable string) (map[string]interface{}, error) {

	data, err := db.GetXData(sqlTable, 1)
	if err != nil {
		return nil, err
	}

	return data[0], err
}

func (db *Database) GetAllData(parameters []string, sqlTable string) ([]map[string]interface{}, error) {

	// Query to get all rows from the table
	rows, err := db.DB.Query(fmt.Sprintf("SELECT COUNT(id) FROM %s;", sqlTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Array to store the result
	result := []map[string]interface{}{}

	// Loop through each row and scan its data into a map
	for rows.Next() {
		var jsonData []byte
		data := make(map[string]interface{})

		err := rows.Scan(&jsonData)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	log.Println("Query successfully finished")

	return result, nil
}
