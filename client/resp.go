package client

import "fmt"
import "strconv"
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

func (this RespArray) Encode() RespEncoded {
	values := make([]RespEncoded, len(this))

	for i, v := range this {
		values[i] = v.Encode()
	}

	return &RespMultiEncoded{
		append(
			[]RespEncoded{&RespTypeSizeEncoded{"*", len(this)}},
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

type RespReply struct {
	text string
}

func (this *RespReply) Integer() (int64, error) {
	contentStartingAt := int64(1)
	for i := 1; this.text[i] != '\r'; i++ {
		contentStartingAt++
	}
	return strconv.ParseInt(this.text[1:contentStartingAt], 10, 64)
}

func (this *RespReply) String() (string, error) {
	contentStartingAt := int64(1)
	for i := 1; this.text[i] != '\r'; i++ {
		contentStartingAt++
	}
	keySize, _ := strconv.ParseInt(this.text[1:contentStartingAt], 10, 64)
	contentStartingAt++
	contentStartingAt++
	contentFinishAt := contentStartingAt + keySize

	return this.text[contentStartingAt:contentFinishAt], nil
}
