package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/drobyshevv/http-server/server"
)

func main() {
	router := server.NewRouter()
	router.SetRoute("POST", "/users", handler)

	cfg := server.Config{
		Network: "tcp",
		Port:    8080,
	}

	srv := server.NewServer(cfg, router)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Run()
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	if err := srv.Shutdown(); err != nil {
		fmt.Printf("shutdown error: %v\n", err)
	}

	if err := <-errCh; err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}

func handler(r *server.Request, w server.ResponseWriter) {
	w.SetStatus(201)
	w.SetBody(r.Body)
}
