package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	//USERNAME in database
	USERNAME = "leonardo"
	//DBNAME name of the data base
	DBNAME = "challenge"
	//PORT of data base
	PORT = "26257"
)

func createDBConnection() (*sql.DB, error) {

	var strconnection string = "user=" + USERNAME + " dbname=" + DBNAME + " sslmode=disable port=" + PORT
	return sql.Open("postgres", strconnection)
}
