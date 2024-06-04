package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	// Set up listener on port 4221 using tcp
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	// Defer connection closing before exiting
	defer l.Close()

	// Accept incoming connection requests
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// Set up buffer size and read incoming data
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading connection", err.Error())
		os.Exit(1)
	}

	// Make string from buf byte array
	req := string(buf)

	// Split the request in parts using the defined separator
	lines := strings.Split(req, "\r\n\r\n")

	// Split the first part to get the path
	path := strings.Split(lines[0], " ")[1]
	// fmt.Println(path)

	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
