package main

import (
	"bufio"
	"os"
	"time"
)

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
