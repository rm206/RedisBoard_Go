package main

import (
	"fmt"
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
		resp := NewResp(conn)

		// read from client
		_, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		// write to client using writer
		writer := NewRespWriter(conn)
		writer.Write(Value{type_of: "string", str: "OK"})
	}
}
