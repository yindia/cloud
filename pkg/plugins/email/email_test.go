package email

import (
	"testing"
)

func TestEmail_Run(t *testing.T) {
	e := &Email{}

	testCases := []struct {
		name       string
		parameters map[string]string
	}{
		{"Valid parameters", map[string]string{"key": "value"}},
		{"Empty parameters", map[string]string{}},
		{"Multiple parameters", map[string]string{"key1": "value1", "key2": "value2"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			err := e.Run(tc.parameters)

			// Test for no error
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

		})
	}
}
