package client

import "fmt"

type RespType interface {
	Encode() RespEncoded
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

type RespArraySmart []RespType

func (a RespArraySmart) Encode() RespEncoded {
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
