package database

import (
	"database/sql"
	"log"
)

var DBConn *sql.DB

func SetupDBConnection() {
	var err error
	DBConn, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/inventorydb")
	if err != nil {
		log.Fatal(err)
	}

}
