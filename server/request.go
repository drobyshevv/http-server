package server

import (
	"bufio"
	"fmt"
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

func parseHeaders() {

}

func parseBody() {

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
}
