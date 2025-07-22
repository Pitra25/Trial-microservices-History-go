package src

import (
	"Trial-microservices-History-go/src/db/mysql"
	"Trial-microservices-History-go/src/types"

	// "Trial-microservices-History-go/src/db/postgres"
	"Trial-microservices-History-go/src/storage"
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// type createrPromise struct {
// 	Calculation string `json:"calculation"`
// 	CreatedAt   string `json:"created_at"`
// }

// func SaveHistory(bode createrPromise) error {
func SaveHistory(Calculation string, CreatedAt string) error {

	// if !postgres.TestConnection() {
	// 	panic("Error conection!!!")
	// }
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL_MYSQL"))
	// if err != nil {
	// 	log.Fatalf("Failed to connect to the database: %v", err)
	// 	return err
	// }
	// defer conn.Close(context.Background())
	// jsonDate, err := json.Marshal(CreatedAt)
	// if err != nil {
	// 	log.Fatalf("Error while coding in JSON: %v", err)
	// 	return err
	// }
	// result, err := conn.Exec(context.Background(), "insert into historys (Calculation, Created_at) values ($1, $2)",
	// 	Calculation, CreatedAt)

	// MySql
	if !mysql.TestConnection() {
		panic("Error conection!!!")
	}

	conn, err := mysql.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return err
	}

	result, err := conn.Exec("insert into historys (Calculation, CreatedAt) values (?, ?)",
		Calculation, CreatedAt)
	if err != nil {
		log.Fatalf("Failed to create history entry: %v", err.Error())
		return err
	}
	defer mysql.Deconect(conn)

	log.Print(result)
	return nil
}

func GetHistory(key string) []*types.Recording {

	storageRecords, err := storage.GetRecordFromHash(key)
	if err != nil && len(storageRecords) != 0 {
		// log.Fatalf("Error get redis recording: %v", err)
		result := make([]*types.Recording, len(storageRecords))
		for i, record := range storageRecords {
			result[i] = &types.Recording{
				ID:          record.ID,
				Calculation: record.Calculation,
				CreatedAt:   record.CreatedAt,
			}
		}

		if len(result) == 0 {
			log.Fatalf("Error search: no records found")
		}
		defer log.Fatalf("Error search ")
	}

	// TEST save redis
	// save := storage.Recording{
	// 	ID:          15,
	// 	Calculation: "1+1=5",
	// 	CreatedAt:   "2025-07-18 20:50:00",
	// }
	// storage.CreatrRecording(save)

	// PostgreSql
	// if !mysql.TestConnection() {
	// 	panic("Error conection!!!")
	// }
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL_MYSQL"))
	// if err != nil {
	// 	log.Fatalf("Failed to connect to the database: %v", err)
	// }
	// defer conn.Close(context.Background())
	// rows, err := conn.Query(context.Background(), "SELECT * from historys")
	// if err != nil {
	// 	log.Fatalf("Error GET recording: %v", err)
	// }
	// defer conn.Close(context.Background())

	// MySql
	if !mysql.TestConnection() {
		panic("Error conection!!!")
	}

	conn, err := mysql.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return []*types.Recording{}
	}

	var rows *sql.Rows
	var errG error
	if key == "" {
		rows, errG = conn.Query("SELECT * from historys")
		if errG != nil {
			log.Fatalln("error GET recording: ", errG)
		}
		defer mysql.Deconect(conn)
	} else {
		numInt, err := strconv.Atoi(key)
		if err != nil {
			log.Fatalln("error converting string to int: ", err)
			return []*types.Recording{}
		}

		rows, errG = conn.Query("SELECT * from historys where id = ?", numInt)
		if errG != nil {
			log.Fatalln("error GET recording: ", errG)
		}
		defer mysql.Deconect(conn)
	}

	var historesArray []*types.Recording

	for rows.Next() {
		histore := &types.Recording{}

		err := rows.Scan(&histore.ID, &histore.Calculation, &histore.CreatedAt)
		if err != nil {
			log.Fatalln("failed to get entry history: ", err)
			continue
		}

		storage.CreatrRecording(&types.Recording{
			ID:          histore.ID,
			Calculation: histore.Calculation,
			CreatedAt:   histore.CreatedAt,
		})

		historesArray = append(historesArray, histore)
	}

	return historesArray
}
