package server

import (
	"bufio"
	"fmt"
	"net"
)

func Server() {
	listenner, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	defer listenner.Close()

	fmt.Println("tcp server is started")

	for {
		conn, err := listenner.Accept()
		if err != nil {
			fmt.Printf("Failed connection: %s", err.Error())
			continue
		}

		go getConnection(conn)
	}
}

func getConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	res, _ := parseRequest(reader)
	fmt.Println(string(res.Body))
}
