package test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestConnDB(t *testing.T) {
	dsn := "root:lian@tcp(127.0.0.1:3306)/"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected!")

	rows, err := db.Query("SHOW DATABASES;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var databaseName string

		if err := rows.Scan(&databaseName); err != nil {
			log.Fatal(err)
		}

		fmt.Println(databaseName)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
