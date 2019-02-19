package main

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func SaveTicketToDB(ticketType, email string) {
	address, username, password, database := GetSQLConfig()

	connectionString := concatStrings(
		username, ":", password, 
		"@(", address, ")/", database,
	)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Printf("Cannot connect to SQL server at address: '%s'", address)
        log.Fatal(err)
	}

	sqlInsert := "INSERT INTO SoldTickets (ticketType, email) VALUES(?, ?)"
	stmt, err := db.Prepare(sqlInsert)
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(ticketType, email)
	if err != nil {
		log.Fatal(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
}