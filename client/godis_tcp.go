package client

import "bufio"
import "net"

type GodisTcp struct {
	host   string
	conn   net.Conn
	reader *bufio.Reader
}

func GodisNew(host string) Godis {
	return &GodisTcp{host, nil, nil}
}

func (g *GodisTcp) Open() error {
	conn, err := net.Dial("tcp", g.host)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(conn)

	g.conn = conn
	g.reader = reader
	return nil
}

func (g *GodisTcp) Set(key Key, value string) error {
	command := &RedisCommandWithArgments{
		&SimpleRedisCommand{"SET"},
		[]Argument{key, value},
	}
	_, err := g.Execute(command)

	if err != nil {
		return err
	}

	return nil
}

func (g *GodisTcp) Get(key Key) (string, error) {
	command := &RedisCommandWithArgments{
		&SimpleRedisCommand{"GET"},
		[]Argument{key},
	}
	fullReturn, err := g.Execute(command)

	if err != nil {
		return "", err
	}

	reply := &RespBulkStringReply{fullReturn}

	return reply.String()
}

func (g *GodisTcp) Exists(key Key, keys ...Key) (int64, error) {
	respCommand := &RedisCommandWithArgments{
		&SimpleRedisCommand{"EXISTS"},
		append([]Argument{key}, keys...),
	}

	fullReturn, err := g.Execute(respCommand)
	if err != nil {
		return -1, err
	}
	reply := &RespIntergerReply{fullReturn}

	return reply.Integer()
}

func (g *GodisTcp) Del(key Key, keys ...Key) (int64, error) {
	respCommand := &RedisCommandWithArgments{
		&SimpleRedisCommand{"DEL"},
		append([]Argument{key}, keys...),
	}

	fullReturn, err := g.Execute(respCommand)
	if err != nil {
		return -1, err
	}

	reply := &RespIntergerReply{fullReturn}

	return reply.Integer()
}

func (g *GodisTcp) Execute(command RedisCommand) (string, error) {
	_, err := g.conn.Write([]byte(command.ToResp().Encode().ToString()))
	if err != nil {
		return "", err
	}

	size := g.reader.Size()
	byteContent := make([]byte, size)
	g.reader.Read(byteContent)

	return string(byteContent), nil
}

func (g *GodisTcp) FlushDb() error {
	if _, err := g.Execute(&SimpleRedisCommand{"FLUSHDB"}); err != nil {
		return err
	}
	return nil
}

func (g *GodisTcp) Close() {
	g.conn.Close()
}
