package server

import (
	"bufio"
	"fmt"
	"net"
)

type Server struct {
	listener net.Listener
	Cfg      Config
}

type Config struct {
	Network string
	Port    int
}

func New(cfg Config) *Server {
	return &Server{
		Cfg: cfg,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen(s.Cfg.Network, fmt.Sprintf(":%d", s.Cfg.Port))
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

		go worker(conn)
	}

}

func worker(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	req, err := parseRequest(reader)
	if err != nil {
		return
	}

	resp := NewResponse(req.Protocol).
		SetStatus(200).
		SetBody(req.Body)

	resp.write(conn)
}
