package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	var directory string
	fmt.Println("Logs from your program will appear here!")
	fmt.Println(os.Args[0])
	if len(os.Args) == 3 && os.Args[1] == "--directory" {
		directory = os.Args[2]
		fmt.Println(directory)
	}
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue 
		}

		go handleConnection(conn, directory) // Go Routine has been called here 
	}
}

func handleConnection(conn net.Conn, directory string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	var length int
	var userAgent string

	// Read the request line
	requestLine, err := reader.ReadString('\n')
	fmt.Printf("Request Line: %s", requestLine)
	if err != nil {
		fmt.Println("Error reading request line:", err.Error())
		return
	}

	// Parse the request line
	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) < 3 {
		fmt.Println("Malformed request line.")
		return
	}
	method, path, version := parts[0], parts[1], parts[2]
	fmt.Printf("Method: %s, Path: %s, Version: %s\n", method, path, version)

	if method == "POST" && strings.HasPrefix(path, "/files/") {
		fileName := path[7:]
		filePath := directory + "/" + fileName

		// Read the request body to get file contents
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Println("Error reading request body:", err.Error())
			return
		}

		// Write file contents to the specified directory
		err = ioutil.WriteFile(filePath, body, 0644)
		if err != nil {
			fmt.Println("Error writing file:", err.Error())
			return
		}

		// Respond with status code 201
		res := fmt.Sprintf("HTTP/1.1 201 Created\r\n\r\n%s", string(body))
		conn.Write([]byte(res))
		return
	}

	for {
		line, err := reader.ReadString('\n')
		fmt.Printf("Raw Header: %q\n", line)

		if err != nil {
			fmt.Println("Error reading headers:", err.Error())
			break
		}
		// Check if the line signifies the end of the headers
		if line == "r\n" {
			break
		}

		// Process each header line. Specifically, look for the User-Agent header.
		if strings.HasPrefix(line, "User-Agent:") {
			//fmt.Println("Hi there")
			userAgent = strings.TrimSpace(strings.TrimPrefix(line, "User-Agent:"))
			length = len(userAgent)
			fmt.Printf("User-Agent details: %s\n", userAgent)
			break // Assuming we're only looking for User-Agent, we can break after finding it
		}
		var res string // declared variable res of type string
		if strings.HasPrefix(path, "/files") {
			fileName := path[7:]
			//fmt.Println(fileName)
			filePath := directory + fileName
			//fmt.Println(filePath)
			if _, err := os.Open(directory); os.IsNotExist(err) {
				fmt.Println("Directory doesn't exists")
			} else {
				fmt.Println("The directory named", directory, "exists")
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					if os.IsNotExist(err) {
						fmt.Println("File does not exist.")
						res := "HTTP/1.1 404 Not Found\r\n\r\n"
						conn.Write([]byte(res))
					} else {
						fmt.Println("error reading file")
						log.Fatal(err)
					}
				} else {
					fileLength := len(fileContent)
					res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", fileLength, string(fileContent))
					conn.Write([]byte(res))
				}
			}
		} else if strings.HasPrefix(path, "/echo") {
			content := path[6:]
			fmt.Println(content)
			result := len(content)
			res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", result, content)
			conn.Write([]byte(res))
		} else if strings.HasPrefix(path, "/user-agent") {
			res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", length, userAgent)
			conn.Write([]byte(res))
		} else if path == "/" {
			res = "HTTP/1.1 200 OK\r\n\r\n"
			conn.Write([]byte(res))
		} else {
			res = "HTTP/1.1 404 Not Found\r\n\r\n"
			conn.Write([]byte(res))
		}
	}
}
