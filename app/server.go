package main

import (
	"bufio"
	"fmt"
	"strings"

	//Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221") //in net package net.Listen listens for the tcp and reserves the port 4221 from anywhere(0.0.0.0)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept() //it's a blocking call mean it waits till client connects
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	var data []byte

	_, err = conn.Read(data) //client connects either HTTP or HTTPS so .Read method is to read  request from a client and it only accepts []byte slice

	if err != nil {
		fmt.Println("Error Reading connection:", err.Error())
		os.Exit(1)
	}
	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n")) // .Write method is write the response to the client in this case we are sending  2200 OK and HTTP/1.1 is the protocol version

	if err != nil {
		fmt.Println("Error Writing response:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(conn) // buifo.reader is for creating a buffered reader

	headers := make(map[string]string)  //created a map that accepts string
	lines, _ := reader.ReadString('\n') //.ReadString reads from the input until the first occurence of delimter in this case 1st line
	parts := strings.Split(lines, " ")  //Split method is used to seperate  by seperator and it returns a string slice
	headers["action"] = parts[0]        //stored 0th element of slice in the headers map as action for example get post put etc
	headers["route"] = parts[1]         //same way 1st element of slice in map by mapping as route or path
	headers["version"] = parts[2]       //same way  2nd element of slice in map by mapping as version

	// 	GET /index.html HTTP/1.1  1st line
	// Host: localhost:4221  --> remaining part of headers
	// User-Agent: curl/7.64.1  --> remaining part of headers

	for {
		line, err := reader.ReadString('\n')
		if err != nil || line != "\r\n" {
			parts = strings.Split(line, " ") //stored remaining part of the headers in the map
			headers[parts[0]] = parts[1]
			fmt.Printf("%v %v\n", parts[0], headers[parts[0]])
		} else {
			break
		}
	}
	if headers["route"] == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}
