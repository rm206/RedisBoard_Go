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

// even though this is supposed to be an append only file, we will clear the file if it is older than "erase_after_days" days. Slighy liberty for use case
func NewAof(path string, erase_after_days int, syncAfterSeconds int) (*AOF, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	cleared, _ := clearFileIfOld(f, erase_after_days)
	if cleared {
		f.WriteString(time.Now().String() + "\n")
	}

	aof := &AOF{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk
	go func() {
		for {
			aof.mutex.Lock()

			aof.file.Sync()

			aof.mutex.Unlock()

			time.Sleep(time.Duration(syncAfterSeconds) * time.Second) // sync every 2 minutes
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
