package server

import (
	"bufio"
	"fmt"
	"net"
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

	req, err := parseRequest(reader)
	if err != nil {
		return
	}

	resp := NewResponse(req.Protocol)
	if err := s.router.Route(req, resp); err != nil {
		return
	}
	resp.write(conn)
}
