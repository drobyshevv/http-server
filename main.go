package main

import "github.com/drobyshevv/http-server/server"

func main() {
	cfg := server.Config{
		Network: "tcp",
		Port:    8080,
	}
	server := server.New(cfg)
	server.Run()
}
