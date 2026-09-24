package server

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// Request представляет HTTP-запрос.
type Request struct {
	Method        string
	RequestTarget string
	Protocol      string
	Headers       map[string]string
	Body          []byte
}

// parseStartLine разбирает стартовую строку HTTP-запроса.
// Возвращает method, request-target и protocol.
func parseStartLine(startLine string) (string, string, string, error) {
	startLine = strings.TrimSpace(startLine)
	parts := strings.SplitN(startLine, " ", 3)

	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid start line: %s", startLine)
	}

	return parts[0], parts[1], parts[2], nil
}

// parseHeaders читает HTTP-заголовки до пустой строки.
func parseHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		// Пустая строка обозначает конец HTTP-заголовков.
		if line == "\r\n" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header")
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		headers[key] = value
	}

	return headers, nil
}

// parseBody читает из потока ровно contentLength байт тела HTTP-запроса.
func parseBody(reader *bufio.Reader, contentLength int) ([]byte, error) {
	body := make([]byte, contentLength)
	remaining := contentLength
	for remaining > 0 {
		n, err := reader.Read(body[contentLength-remaining:])
		if err != nil {
			return nil, err
		}

		remaining -= n
	}
	return body, nil
}

// parseRequest разбирает HTTP-запрос из потока.
func parseRequest(reader *bufio.Reader) (*Request, error) {
	startLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read start line: %w", err)
	}

	method, reqTarget, protocol, err := parseStartLine(startLine)
	if err != nil {
		return nil, err
	}

	headers, err := parseHeaders(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse headers: %w", err)
	}

	// Content-Length определяет длину тела, если оно передано.
	val, bodyExists := headers["Content-Length"]

	var body []byte = nil
	if bodyExists {
		contentLength, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Content-Length to integer: %w", err)
		}
		body, err = parseBody(reader, contentLength)
		if err != nil {
			return nil, fmt.Errorf("failed to parse request body: %w", err)
		}
	}

	return &Request{
		Method:        method,
		RequestTarget: reqTarget,
		Protocol:      protocol,
		Headers:       headers,
		Body:          body,
	}, nil
}
