package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vedsatt/laboratory-works-checker/internal/models"
)

func SubmitLabHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method, use the POST method", http.StatusMethodNotAllowed)
		log.Println("invalid request method")
		return
	}

	var body *models.LabRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "error with decoding request body", http.StatusInternalServerError)
		log.Printf("error with decoding request body: %v", err)
		return
	}

	// fmt.Printf("code:\n%v\ntask1: %v\ntask2: %v\ntask3: %v\n", body.Code, body.Task1, body.Task2, body.Task3)
}
