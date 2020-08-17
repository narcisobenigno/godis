package client

type RedisCommand interface {
	ToResp() Resp
}

type SimpleRedisCommand struct {
	name string
}

func (this *SimpleRedisCommand) ToResp() Resp {
	return &RespBulkString{this.name}
}

type RedisCommandWithArgments struct {
	command   RedisCommand
	arguments []Argument
}

func (this *RedisCommandWithArgments) ToResp() Resp {
	keysAsBulkString := make([]Resp, len(this.arguments))
	for i, k := range this.arguments {
		keysAsBulkString[i] = &RespBulkString{k}
	}
	return RespArray(append(
		[]Resp{this.command.ToResp()},
		keysAsBulkString...,
	))
}
