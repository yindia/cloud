package query

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var PLUGIN_NAME = "run_query"

type Query struct {
}

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// run_query executes a query and fails 20% of the time
func (q *Query) Run(parameters map[string]string) error {
	if seededRand.Float64() < 0.2 { // 20% chance to fail
		return fmt.Errorf("query failed")
	}

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
