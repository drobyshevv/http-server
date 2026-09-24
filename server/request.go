package server

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Request struct {
	Method        string
	RequestTarget string
	Protocol      string
	Headers       map[string]string
	Body          []byte
}

func parseStartLine(startLine string) (string, string, string, error) {
	startLine = strings.TrimSpace(startLine)
	parts := strings.SplitN(startLine, " ", 3)

	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("invalid start line %s", startLine)
	}

	return parts[0], parts[1], parts[2], nil
}

func parseHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if line == "\r\n" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		headers[key] = value
	}

	return headers, nil
}

func parseBody(reader *bufio.Reader, contentLength int) ([]byte, error) {
	body := make([]byte, contentLength)
	remaining := contentLength
	for remaining > 0 {
		n, err := reader.Read(body)
		if err != nil {
			return nil, err
		}

		remaining -= n
	}
	return body, nil
}

func parseRequest(reader *bufio.Reader) (*Request, error) {
	startLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read start line %w", err)
	}

	method, reqTarget, protocol, err := parseStartLine(startLine)
	if err != nil {
		return nil, err
	}

	headers, err := parseHeaders(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse headers %w", err)
	}

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
