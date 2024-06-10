package main

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type AOF struct {
	file  *os.File
	rd    *bufio.Reader
	mutex sync.Mutex
}

func NewAof(path string) (*AOF, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &AOF{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.mutex.Lock()

			aof.file.Sync()

			aof.mutex.Unlock()

			time.Sleep(time.Second) // sync every 2 minutes
		}
	}()

	return aof, nil
}

// close file when server shuts down
func (aof *AOF) Close() error {
	aof.mutex.Lock()
	defer aof.mutex.Unlock()

	return aof.file.Close()
}

func (aof *AOF) Write(value Value) error {
	aof.mutex.Lock()
	defer aof.mutex.Unlock()

	_, err := aof.file.Write(value.Serialize())
	if err != nil {
		return err
	}

	return nil
}
