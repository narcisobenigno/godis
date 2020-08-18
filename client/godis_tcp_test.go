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
	(&SetCommand{godis, KEY, VALUE}).Execute()
	actual, _ := (&GetCommand{godis, KEY}).Execute()
	if actual != VALUE {
		t.Errorf("Expected Get %s but got %s", VALUE, actual)
	}

	(&DelCommand{godis, KEY, []Key{}}).Execute()
}

func TestThatExistsReturnsZeroIfNotInsertedKey(t *testing.T) {
	actual, _ := (&ExistsCommand{godis, KEY, []Key{}}).Execute()
	if actual != 0 {
		t.Errorf("Expected Exists 0 but got %d", actual)
	}
}

func TestThatExistsReturns10When10KeysExists(t *testing.T) {
	(&SetCommand{godis, KEY, VALUE}).Execute()
	(&SetCommand{godis, KEY2, VALUE}).Execute()
	actual, _ := (&ExistsCommand{godis, KEY, []Key{KEY2}}).Execute()
	if actual != 2 {
		t.Errorf("Expected Exists 2 but got %d", actual)
	}
	(&DelCommand{godis, KEY, []Key{}}).Execute()
	(&DelCommand{godis, KEY2, []Key{}}).Execute()
}

func TestThatDeletes2Keys(t *testing.T) {
	(&SetCommand{godis, KEY, VALUE}).Execute()
	(&SetCommand{godis, KEY2, VALUE}).Execute()
	(&DelCommand{godis, KEY, []Key{KEY2}}).Execute()
	actual, _ := (&ExistsCommand{godis, KEY, []Key{KEY2}}).Execute()
	if actual != 0 {
		t.Errorf("Expected Exists 0 but got %d", actual)
	}
}

func teardown() {
	(&FlushDbCommand{godis}).Execute()
	godis.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
