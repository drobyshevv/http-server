# HTTP Server

[English](README.md) | [Русский](README.ru.md)

Небольшой HTTP/1.1 сервер, реализованный на Go только с использованием стандартной библиотеки.

Проект реализует базовый жизненный цикл HTTP-сервера: приём TCP-соединений, разбор HTTP-запросов, маршрутизацию запросов к обработчикам и формирование HTTP-ответов.

## Возможности

* TCP-сервер на стандартной библиотеке Go
* Разбор HTTP/1.1 запросов
* Разбор стартовой строки запроса
* Разбор HTTP-заголовков
* Разбор тела запроса через `Content-Length`
* Формирование HTTP-ответов
* HTTP status codes и headers
* Простой роутер по методу и пути
* Поддержка persistent TCP connections
* Тайм-ауты соединения, чтения и записи
* Graceful shutdown
* Unit и integration-style тесты

## Структура проекта

```text
.
├── client/
│   └── main.go
├── server/
│   ├── connection.go
│   ├── connection_test.go
│   ├── request.go
│   ├── request_test.go
│   ├── response.go
│   ├── response_test.go
│   ├── router.go
│   ├── router_test.go
│   ├── server.go
│   └── server_test.go
├── main.go
├── go.mod
├── README.md
└── README.ru.md
```

### `server/`

Содержит реализацию HTTP-сервера:

* `request.go` — разбор HTTP-запросов
* `response.go` — формирование и запись HTTP-ответов
* `router.go` — маршрутизация запросов
* `connection.go` — обработка TCP-соединений
* `server.go` — жизненный цикл сервера и graceful shutdown
* `*_test.go` — тесты компонентов сервера

### `client/`

Содержит небольшой TCP-клиент для отправки raw HTTP-запроса на сервер.

### `main.go`

Содержит точку входа приложения и конфигурацию сервера.

## Запуск

Запустить сервер:

```bash
go run .
```

Сервер слушает:

```text
localhost:8080
```

В другом терминале запустить пример клиента:

```bash
go run ./client
```

## Пример запроса

Пример клиента отправляет raw HTTP-запрос:

```http
POST /users HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 24
Connection: close

{"name":"John","age":25}
```

Сервер отвечает:

```http
HTTP/1.1 201 Created
Content-Length: 24

{"name":"John","age":25}
```

## Тесты

Запустить все тесты:

```bash
go test ./...
```

Запустить тесты пакета `server`:

```bash
go test ./server
```

Также можно запустить статический анализ:

```bash
go vet ./...
```

## Требования

* Go 1.25+
* Внешние зависимости отсутствуют
