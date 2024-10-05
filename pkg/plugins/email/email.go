package email

import (
	"os"
	"strconv"
	"time"
)

var PLUGIN_NAME = "send_email"

type Email struct {
}

func (e *Email) Run(parameters map[string]string) error {
	// Get timeout from TASK_TIME_OUT env variable or use 10 seconds as default
	timeout := 10
	if timeoutStr := os.Getenv("TASK_TIME_OUT"); timeoutStr != "" {
		if parsedTimeout, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = parsedTimeout
		}
	}
	time.Sleep(time.Duration(timeout) * time.Second)
	return nil
}
