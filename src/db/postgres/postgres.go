package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func TestConnection() bool {

	database_url := os.Getenv("DATABASE_URL_POSTGRES")

	if database_url == "" {
		log.Fatalf("Database connection string missing")
		return false
	}

	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return false
	}
	defer conn.Close(context.Background())

	// Example query to test connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
		return false
	}

	log.Println("Connected to: ", version)
	return true
}

func Connect() (*pgx.Conn, error) {
	database_url := os.Getenv("DATABASE_URL_POSTGRES")
	if database_url == "" {
		return nil, fmt.Errorf("Database connection string missing")
	}

	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	return conn, nil
}

func Deconect(conn *pgx.Conn) {
	conn.Close(context.Background())
}
