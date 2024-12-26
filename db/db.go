package db

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

// InitDB initializes the database connection and creates tables if necessary
func InitDB() {
	var err error
	// PostgreSQL connection DSN
	dsn := "postgresql://postgres:mlgusovnMIucjAvgKfVmTLEZrADmzOHR@autorack.proxy.rlwy.net:11015/railway"

	// Open a connection to the PostgreSQL database
	DB, err = sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Error opening PostgreSQL database: %v", err)
		return
	}

	_, err = DB.Exec("SET CLIENT_ENCODING TO 'UTF8';")
	if err != nil {
		log.Fatal("Error setting encoding:", err)
	}

	// Verify the connection by pinging the database
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging PostgreSQL database: %v", err)
		return
	}
	log.Println("Successfully connected to the PostgreSQL database")
}
