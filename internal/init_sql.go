package internal

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() {
	db, err := sql.Open(
		"mysql",
		"root:lian@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local",
	)
	if err != nil {
		fmt.Println("open database failed:", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("database connect failed:", err)
		return
	}

	fmt.Println("database connect success")

	sqlBytes, err := os.ReadFile("db/init.sql")
	if err != nil {
		fmt.Println("read sql file failed:", err)
		return
	}

	sqlContent := string(sqlBytes)

	_, err = db.Exec(sqlContent)
	if err != nil {
		fmt.Println("execute sql failed:", err)
		return
	}

	fmt.Println("execute sql success")
}
