package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Adapter struct {
	DB *bun.DB
}

func NewAdapter(dataSourceName string) (*Adapter, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dataSourceName)))

	db := bun.NewDB(sqldb, pgdialect.New())

	// Optional: set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Ping the database to verify the connection
	err := db.Ping()
	if err != nil {
		log.Fatalf("DB ping failure %v", err)
	}

	return &Adapter{DB: db}, nil
}

func (da *Adapter) CloseDbConnection() {
	err := da.DB.Close()
	if err != nil {
		log.Fatalf("DB close failure: %v", err)
	}
}
