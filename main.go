package main

import (
	"Trial-microservices-History-go/src"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	http.HandleFunc("/histore", historeHandler)
	http.HandleFunc("/calculator", calculatorHandler)

	fmt.Println("Starting server at port 8080")
	errS := http.ListenAndServe(":8080", nil)
	if errS != nil {
		fmt.Println("Error starting the server:", errS)
	}
}

func historeHandler(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "id")

	queryParams := r.URL.Query()
	// Получаем значение параметра "name"
	idStr := queryParams.Get("id")

	// if idStr == "" {
	// 	render.Status(r, http.StatusBadRequest)
	// 	return
	// }

	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	render.Status(r, http.StatusBadRequest)
	// 	return
	// }

	result := src.GetHistory(idStr)

	render.Status(r, 200)
	render.JSON(w, r, result)

}

type Recording struct {
	Calculation string `json:"Calculation"`
	CreatedAt   string `json:"CreatedAt"`
}

func calculatorHandler(w http.ResponseWriter, r *http.Request) {
	// Read request body
	bodyByt, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON into struct
	var record Recording
	err = json.Unmarshal(bodyByt, &record)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 2025-07-18 20:50:00.623091+00

	// Validate required fields
	if record.Calculation == "" || record.CreatedAt == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Save the record
	err = src.SaveHistory(record.Calculation, record.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to save record", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Recording saved successfully",
	})
}
