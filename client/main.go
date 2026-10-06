package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	request := "POST /users HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"Content-Type: application/json\r\n" +
		"Content-Length: 24\r\n" +
		"Connection: close\r\n" +
		"\r\n" +
		`{"name":"John","age":25}`

	if _, err := conn.Write([]byte(request)); err != nil {
		panic(err)
	}

	fmt.Println("request sent")

	response := make([]byte, 4096)

	n, err := conn.Read(response)
	if err != nil {
		panic(err)
	}

	fmt.Printf("response:\n%s\n", response[:n])
}
