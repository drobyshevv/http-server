package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"
)

func (s *Server) serveConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	stop := context.AfterFunc(ctx, func() {
		conn.Close()
	})
	defer stop()

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
