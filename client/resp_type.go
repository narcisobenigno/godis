package client

import "fmt"
import "strings"

type RespType interface {
	Encode() RespEncoded
}

type RespEncoded interface {
	ToString() string
}

type RespBulkString struct {
	value string
}

func (bulkString *RespBulkString) Encode() RespEncoded {
	identifier := fmt.Sprintf("$%d", len(bulkString.value))

	return &RespMultiEncoded{
		[]RespEncoded{
			&RespSingleEncoded{identifier},
			&RespSingleEncoded{bulkString.value},
		},
	}
}

type RespArray []RespType

func (a RespArray) Encode() RespEncoded {
	values := make([]RespEncoded, len(a))

	for i, v := range a {
		values[i] = v.Encode()
	}

	identifier := fmt.Sprintf("*%d", len(values))
	return &RespMultiEncoded{
		append(
			[]RespEncoded{&RespSingleEncoded{identifier}},
			values...,
		),
	}
}

type RespSingleEncoded struct {
	value string
}

func (encoded *RespSingleEncoded) ToString() string {
	const RESP_TERMINATOR = "\r\n"
	return encoded.value + RESP_TERMINATOR
}

type RespMultiEncoded struct {
	values []RespEncoded
}

func (encoded *RespMultiEncoded) ToString() string {
	values := make([]string, len(encoded.values))

	for i, v := range encoded.values {
		values[i] = v.ToString()
	}
	return strings.Join(values, "")
}
