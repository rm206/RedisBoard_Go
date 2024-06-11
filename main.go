package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	erase_after_days := 5
	flushMapAfterSeconds := 10 * 60
	syncAfterSeconds := 120

	if len(os.Args) > 1 {
		temp := os.Args[1]
		temp_converted, err := strconv.Atoi(temp)
		if err != nil {
			fmt.Println("Invalid argument, defaulting to erasing dump in 5 days")
		} else {
			erase_after_days = temp_converted
		}
	} else {
		fmt.Println("No argument provided, defaulting to erasing dump in 5 days")
	}

	fmt.Println("Listening on port 6379")

	// Listen on port 6379
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dump_path_rdb := "dump.rdb"
	rdb, err := NewRdb(dump_path_rdb, erase_after_days, flushMapAfterSeconds, syncAfterSeconds)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rdb.Close()

	dump_path_aof := "dump.aof"
	aof, err := NewAof(dump_path_aof, erase_after_days, syncAfterSeconds)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer aof.Close()

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
		val, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		if val.type_of != "array" {
			fmt.Println("Expected array")
			continue
		}
		if len(val.array) < 1 {
			fmt.Println("Expected at least one argument")
			continue
		}

		// get the command
		cmd := strings.ToUpper(val.array[0].bulk)
		args := val.array[1:]

		// write to client using writer
		writer := NewRespWriter(conn)

		handler, ok := Handlers[cmd]
		if !ok {
			fmt.Println("Command not found")
			writer.Write(Value{type_of: "string", str: "ERR unknown command '" + cmd + "'"})
			continue
		}

		if cmd == "SET" || cmd == "HSET" {
			aof.Write(val)
		}

		// call the handler
		res := handler(args)

		writer.Write(res)
	}
}
