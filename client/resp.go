package client

import "fmt"
import "strings"

type Resp interface {
	Encode() RespEncoded
}

type RespEncoded interface {
	ToString() string
}

type RespBulkString struct {
	value string
}

func (bulkString *RespBulkString) Encode() RespEncoded {
	return &RespMultiEncoded{
		[]RespEncoded{
			&RespTypeSizeEncoded{"$", len(bulkString.value)},
			&RespSingleEncoded{bulkString.value},
		},
	}
}

type RespArray []Resp

func (a RespArray) Encode() RespEncoded {
	values := make([]RespEncoded, len(a))

	for i, v := range a {
		values[i] = v.Encode()
	}

	return &RespMultiEncoded{
		append(
			[]RespEncoded{&RespTypeSizeEncoded{"*", len(values)}},
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

type RespTypeSizeEncoded struct {
	respType string
	size     int
}

func (this *RespTypeSizeEncoded) ToString() string {
	simpleEncoding := RespSingleEncoded{
		fmt.Sprintf("%s%d", this.respType, this.size),
	}
	return simpleEncoding.ToString()
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
