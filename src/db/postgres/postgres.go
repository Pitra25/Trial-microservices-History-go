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
		log.Fatalln("database connection string missing")
		return false
	}

	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalln("failed to connect to the database: ", err)
		return false
	}
	defer conn.Close(context.Background())

	// Example query to test connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalln("query failed: ", err)
		return false
	}

	log.Println("Connected to: ", version)
	return true
}

func Connect() (*pgx.Conn, error) {
	database_url := os.Getenv("DATABASE_URL_POSTGRES")
	if database_url == "" {
		return nil, fmt.Errorf("database connection string missing")
	}

	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalln("failed to connect to the database: ", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	return conn, nil
}

func Deconect(conn *pgx.Conn) {
	conn.Close(context.Background())
}
