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

	defer db.DB.Close()

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

	defer db.DB.Close()

	return err
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
	_, err = db.DB.Exec(query, jsonData)
	if err != nil {
		log.Fatalf("Error inserting into table: %s", err)
	}

	return nil
}

func (db *Database) GetXData(sqlTable string, numberOfData int) (map[string]interface{}, error) {
	// Query to get the last row from the table
	row := db.DB.QueryRow("SELECT data FROM "+sqlTable+" ORDER BY date DESC LIMIT %d", numberOfData)

	// Variables to store the JSONB data and to decode it into a map[string]interface{}
	var jsonData []byte
	data := make(map[string]interface{})

	// Scan the data into the jsonData variable
	err := row.Scan(&jsonData)
	if err != nil {
		return nil, err
	}

	// Decode the JSONB data into the data map
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetLastData(sqlTable string, numberOfData int) (map[string]interface{}, error) {

	data, err := db.GetXData(sqlTable, 1)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (db *Database) GetAllData(parameters []string, sqlTable string) (row_output *sql.Rows, status int) {

	myFormattedParameters := "*"
	if parameters != nil {
		myFormattedParameters = strings.Join(parameters, ", ")
	}

	sqlStatement := fmt.Sprintf("SELECT COUNT(id) FROM %s;", sqlTable)
	sqlStatement = fmt.Sprintf("SELECT %s FROM %s", myFormattedParameters, sqlTable)

	// get data from table
	var err error
	row_output, err = db.DB.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer row_output.Close()

	log.Println("Query successfully finished")
	return row_output, status
}
