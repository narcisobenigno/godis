package client

import (
	"os"
	"testing"
)

const KEY = "TEST_KEY"
const VALUE = "TEST_VALUE"

var godis Godis

func setup() {
	godis = GodisNew("localhost:6379")
	godis.Open()
}

func TestThatSetSetsKeyWithValue(t *testing.T) {
	godis.Set(KEY, VALUE)
	if actual, _ := godis.Get(KEY); actual != VALUE {
		t.Errorf("Expected Get %s but got %s", VALUE, actual)
	}

	godis.Del(KEY)
}

func TestThatExistsReturnsZeroIfNotInsertedKey(t *testing.T) {
	if actual, _ := godis.Exists(KEY); actual != 0 {
		t.Errorf("Expected Exists 0 but got %d", actual)
	}
}

func teardown() {
	godis.FlushDb()
	godis.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
