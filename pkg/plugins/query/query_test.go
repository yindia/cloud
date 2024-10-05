package query

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	q := &Query{}

	// Test case for successful execution
	t.Run("Successful execution", func(t *testing.T) {
		start := time.Now()
		err := q.Run(map[string]string{"success": "true"})
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if duration < 3*time.Second {
			t.Errorf("Expected execution time of at least 3 seconds, got %v", duration)
		}
	})

	// Test case for failure
	t.Run("Failure", func(t *testing.T) {
		maxRetries := 10
		failureOccurred := false

		for i := 0; i < maxRetries; i++ {
			err := q.Run(map[string]string{})
			if err != nil {
				failureOccurred = true
				if err.Error() != "query failed" {
					t.Errorf("Expected 'query failed', got %v", err)
				}
				break
			}
		}

		if !failureOccurred {
			t.Errorf("Expected at least one failure, but all attempts succeeded")
		}
	})
}
