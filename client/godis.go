package client

type Key = string
type Argument = string

type Godis interface {
	Open() error
	Execute(command RedisCommand) (*RespReply, error)
	Close()
}
