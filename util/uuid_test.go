package util

import (
	"testing"
)

func TestPing(t *testing.T) {
	testID := GetUUIDv4()
	if testID == "" || len(testID) < 10 {
		t.Fatal("Invalid ID generated")
	}
}
