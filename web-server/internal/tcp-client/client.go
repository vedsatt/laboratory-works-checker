package tcpclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type LabWork struct {
	ID    int    `json:"id"`
	Code  string `json:"code"`
	Task1 int    `json:"task1"`
	Task2 int    `json:"task2"`
	Task3 int    `json:"task3"`
}

type CheckerResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func Connect() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Устанавливаем таймауты
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	labWork := LabWork{
		ID:    1,
		Code:  "package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}",
		Task1: 5,
		Task2: 3,
		Task3: 4,
	}

	jsonData, err := json.Marshal(labWork)
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return
	}

	// Отправляем данные с \n
	_, err = conn.Write(append(jsonData, '\n'))
	if err != nil {
		fmt.Println("Error sending:", err)
		return
	}

	// Читаем ответ
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	// Десериализуем
	var result CheckerResponse
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		fmt.Println("Error unmarshaling:", err)
		return
	}

	fmt.Printf("Response: %+v\n", result)
}
