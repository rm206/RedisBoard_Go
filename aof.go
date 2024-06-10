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

// clear the file if it is older than "erase_after_days" days
func clearFileIfOld(f *os.File, erase_after_days int) (bool, error) {
	// Read the first line of the file
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		// File is empty, no need to clear
		return false, nil
	}
	firstLine := scanner.Text()

	time_str_format := "2006-01-02 15:04:05.999999 -0700 MST m=+0.000000000"
	old_time, _ := time.Parse(time_str_format, firstLine)
	date_days_ago := time.Now().AddDate(0, 0, -erase_after_days)

	if old_time.Before(date_days_ago) {
		// Clear the file
		f.Truncate(0)
		f.Seek(0, 0)
		return true, nil
	}

	return false, nil
}

func NewAof(path string, erase_after_days int) (*AOF, error) {
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

			time.Sleep(120 * time.Second) // sync every 2 minutes
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
