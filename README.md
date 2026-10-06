# HTTP Server

[English](README.md) | [Русский](README.ru.md)

A small HTTP/1.1 server implemented in Go using only the standard library.

The project implements the basic HTTP server lifecycle: accepting TCP connections, parsing HTTP requests, routing requests to handlers, and writing HTTP responses.

## Features

* TCP server built with the Go standard library
* HTTP/1.1 request parsing
* Request start-line parsing
* HTTP header parsing
* Request body parsing with `Content-Length`
* HTTP response construction
* HTTP status codes and headers
* Simple method + path router
* Persistent TCP connections
* Connection and read/write timeouts
* Graceful shutdown
* Unit and integration-style tests

## Project Structure

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

Contains the HTTP server implementation:

* `request.go` — HTTP request parsing
* `response.go` — HTTP response construction and writing
* `router.go` — request routing
* `connection.go` — TCP connection handling
* `server.go` — server lifecycle and graceful shutdown
* `*_test.go` — tests for the server components

### `client/`

Contains a small TCP client used to send a raw HTTP request to the server.

### `main.go`

Contains the application entry point and server configuration.

## Run

Start the server:

```bash
go run .
```

The server listens on:

```text
localhost:8080
```

Run the example client in another terminal:

```bash
go run ./client
```

## Example Request

The example client sends a raw HTTP request:

```http
POST /users HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 24
Connection: close

{"name":"John","age":25}
```

The server responds with:

```http
HTTP/1.1 201 Created
Content-Length: 24

{"name":"John","age":25}
```

## Tests

Run all tests:

```bash
go test ./...
```

Run tests for the server package:

```bash
go test ./server
```

You can also run static analysis with:

```bash
go vet ./...
```

## Requirements

* Go 1.25+
* No external dependencies
