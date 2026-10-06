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

// serveConn handles a single TCP connection.
//
// The method sequentially receives HTTP requests, passes them to the router,
// and writes HTTP responses to the connection. The connection is terminated
// when the client closes it, a timeout occurs, an error is encountered,
// or the server context is cancelled.
func (s *Server) serveConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	// Close the connection when the server context is cancelled
	// to interrupt a blocked Read.
	stop := context.AfterFunc(ctx, func() {
		conn.Close()
	})
	defer stop()

	// Prevent a panic in a single connection handler
	// from terminating the entire server.
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

		// If the client explicitly sends Connection: close, close the connection
		// after sending the current response.
		if req.Headers["Connection"] == "close" {
			fmt.Println(req.Headers["Connection"])
			return
		}
	}
}
