package client

import (
	"testing"
)

func TestAllowClientSetAndGet(t *testing.T) {
	client := GodisNew("localhost:6379")
	client.Open()

	client.Set("TEST_KEY", "TEST_VALUE")
	if actual, _ := client.Get("TEST_KEY"); actual != "TEST_VALUE" {
		t.Errorf("Expected Get %s but got %s", "TEST_VALUE", actual)
	}

	defer client.Close()
}
