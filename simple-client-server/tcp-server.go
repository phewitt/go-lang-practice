package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("launching server...")

	ln, _ := net.Listen("tcp", ":8081")

	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')

		fmt.Print("Message Received:", string(message))

		newMessage := strings.ToUpper(message)
		conn.Write([]byte(newMessage + "\n"))
	}
}
