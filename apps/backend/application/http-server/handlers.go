package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/vedsatt/laboratory-works-checker/apps/backend/application/checker"
	"github.com/vedsatt/laboratory-works-checker/apps/backend/application/models"
)

func SubmitLabHandler(w http.ResponseWriter, r *http.Request) {
	// Пропускается только метод POST
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("invalid request method: %v, use the POST method", r.Method), http.StatusMethodNotAllowed)
		log.Println("invalid request method")
		return
	}

	// Декодируем тело запроса в структуру
	var body *models.LabRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "system error with decoding request body", http.StatusInternalServerError)
		log.Printf("System error with decoding request body: %v", err)
		return
	}

	// Формируем структуру ответа с айди запроса
	resp := &models.Response{}
	resp.ID = body.ID

	// Создаем новый объект Checker
	checker, err := checker.New(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		log.Printf("System error with creating checker obj: %v", err)
		return
	}
	log.Printf("Request with id: %v, labNum %v, tasks: %v", body.ID, body.LabNum, body.Tasks)

	// Отправляем работу на проверку
	msg, err := checker.Check()
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		log.Printf("System error with checker: %v", err)
		return
	}
	resp.ResMsg = msg

	// Возвращаем ответ
	log.Printf("Response was sent")
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
