package query

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Database struct {
	DB         *sql.DB
	sqlStament string
}

func NewDb() *Database {
	return &Database{nil, ""}
}

func (db *Database) connect_to_database(user, pw, dbName string) error {
	// Connect to database
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pw, dbName)

	var err error
	db.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connection to database %s successfully completed.", dbName)

	defer db.DB.Close()

	return nil
}

func (db *Database) send_data(sqlTable string, parameters []string, values []string) (id int) {

	// Process parameters for staments in sql
	myFormattedParameters := strings.Join(parameters, ", ")

	// Process values for staments in sql
	myFormattedValues := strings.Join(values, ", ")

	sqlStatement := fmt.Sprintf(`
	INSERT INTO %s (%s)
	VALUES (%s)
	RETURNING id`, sqlTable, myFormattedParameters, myFormattedValues)

	var err error
	err = db.DB.QueryRow(sqlStatement).Scan(&id)
	if err != nil {
		panic(err)
	}

	fmt.Println("New record ID is:", id)

	return id
}

func (db *Database) get_all_data(parameters []string, sqlTable string) (row_output *sql.Rows, status int) {

	myFormattedParameters := "*"
	if parameters != nil {
		myFormattedParameters = strings.Join(parameters, ", ")
	}
	sqlStatement := fmt.Sprintf("SELECT %s FROM %s", myFormattedParameters, sqlTable)

	// get data from table
	var err error
	row_output, err = db.DB.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer row_output.Close()

	fmt.Println("Query successfully finished")
	return row_output, status
}
