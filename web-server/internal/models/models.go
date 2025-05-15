package models

type (
	Request struct {
		ID    int    `json:"id"`
		Code  string `json:"code"`
		Task1 int    `json:"task1"`
		Task2 int    `json:"task2"`
		Task3 int    `json:"task3"`
	}

	LabWork struct {
		ID    int    `json:"id"`
		Code  string `json:"code"`
		Task1 int    `json:"task1"`
		Task2 int    `json:"task2"`
		Task3 int    `json:"task3"`
	}

	CheckerResponse struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
)
