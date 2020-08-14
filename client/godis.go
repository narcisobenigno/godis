package client

type Godis interface {
	Open() error
	Close()
	Set(key, value string) error
	Get(key string) (string, error)
}
