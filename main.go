package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	fmt.Println("Listening on TCP 5000")

	l, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	// Close the connection
	defer l.Close()

	for {

		// accept connections to our tcp server
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		go handleUserConnection(conn)

	}

}

func handleUserConnection(c net.Conn) {
	defer c.Close()

	for {
		input, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v", err)
			break
		}
		sendBackToClient(c)
		fmt.Println(input)
	}

}

func sendBackToClient(c net.Conn) {

	// lets send some data back over the connection
	if _, err := io.WriteString(c, "oh no you didnt\n"); err != nil {
		fmt.Printf("Error: %v", err)
	}

}
