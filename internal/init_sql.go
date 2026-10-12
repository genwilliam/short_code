package internal

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func InitDB(config DatabaseConfig) (*sql.DB, error) {
	driverConfig, err := mysql.ParseDSN(config.DSN())
	if err != nil {
		return nil, err
	}

	schema, err := os.ReadFile(config.SchemaFile)
	if err != nil {
		return nil, err
	}

	bootstrapConfig := *driverConfig
	bootstrapConfig.DBName = ""
	bootstrapConfig.MultiStatements = true
	bootstrapDB, err := sql.Open("mysql", bootstrapConfig.FormatDSN())
	if err != nil {
		return nil, err
	}
	defer bootstrapDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := bootstrapDB.PingContext(ctx); err != nil {
		return nil, err
	}
	if _, err := bootstrapDB.ExecContext(ctx, string(schema)); err != nil {
		return nil, err
	}

	driverConfig.DBName = "short_url"
	driverConfig.MultiStatements = false
	db, err := sql.Open("mysql", driverConfig.FormatDSN())
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
