package server

import (
	"net"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	cfg := Config{
		Network: "tcp",
		Port:    8080,
	}

	router := NewRouter()

	s := NewServer(cfg, router)

	if s.cfg != cfg {
		t.Errorf("unexpected config: got %+v, want %+v", s.cfg, cfg)
	}

	if s.router != router {
		t.Error("router was not set correctly")
	}

	if s.shutdownCtx == nil {
		t.Error("shutdownCtx is nil")
	}

	if s.cancelShutdown == nil {
		t.Error("cancelShutdown is nil")
	}
}

func TestServerRun(t *testing.T) {
	router := NewRouter()

	s := NewServer(Config{
		Network: "tcp",
		Port:    0,
	}, router)

	runErr := make(chan error, 1)

	go func() {
		runErr <- s.Run()
	}()

	<- s.Ready()

	defer s.listener.Close()

	conn, err := net.Dial("tcp", s.listener.Addr().String())
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
	conn.Close()

	if err := s.Shutdown(); err != nil {
		t.Fatalf("failed to shutdown server: %v", err)
	}

	select {
	case err := <-runErr:
		if err != nil {
			t.Errorf("unexpected Run error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Run did not stop after shutdown")
	}
}
