package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// Server represents the HTTP server and manages its lifecycle,
// TCP connections, and graceful shutdown.
type Server struct {
	cfg      Config
	listener net.Listener
	router   *Router
	// wg tracks all active server connections.
	wg sync.WaitGroup
	// shutdownCtx is used to notify active connections
	// that the server is shutting down.
	shutdownCtx context.Context
	// ready is closed once the listener has been initialized.
	ready chan struct{}
	// cancelShutdown cancels shutdownCtx and signals active connections
	// to shut down.
	cancelShutdown context.CancelFunc
}

// Config contains the TCP server configuration.
type Config struct {
	Network string
	Port    int
}

func (s *Server) Ready() <-chan struct{} {
	return s.ready
}

// NewServer creates a new Server with the specified configuration and router.
func NewServer(cfg Config, router *Router) *Server {
	// Create the server lifecycle context.
	// It will be cancelled during graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		cfg:            cfg,
		router:         router,
		wg:             sync.WaitGroup{},
		ready:          make(chan struct{}),
		shutdownCtx:    ctx,
		cancelShutdown: cancel,
	}
}

// Run starts the server and begins accepting incoming TCP connections.
//
// The method blocks until the listener is closed.
func (s *Server) Run() error {
	listener, err := net.Listen(s.cfg.Network, fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.listener = listener
	close(s.ready)

	// Wait for new TCP connections.
	// The loop terminates when the listener is closed during shutdown.
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			continue
		}
		s.wg.Go(func() { s.serveConn(s.shutdownCtx, conn) })
	}

	return nil
}

// Shutdown performs a graceful server shutdown.
//
// The server stops accepting new connections and waits for
// already active connections to finish. If they do not finish
// within 5 seconds, a timeout error is returned.
func (s *Server) Shutdown() error {
	// Close the listener to stop accepting new connections.
	// This also unblocks Accept() in Run().
	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	// Wait for all active connections to finish in a separate goroutine
	// so that a maximum wait time can be enforced.
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	// Cancel the server context and notify active connections
	// that they should terminate.
	s.cancelShutdown()

	select {
	case <-done:
		fmt.Println("server stopped")
	case <-time.After(time.Second * 5):
		return fmt.Errorf("graceful shutdown timeout")
	}
	return nil
}
