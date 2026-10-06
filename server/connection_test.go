package server

import (
	"bufio"
	"context"
	"net"
	"testing"
)

func TestServeConn(t *testing.T) {
	router := NewRouter()

	handler := func(r *Request, w ResponseWriter) {
		w.SetBody([]byte("hello"))
	}

	router.SetRoute("GET", "/", handler)

	s := &Server{
		router: router,
	}

	client, conn := net.Pipe()
	defer client.Close()

	go s.serveConn(context.Background(), conn)

	request := "GET / HTTP/1.1\r\nConnection: close\r\n\r\n"

	_, err := client.Write([]byte(request))
	if err != nil {
		t.Fatalf("failed to write request: %v", err)
	}

	reader := bufio.NewReader(client)

	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if line != "HTTP/1.1 200 OK\r\n" {
		t.Errorf("unexpected status line: %q", line)
	}
}

func TestServeConnRouteNotFound(t *testing.T) {
	router := NewRouter()

	s := &Server{
		router: router,
	}

	client, conn := net.Pipe()
	defer client.Close()

	go s.serveConn(context.Background(), conn)

	request := "GET /unknown HTTP/1.1\r\nConnection: close\r\n\r\n"

	_, err := client.Write([]byte(request))
	if err != nil {
		t.Fatalf("failed to write request: %v", err)
	}

	buf := make([]byte, 1)
	_, err = client.Read(buf)

	if err == nil {
		t.Error("expected connection to be closed")
	}
}

func TestServeConnInvalidRequest(t *testing.T) {
	router := NewRouter()

	s := &Server{
		router: router,
	}

	client, conn := net.Pipe()
	defer client.Close()

	go s.serveConn(context.Background(), conn)

	request := "GET / HTTP/1.0\r\n\r\n"

	_, err := client.Write([]byte(request))
	if err != nil {
		t.Fatalf("failed to write request: %v", err)
	}

	buf := make([]byte, 1)
	_, err = client.Read(buf)

	if err == nil {
		t.Error("expected connection to be closed")
	}
}
