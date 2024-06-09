package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	// data types
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// struct for commands and arguments for serialize/deserialize
// typical request/reponse will be an array of Value
type Value struct {
	type_of string
	bulk    string
	array   []Value
}

type Resp struct {
	reader *bufio.Reader
}

// creates a new Resp struct
func NewResp(reader io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(reader)}
}

// reads a line from the buffer, ok for utf8 strings since the purpose is to read till crlf
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			break
		}
	}

	return line[:len(line)-2], n, nil
}

// reads the integer from the buffer
func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	integer_64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}

	return int(integer_64), n, nil
}

// reads an array of Value, call to Read in the loop handles getting the values
func (r *Resp) readArray() (Value, error) {
	v := Value{type_of: "array"}

	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// read all the elements of the array
	v.array = make([]Value, length)
	for i := 0; i < length; i++ {
		value, err := r.Read()
		if err != nil {
			return v, err
		}

		v.array[i] = value
	}

	return v, nil
}

// reads a bulk string from the buffer
func (r *Resp) readBulk() (Value, error) {
	v := Value{type_of: "bulk"}

	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}
	if length <= 0 {
		return v, nil
	}

	bulk := make([]byte, length)

	// read the actual string
	r.reader.Read(bulk)
	v.bulk = string(bulk)

	// read the crlf
	r.readLine()

	return v, nil
}

// handles reading the different types of data from the buffer
func (r *Resp) Read() (Value, error) {
	type_of, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch type_of {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("ERROR: unknown type %v\n", string(type_of))
		return Value{}, nil
	}
}
