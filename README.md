# laboratory-works-checker

## О проекте
Проект представляет собой сервер для проверки лабораторных работ по направлению "Алгоритмизация и программирование" первого курса ИВТ ВШЭ.


## Используемые технологии
- Golang
- JS/CSS/HTML
- JSON
- Python
- REST API
- TLS
- Docker & Docker Compose

## Структура проекта
```
.
├── README.md
├── apps
│   ├── backend
│   │   ├── application
│   │   │   ├── checker
│   │   │   │   ├── checker.go
│   │   │   │   ├── monolit.go
│   │   │   │   ├── splited.go
│   │   │   │   └── validator.go
│   │   │   ├── cleaner
│   │   │   │   └── cleaner.go
│   │   │   ├── config
│   │   │   │   └── config.go
│   │   │   ├── http-server
│   │   │   │   ├── handlers.go
│   │   │   │   ├── middlewares.go
│   │   │   │   └── server.go
│   │   │   └── models
│   │   │       └── models.go
│   │   └── cmd
│   │       └── cmd.go
│   └── frontend
│       ├── index.html
│       ├── lab.html
│       ├── script.js
│       └── styles.css
├── labs
│   ├── lab1
│   │   ├── config.json
│   │   ├── description.md
│   │   ├── task1
│   │   │   ├── var1
│   │   │   │   └── solution.py
│   │   │   ├── var2
│   │   │   │   └── solution.py
│   │   │   ├── var3
│   │   │   │   └── solution.py
│   │   │   └── var5
│   │   │       └── solution.py
│   │   ├── task2
│   │   │   ├── var1
│   │   │   │   └── solution.py
│   │   │   ├── var2
│   │   │   │   └── solution.py
│   │   │   ├── var3
│   │   │   │   └── solution.py
│   │   │   └── var8
│   │   │       └── solution.py
│   │   └── task3
│   │       ├── var1
│   │       │   └── solution.py
│   │       ├── var2
│   │       │   └── solution.py
│   │       ├── var3
│   │       │   └── solution.py
│   │       └── var4
│   │           └── solution.py
│   ├── lab10
│   │   └── description.md
│   ├── lab11
│   │   └── description.md
│   ├── lab2
│   │   ├── config.json
│   │   ├── description.md
│   │   ├── task1
│   │   │   ├── var1
│   │   │   │   └── tests.json
│   │   │   └── var2
│   │   │       └── tests.json
│   │   └── task2
│   │       ├── var1
│   │       │   └── tests.json
│   │       └── var2
│   │           └── tests.json
│   ├── lab3
│   │   └── description.md
│   ├── lab4
│   │   └── description.md
│   ├── lab5
│   │   └── description.md
│   ├── lab6
│   │   └── description.md
│   ├── lab7
│   │   └── description.md
│   ├── lab8
│   │   └── description.md
│   └── lab9
│       └── description.md
└── go.mod
```

## Требования
Перед запуском убедитесь, что у вас установлено:
### Вариант с Docker (рекомендуется)
- **Docker** версия 24.0 или выше  
  [Скачать Docker](https://www.docker.com/get-started)
- **Docker Compose** (обычно входит в Docker)
### Вариант без Docker

- **Go** версия 1.20+  
  [Скачать Go](https://golang.org/dl/)
- **Интерпретатор Python**  
  [Скачать интерпретатор Python](https://www.python.org/downloads/)


## Установка и запуск проекта
### Установка
1. Клонируйте репозиторий
```sh
git clone https://github.com/vedsatt/laboratory-works-checker.git
cd ./laboratory-works-checker
```
Установите зависимости
```sh
go mod tidy
```
### Запуск
1. **Укажите переменные среды (легче всего через Docker)**:
```sh
PORT=<ЗНАЧЕНИЕ>
HTTPS_PORT=<ЗНАЧЕНИЕ>
HOST_URL=<ЗНАЧЕНИЕ> // "localhost" без Docker; "backend" с Docker
```
**Примечание:** если не создать файл и/или не указать определенные значения, то программа установит неуказанные значения по умолчанию (для адреса подключения к оркестратору будет установлено значение для локального подключения)

### С помощью Go:

1. **Запустите backend-сервер:**
```sh
go run apps/backend/cmd/cmd.go
```
2. **Запустите frontend-сервер:**
```
тут команды для запуска
```
### С помощью Docker:
**Разверните контейнеры с помощью docker-compose:**
```sh
docker-compose up --build
```

## Дизайн системы
В коде подробна расписана каждая функция backend-а.
Чтобы посмотреть структурированную документацию:
1. Установите godoc:
```bash
go install golang.org/x/tools/cmd/godoc@latest
```
2. Запустите локальный сайт с документацией:
```bash
godoc -http=:6060
```
3. Перейдите на этот сайт: http://localhost:6060/pkg/github.com/vedsatt/laboratory-works-checker
## Контакты
Если у вас возникли вопросы или предложения, можете написать их в issue или написать мне лично в тг  
[![Telegram](https://raw.githubusercontent.com/CLorant/readme-social-icons/refs/heads/main/small/filled/telegram.svg)](https://t.me/sigmatemik52) [**Telegram**: @sigmatemik52](https://t.me/sigmatemik52)
