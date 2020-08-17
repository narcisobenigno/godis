package client

type RedisCommand interface {
	ToResp() Resp
}

type SimpleRedisCommand struct {
	name string
}

func (command *SimpleRedisCommand) ToResp() Resp {
	return &RespBulkString{command.name}
}

type WithArgumentRedisCommand struct {
	command   RedisCommand
	arguments []Argument
}

func (command *WithArgumentRedisCommand) ToResp() Resp {
	keysAsBulkString := make([]Resp, len(command.arguments))
	for i, k := range command.arguments {
		keysAsBulkString[i] = &RespBulkString{k}
	}
	return RespArray(append(
		[]Resp{command.command.ToResp()},
		keysAsBulkString...,
	))
}
