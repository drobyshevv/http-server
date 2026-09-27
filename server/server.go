package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"
)

type Server struct {
	listener net.Listener
	router   *Router
	cfg      Config
}

type Config struct {
	Network string
	Port    int
}

func NewServer(cfg Config, router *Router) *Server {
	return &Server{
		cfg:    cfg,
		router: router,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen(s.cfg.Network, fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	defer listener.Close()
	s.listener = listener

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		go s.worker(conn)
	}

}

func (s *Server) worker(conn net.Conn) {
	defer conn.Close()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic in worker: %v\n", r)
		}
	}()

	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		req, err := parseRequest(reader)
		if err != nil {
			if errors.Is(err, ErrConnectionClosed) {
				fmt.Println("connection closed")
				return
			}
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				fmt.Println("timeout:", err)
				return
			}
			fmt.Println("request error:", err)
			return
		}

		resp := NewResponse(req.Protocol)
		if err := s.router.Route(req, resp); err != nil {
			fmt.Println(err)
			return
		}

		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		err = resp.write(conn)
		if err != nil {
			if errors.Is(err, syscall.ECONNRESET) {
				fmt.Println("connection reset by peer")
				return
			}
			if errors.Is(err, syscall.EPIPE) {
				fmt.Println("broken pipe")
				return
			}
			var netErr *net.OpError
			if errors.As(err, &netErr) {
				fmt.Printf("net error%s: %v\n", netErr.Op, netErr.Err)
				return
			}
			fmt.Println("request error")
			return
		}

		if req.Headers["Connection"] == "close" {
			fmt.Println(req.Headers["Connection"])
			return
		}
	}
}
