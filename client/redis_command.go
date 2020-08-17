package client

type RedisCommand interface {
	ToResp() RespType
}

type SimpleRedisCommand struct {
	name string
}

func (command *SimpleRedisCommand) ToResp() RespType {
	return &RespBulkString{command.name}
}

type WithArgumentRedisCommand struct {
	command   RedisCommand
	arguments []Argument
}

func (command *WithArgumentRedisCommand) ToResp() RespType {
	keysAsBulkString := make([]RespType, len(command.arguments))
	for i, k := range command.arguments {
		keysAsBulkString[i] = &RespBulkString{k}
	}
	return RespArray(append(
		[]RespType{command.command.ToResp()},
		keysAsBulkString...,
	))
}
