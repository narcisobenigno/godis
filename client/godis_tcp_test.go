package client

import (
	"os"
	"testing"
)

const KEY = "KEY"
const KEY2 = "KEY2"
const VALUE = "VALUE"

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

func TestThatExistsReturns10When10KeysExists(t *testing.T) {
	godis.Set(KEY, VALUE)
	godis.Set(KEY2, VALUE)
	if actual, _ := godis.Exists(KEY, KEY2); actual != 2 {
		t.Errorf("Expected Exists 2 but got %d", actual)
	}
	godis.Del(KEY)
	godis.Del(KEY2)
}

func TestThatDeletes2Keys(t *testing.T) {
	godis.Set(KEY, VALUE)
	godis.Set(KEY2, VALUE)
	godis.Del(KEY, KEY2)
	if actual, _ := godis.Exists(KEY, KEY2); actual != 0 {
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
