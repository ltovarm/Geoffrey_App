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
	Err        error
}

func (bd *Database) connect_to_database(user, pw, dbName string) {
	// Connect to database
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pw, dbName)

	bd.DB, bd.Err = sql.Open("postgres", connStr)
	if bd.Err != nil {
		panic(bd.Err)
	}

	fmt.Printf("Connection to database %s successfully completed.", dbName)

	defer bd.DB.Close()
}

func (db Database) send_data(sqlTable string, parameters []string, values []string) (id int) {

	// Process parameters for staments in sql
	myFormattedParameters := strings.Join(parameters, ", ")

	// Process values for staments in sql
	myFormattedValues := strings.Join(values, ", ")

	sqlStatement := fmt.Sprintf(`
	INSERT INTO %s (%s)
	VALUES (%s)
	RETURNING id`, sqlTable, myFormattedParameters, myFormattedValues)

	db.Err = db.DB.QueryRow(sqlStatement).Scan(&id)
	if db.Err != nil {
		panic(db.Err)
	}

	fmt.Println("New record ID is:", id)

	return id
}

func (db Database) get_all_data(parameters []string, sqlTable string) (row_output, status) {

	myFormattedParameters := "*"
	if parameters != nil {
		myFormattedParameters = strings.Join(parameters, ", ")
	}
	sqlStatement := fmt.Sprintf("SELECT %s FROM %s", myFormattedParameters, sqlTable)

	// get data from table
	row_output, db.Err = db.DB.Query(sqlStatement)
	if db.Err != nil {
		panic(db.Err)
	}
	defer row_output.Close()

	fmt.Println("Query successfully finished")
}
