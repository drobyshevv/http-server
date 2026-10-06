package server

import (
	"bufio"
	"errors"
	"strings"
	"testing"
)

func TestParseStartLine(t *testing.T) {
	input := "GET /users HTTP/1.1"
	method, requestTarget, protocol, err := parseStartLine(input)

	if method != "GET" {
		t.Errorf("wrong method: %s", method)
	}
	if requestTarget != "/users" {
		t.Errorf("wrong request target: %s", requestTarget)
	}
	if protocol != "HTTP/1.1" {
		t.Errorf("wrong protocol: %s", protocol)
	}
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestParseStartLineInvalid(t *testing.T) {
	input := "GET /users"
	_, _, _, err := parseStartLine(input)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParseStartLineSuperfluousData(t *testing.T) {
	input := "GET /users HTTP/1.1 something"
	_, _, _, err := parseStartLine(input)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParseStartLineMultipleSpaces(t *testing.T) {
	input := "GET  /users HTTP/1.1"
	_, _, _, err := parseStartLine(input)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParseStartLineInvalidProtocol(t *testing.T) {
	input := "GET /users FTP/1.1"
	_, _, _, err := parseStartLine(input)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParseHeaders(t *testing.T) {
	input := "Host: example.com\r\n" +
		"Content-Type: application/json\r\n" +
		"\r\n"

	reader := strings.NewReader(input)

	headers, err := parseHeaders(bufio.NewReader(reader))

	if headers["Host"] != "example.com" || headers["Content-Type"] != "application/json" {
		t.Errorf("wrong request headers: %s", headers)
	}
	if len(headers) != 2 {
		t.Errorf("wrong request headers: %s", headers)
	}
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestParseHeadersEmpty(t *testing.T) {
	input := "\r\n"

	reader := strings.NewReader(input)

	headers, err := parseHeaders(bufio.NewReader(reader))

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if len(headers) != 0 {
		t.Errorf("expected empty headers, got %v", headers)
	}
}

func TestParseHeadersInvalid(t *testing.T) {
	input := "Host example.com\r\n\r\n"

	reader := strings.NewReader(input)

	_, err := parseHeaders(bufio.NewReader(reader))

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParseHeadersTrimSpaces(t *testing.T) {
	input := "Host:   example.com   \r\n" +
		"Content-Type:   application/json   \r\n" +
		"\r\n"

	reader := bufio.NewReader(strings.NewReader(input))

	headers, err := parseHeaders(reader)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if headers["Host"] != "example.com" {
		t.Errorf("wrong Host header: %q", headers["Host"])
	}

	if headers["Content-Type"] != "application/json" {
		t.Errorf("wrong Content-Type header: %q", headers["Content-Type"])
	}
}

func TestParseHeadersValueContainsColon(t *testing.T) {
	input := "Host: example.com:8080\r\n\r\n"

	reader := strings.NewReader(input)

	headers, err := parseHeaders(bufio.NewReader(reader))

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if headers["Host"] != "example.com:8080" {
		t.Errorf("wrong Host header: %q", headers["Host"])
	}
}

func TestParseBody(t *testing.T) {
	input := "hello world"
	reader := bufio.NewReader(strings.NewReader(input))

	body, err := parseBody(reader, 11)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if string(body) != "hello world" {
		t.Errorf("wrong body: %s", body)
	}
}

func TestParseBodyEmpty(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(""))

	body, err := parseBody(reader, 0)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if len(body) != 0 {
		t.Errorf("expected empty body, got %v", body)
	}
}

func TestParseBodyUnexpectedEOF(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("hello"))

	_, err := parseBody(reader, 10)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParseBodyReadsExactLength(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("hello world"))

	body, err := parseBody(reader, 5)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if string(body) != "hello" {
		t.Errorf("wrong body: %s", body)
	}
}

func TestParseRequest(t *testing.T) {
	input := "GET /users HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Content-Type: application/json\r\n" +
		"\r\n"

	reader := bufio.NewReader(strings.NewReader(input))

	req, err := parseRequest(reader)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if req.Method != "GET" {
		t.Errorf("wrong method: %s", req.Method)
	}

	if req.RequestTarget != "/users" {
		t.Errorf("wrong request target: %s", req.RequestTarget)
	}

	if req.Protocol != "HTTP/1.1" {
		t.Errorf("wrong protocol: %s", req.Protocol)
	}

	if req.Headers["Host"] != "example.com" {
		t.Errorf("wrong Host header: %s", req.Headers["Host"])
	}

	if req.Body != nil {
		t.Errorf("expected nil body, got %v", req.Body)
	}
}

func TestParseRequestWithBody(t *testing.T) {
	input := "POST /users HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Content-Length: 5\r\n" +
		"\r\n" +
		"hello"

	reader := bufio.NewReader(strings.NewReader(input))

	req, err := parseRequest(reader)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if string(req.Body) != "hello" {
		t.Errorf("wrong body: %s", req.Body)
	}
}

func TestParseRequestConnectionClosed(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(""))

	_, err := parseRequest(reader)

	if !errors.Is(err, ErrConnectionClosed) {
		t.Errorf("expected ErrConnectionClosed, got %v", err)
	}
}

func TestParseRequestInvalidContentLength(t *testing.T) {
	input := "POST /users HTTP/1.1\r\n" +
		"Content-Length: abc\r\n" +
		"\r\n"

	reader := bufio.NewReader(strings.NewReader(input))

	_, err := parseRequest(reader)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
