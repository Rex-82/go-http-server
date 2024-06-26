package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	// Set up listener on port 4221 using tcp
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	// Infinite "while" loop to handle multiple connections
	for {

		fmt.Println("waiting for connection...")
		// Accept incoming connection requests
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		// Create new goroutine to handle connections concurrently
		go handleConnection(conn)
	}

}

// Handles each connection individually
func handleConnection(conn net.Conn) {
	// Defer connection closing before exiting
	defer conn.Close()

	// Set up buffer size and read incoming data
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading connection", err.Error())
		os.Exit(1)
	}

	// Make string from buf byte array
	req := string(buf)

	// Split the request in parts using the defined separator
	lines := strings.Split(req, "\r\n\r\n")

	// Split the first part to get path and method
	splitLine := strings.Split(lines[0], " ")
	path := splitLine[1]
	method := splitLine[0]

	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {

		endpoint := strings.Split(path, "/")

		switch endpoint[1] {
		case "echo": // case "echo" -> Server should respond with enpoint content. (eg. echo/hello -> "hello")
			content := "HTTP/1.1 200 OK" +
				"\r\n" +
				"Content-Type: text/plain\r\n" +
				"Content-Length: " + fmt.Sprintf("%d\r\n", utf8.RuneCountInString(endpoint[2])) +
				"\r\n" +
				endpoint[2]

			// fmt.Println(content)

			conn.Write([]byte(content))

		case "user-agent": // case "user-agent" -> Server should respond with user agent infos as the response body
			userAgent := splitLine[4]

			fmt.Println(lines)

			content := "HTTP/1.1 200 OK" +
				"\r\n" +
				"Content-Type: text/plain\r\n" +
				"Content-Length: " + fmt.Sprintf("%d\r\n", utf8.RuneCountInString(userAgent)) +
				"\r\n" +
				userAgent

			conn.Write([]byte(content))

		case "files": // case "files" -> Server should handle files in regard to the specified method (GET or POST)

			// Base dir path where codecrafetrs's tester expects files to be
			dirPath := "/tmp/data/codecrafters.io/http-server-tester/"

			// File path composed of dir path specified above and path specified by request
			filePath := dirPath + endpoint[2]

			// Get file info (if any)
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				// If there's an "error" (file does not exist) check request's method

				if method == "POST" {
					// "POST" method: the server will try to create the file with specified content in it

					// Create the directory structure if it doesn't exist
					if err := os.MkdirAll(dirPath, 0755); err != nil {
						fmt.Println("Error creating directory:", err)
						conn.Write([]byte("HTTP/1.1 500 Internal Server Error"))
						os.Exit(1)
					}

					// Now, create the file
					_, err := os.Create(filePath)
					if err != nil {
						fmt.Println("Creation error:", err)
						conn.Write([]byte("HTTP/1.1 500 Internal Server Error"))
						os.Exit(1)
					}

					body := strings.Trim(lines[1], "\x00")

					err = os.WriteFile(filePath, []byte(body), 0664)
					if err != nil {
						fmt.Println("Write error:", err)
						conn.Write([]byte("HTTP/1.1 500 Internal Server Error"))
						os.Exit(1)
					}

					conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
				} else {
					conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
				}

			} else {
				// "GET" method: the server will try to read specified file's content

				fileSize := fileInfo.Size()

				// Read file's content and store it to use in response body
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println(string(fileContent))

				content := "HTTP/1.1 200 OK" +
					"\r\n" +
					"Content-Type: application/octet-stream" +
					"\r\n" +
					"Content-Length: " + fmt.Sprintf("%d\r\n", fileSize) +
					"\r\n" +
					string(fileContent)

				conn.Write([]byte(content))
			}

		default:
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
	}
}
