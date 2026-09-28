package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	cfg      Config
	listener net.Listener
	router   *Router
	wg       *sync.WaitGroup
}

type Config struct {
	Network string
	Port    int
}

func NewServer(cfg Config, router *Router, wg *sync.WaitGroup) *Server {
	return &Server{
		cfg:    cfg,
		router: router,
		wg:     wg,
	}
}

func (s *Server) run() error {
	listener, err := net.Listen(s.cfg.Network, fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.listener = listener

	sig, stopSig := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSig()

	shutdownCtx, cancelShutdownCtx := context.WithCancel(context.Background())
	defer cancelShutdownCtx()

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					break
				}
				continue
			}
			s.wg.Go(func() { s.serveConn(shutdownCtx, conn) })
		}
	}()

	<-sig.Done()
	stopSig()

	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	cancelShutdownCtx()

	select {
	case <-done:
		fmt.Println("server stopped")
	case <-time.After(time.Second * 5):
		return fmt.Errorf("graceful shutdown timeout")
	}
	return nil
}

func (s *Server) MustRun() {
	err := s.run()
	if err != nil {
		panic(err)
	}
}

func (s *Server) Shutdown() error
