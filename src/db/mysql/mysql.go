package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func TestConnection() bool {
	// Формат строки подключения: "username:password@tcp(127.0.0.1:3306)/dbname"
	databaseURL := os.Getenv("DATABASE_URL_MYSQL")

	if databaseURL == "" {
		log.Fatalf("Database connection string missing")
		return false
	}

	// Открываем соединение с MySQL
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return false
	}
	defer db.Close()

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return false
	}

	// Пример запроса для проверки соединения
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
		return false
	}

	log.Println("Connected to MySQL version:", version)
	return true
}

func Connect() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL_MYSQL")

	if databaseURL == "" {
		return nil, fmt.Errorf("Database connection string missing")
	}

	// Открываем соединение с MySQL
	conn, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil, err
	}
	//defer conn.Close()

	return conn, nil
}

func Deconect(conn *sql.DB) {
	conn.Close()
}
