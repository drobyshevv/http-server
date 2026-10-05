package server

import (
	"net"
	"testing"
)

func TestNewResponse(t *testing.T) {
	resp := NewResponse("HTTP/1.1")

	if resp.Protocol != "HTTP/1.1" {
		t.Errorf("wrong protocol: %s", resp.Protocol)
	}

	if resp.Code != 200 {
		t.Errorf("wrong status code: %d", resp.Code)
	}

	if resp.Header["Content-Length"] != "0" {
		t.Errorf("wrong Content-Length: %s", resp.Header["Content-Length"])
	}

	if resp.Body != nil {
		t.Errorf("expected nil body, got %v", resp.Body)
	}
}

func TestResponseStatus(t *testing.T) {
	resp := NewResponse("HTTP/1.1")

	resp.SetStatus(404)

	if resp.Status() != 404 {
		t.Errorf("wrong status: %d", resp.Status())
	}
}

func TestResponseSetHeader(t *testing.T) {
	resp := NewResponse("HTTP/1.1")

	resp.SetHeader("Content-Type", "application/json")

	if resp.Header["Content-Type"] != "application/json" {
		t.Errorf("wrong Content-Type: %s", resp.Header["Content-Type"])
	}
}

func TestResponseSetBody(t *testing.T) {
	resp := NewResponse("HTTP/1.1")

	body := []byte("hello")
	resp.SetBody(body)

	if string(resp.Body) != "hello" {
		t.Errorf("wrong body: %s", resp.Body)
	}

	if resp.Header["Content-Length"] != "5" {
		t.Errorf("wrong Content-Length: %s", resp.Header["Content-Length"])
	}
}

func TestStatusText(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{200, "OK"},
		{201, "Created"},
		{204, "No Content"},
		{400, "Bad Request"},
		{404, "Not Found"},
		{500, "Internal Server Error"},
	}

	for _, tt := range tests {
		got := statusText(tt.code)

		if got != tt.want {
			t.Errorf("statusText(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestStatusTextUnknown(t *testing.T) {
	if got := statusText(999); got != "" {
		t.Errorf("statusText(999) = %q, want empty string", got)
	}
}

func TestWriteResponcse(t *testing.T) {
	resp := &Response{
		Protocol: "HTTP/1.1",
		Code:     200,
		Header: map[string]string{
			"Content-Length": "5",
		},
		Body: []byte("hello"),
	}

	client, server := net.Pipe()
	defer server.Close()

	go func() {
		defer client.Close()
		resp.write(client)
	}()

	buf := make([]byte, 50)
	n, _ := server.Read(buf)
	str := string(buf[:n])

	if str != "HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\n" {
		t.Errorf("unexpected response: %q", str)
	}

	bufBody := make([]byte, 50)
	nBody, _ := server.Read(bufBody)
	strBody := string(bufBody[:nBody])

	if strBody != "hello" {
		t.Errorf("unexpected body: %q", strBody)
	}
}
