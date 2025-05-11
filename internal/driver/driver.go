package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB holds the database connection pool
// in case I ever want to change the database
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL creates database pool from postgres
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(maxDbLifetime)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetMaxOpenConns(maxOpenDbConn)

	dbConn.SQL = db

	err = TestDB(db)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// TestDB pings the database
func TestDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase creates a new database
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
