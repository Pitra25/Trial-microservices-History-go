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
		log.Fatalln("database connection string missing")
		return false
	}

	// Открываем соединение с MySQL
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatalln("failed to connect to the database: ", err)
		return false
	}
	defer db.Close()

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatalln("failed to ping database: ", err)
		return false
	}

	// Пример запроса для проверки соединения
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatalln("query failed: ", err)
		return false
	}

	log.Println("connected to MySQL version: ", version)
	return true
}

func Connect() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL_MYSQL")

	if databaseURL == "" {
		return nil, fmt.Errorf("database connection string missing")
	}

	// Открываем соединение с MySQL
	conn, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatalln("failed to connect to the database: ", err)
		return nil, err
	}
	//defer conn.Close()

	return conn, nil
}

func Deconect(conn *sql.DB) {
	conn.Close()
}
