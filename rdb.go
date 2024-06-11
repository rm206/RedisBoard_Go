package main

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type RDB struct {
	file  *os.File
	rd    *bufio.Reader
	mutex sync.Mutex
}

// will open a file and clear it if it is older than "erase_after_days" days
func NewRdb(path string, erase_after_days int, flushMapAfterSeconds int, syncAfterSeconds int) (*RDB, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	cleared, _ := clearFileIfOld(f, erase_after_days)
	if cleared {
		f.WriteString(time.Now().String() + "\n")
	}

	rdb := &RDB{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// Start a goroutine to flush map to disk
	go func() {
		for {
			rdb.mutex.Lock()

			rdb.Write()

			rdb.mutex.Unlock()

			time.Sleep(time.Duration(flushMapAfterSeconds) * time.Second) // flush every flushMapAfterSeconds
		}
	}()

	// Start a goroutine to sync RDB to disk
	go func() {
		for {
			rdb.mutex.Lock()

			rdb.file.Sync()

			rdb.mutex.Unlock()

			time.Sleep(time.Duration(syncAfterSeconds) * time.Second) // sync every syncAfterSeconds
		}
	}()

	return rdb, nil
}

// close file when server shuts down
func (rdb *RDB) Close() error {
	rdb.mutex.Lock()
	defer rdb.mutex.Unlock()

	return rdb.file.Close()
}

func (rdb *RDB) Write() error {
	rdb.mutex.Lock()
	defer rdb.mutex.Unlock()

	// write to file
	_, err := rdb.file.WriteString(getHSETs_mapString())
	if err != nil {
		return err
	}

	// clear HSET_map
	clearHSETs_map()

	return nil
}
