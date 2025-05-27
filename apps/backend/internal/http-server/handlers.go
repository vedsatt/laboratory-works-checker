package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/vedsatt/laboratory-works-checker/apps/backend/internal/checker"
	"github.com/vedsatt/laboratory-works-checker/apps/backend/internal/models"
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

	resp := &models.CheckerResponse{}
	resp.ID = body.ID

	checker, err := checker.New(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	}
	msg, err := checker.Check()
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	}
	resp.ResMsg = msg

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
