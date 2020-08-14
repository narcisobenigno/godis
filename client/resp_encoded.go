package client

import "strings"

type RespEncoded interface {
	ToString() string
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
