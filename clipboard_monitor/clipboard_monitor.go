package clipboard_monitor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/go-redis/redis/v8"
)

func getK() (string, string) {
	fullTime := strings.Split(time.Now().String(), " ")
	day := fullTime[0]
	hour := strings.Split(fullTime[1], ":")[0]

	return day, hour
}

func Monitor() {
	var lastText string
	var dayText string
	var hourText string

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Use default Addr
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})
	ctx := context.Background()

	for {
		currentText, err := clipboard.ReadAll()
		if err != nil {
			fmt.Printf("Error reading clipboard: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		if currentText != lastText {
			// fmt.Printf("Clipboard changed: %s\n", currentText)
			dayText, hourText = getK()
			redisClient.HSet(ctx, dayText, hourText, currentText)
			lastText = currentText
		}

		time.Sleep(time.Second) // Poll every 1 second
	}
}
