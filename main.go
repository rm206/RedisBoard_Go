package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Listening on port 6379")

	// Listen on port 6379
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// close connection when done
	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		// read from client
		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}

		// write to client
		conn.Write([]byte("+OK\r\n"))
	}

}
