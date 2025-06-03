// Package models - структуры, которые используются в разных пакетах
package models

type (
	// Представляет собой структуру для принятия лабораторных работ
	//	Поля:
	// 	ID: уникальный айди запроса, который генерируется рандомно
	// 	LabNum: номер лабораторной работы
	// 	Code: стркоа с кодом с преобразованными escape-последовательностями
	// 	Tasks: мапа с вариациями заданий. Ключ - номер части, значение - вариация этой части для конкретного варианта
	// 	CurrTask int: в случае с разбитыми лабораторными работами отправляется номер конкретной части (1-3)
	// }
	// Пример:
	// ID: 124534
	// LabNum: 1
	// Code: #include <stdio.h>\\n\\nint main() {\\n    printf("Hello world");\\n    return 0;\\n}"
	// Tasks: {"task1":4, "task2":1, "task3":9}
	//	CurrTask: 0
	//
	LabRequest struct {
		ID       int            `json:"id"`
		LabNum   int            `json:"lab_number"`
		Code     string         `json:"code"`
		Tasks    map[string]int `json:"tasks"`
		CurrTask int            `json:"task"`
	}

	// Структура для отправки ответа
	// Response struct {
	// 	ID: уникальный айди запроса
	// 	ResMsg: результат проверки (если все корректно - "ОК", в ином случае - конкретное сообщение об ошибке)
	// }
	Response struct {
		ID     int    `json:"id"`
		ResMsg string `json:"res_msg"`
	}
)
