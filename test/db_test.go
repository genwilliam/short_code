package test

import (
	"database/sql"
	"fmt"
	"testing"
  _ "github.com/go-sql-driver/mysql"
)

type dbTest struct {
	db *sql.DB
}

func DB(db *sql.DB) *dbTest {
	return &dbTest{db: db}
}

func TestRun(t *testing.T) {
	dbs, err := sql.Open("mysql", "root:lian@tcp(127.0.0.1:3306)/")
	if err != nil {
		t.Fatal(err)
	}
	defer dbs.Close()

	if err := dbs.Ping(); err != nil {
		t.Fatal(err)
	}

	fmt.Println("Database connected!")
	DB(dbs)
}
