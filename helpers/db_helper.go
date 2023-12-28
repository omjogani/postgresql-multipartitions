package helper

import (
	"database/sql"
	"log"
)

func ConnectToDb(connectionString string) *sql.DB {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("ERROR IN CONNECTION TO POSTGRE: ", err)
	}
	return db
}
