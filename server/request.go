package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	ErrConnectionClosed = errors.New("connection closed")
)

// Request represents an HTTP request.
type Request struct {
	Method        string
	RequestTarget string
	Protocol      string
	Headers       map[string]string
	Body          []byte
}

// parseStartLine parses the HTTP request start line.
// Returns the method, request-target, and protocol.
func parseStartLine(startLine string) (string, string, string, error) {
	startLine = strings.TrimSpace(startLine)
	parts := strings.Split(startLine, " ")

	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid start line: %s", startLine)
	}

	if parts[2] != "HTTP/1.1" {
		return "", "", "", fmt.Errorf("invalid protocol: %s", startLine)
	}

	return parts[0], parts[1], parts[2], nil
}

// parseHeaders reads HTTP headers until an empty line.
func parseHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		// The empty line indicates the end of the HTTP headers.
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

// parseBody reads exactly contentLength bytes from the HTTP request body.
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

// parseRequest parses an HTTP request from the stream.
func parseRequest(reader *bufio.Reader) (*Request, error) {
	startLine, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) && len(startLine) == 0 {
			return nil, ErrConnectionClosed
		}

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

	// Content-Length specifies the length of the request body, if present.
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
