package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	cfg            Config
	listener       net.Listener
	router         *Router
	wg             sync.WaitGroup
	shutdownCtx    context.Context
	cancelShutdown context.CancelFunc
}

type Config struct {
	Network string
	Port    int
}

func NewServer(cfg Config, router *Router) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		cfg:            cfg,
		router:         router,
		wg:             sync.WaitGroup{},
		shutdownCtx:    ctx,
		cancelShutdown: cancel,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen(s.cfg.Network, fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.listener = listener

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

func (s *Server) Shutdown() error {
	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	s.cancelShutdown()

	select {
	case <-done:
		fmt.Println("server stopped")
	case <-time.After(time.Second * 5):
		return fmt.Errorf("graceful shutdown timeout")
	}
	return nil
}
