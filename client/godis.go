package client

type Key = string
type Argument = string

type Godis interface {
	Open() error
	Set(key Key, value string) error
	Get(key Key) (string, error)
	Exists(key Key, keys ...Key) (int64, error)
	Del(key Key, keys ...Key) (int64, error)
	FlushDb() error
	Execute(command RedisCommand) (*RespReply, error)
	Close()
}
