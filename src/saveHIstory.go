package src

import (
	postgres "Trial-microservices-History-go/src/db"
	"Trial-microservices-History-go/src/storage"
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"

	_ "github.com/go-sql-driver/mysql"
)

// type createrPromise struct {
// 	Calculation string `json:"calculation"`
// 	CreatedAt   string `json:"created_at"`
// }

// func SaveHistory(bode createrPromise) error {
func SaveHistory(Calculation string, CreatedAt string) error {

	if !postgres.TestConection() {
		panic("Error conection!!!")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return err
	}
	defer conn.Close(context.Background())

	result, err := conn.Exec(context.Background(), "insert into Histores (created_at, Calculation) values (?, ?)",
		CreatedAt, Calculation)

	if err != nil {
		log.Fatalf("Failed to create history entry: %v", err)
		return err
	}

	//storage.CreatrRecording()

	log.Print(result.RowsAffected())
	return nil
}

type Recording struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	Calculation string `json:"calculation"`
}

func GetHistory(key string) []Recording {

	storageRecords, err := storage.GetRecordFromHash(key)
	if err != nil {
		log.Fatalf("Error get redis recording: %v", err)
	}

	result := make([]Recording, len(storageRecords))
	for i, record := range storageRecords {
		result[i] = Recording{
			ID:          record.ID,
			Calculation: record.Calculation,
			CreatedAt:   record.CreatedAt,
		}
	}

	if len(result) == 0 {
		log.Fatalf("Error search: no records found")
	}
	defer log.Fatalf("Error search ")

	// PostgreSql

	if !postgres.TestConection() {
		panic("Error conection!!!")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select * from Histores")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	var historesArray []Recording

	for rows.Next() {
		var histore Recording

		err := rows.Scan(&histore.ID, &histore.Calculation, &histore.CreatedAt)
		if err != nil {
			log.Fatalf("Failed to get entry history: %v", err)
			continue
		}

		historesArray = append(historesArray, histore)
	}

	return historesArray
}
