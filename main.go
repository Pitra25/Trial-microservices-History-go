package main

import (
	"Trial-microservices-History-go/src"
	"Trial-microservices-History-go/src/types"
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

	http.HandleFunc("/history", historeHandler)
	http.HandleFunc("/save", calculatorHandler)

	errS := http.ListenAndServe(":8080", nil)
	if errS != nil {
		fmt.Println("Error starting the server:", errS)
	}

	fmt.Println("Starting server at port 8080")

}

func historeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(404)
		w.Write([]byte("Post method NOT"))
		return
	}

	queryParams := r.URL.Query()
	idStr := queryParams.Get("id")

	result := src.GetHistory(idStr)

	render.Status(r, 200)
	render.JSON(w, r, result)

}

func calculatorHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(404)
		w.Write([]byte("Get method NOT"))
		return
	}

	// Read request body
	bodyByt, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON into struct
	var record *types.BodyStructure
	err = json.Unmarshal(bodyByt, &record)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

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
