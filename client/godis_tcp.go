package client

import "bufio"
import "net"
import "strconv"

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
	command := &WithArgumentRedisCommand{
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
	command := &WithArgumentRedisCommand{
		&SimpleRedisCommand{"GET"},
		[]Argument{key},
	}
	fullReturn, err := g.Execute(command)

	if err != nil {
		return "", err
	}

	contentStartingAt := int64(1)
	for i := 1; fullReturn[i] != '\r'; i++ {
		contentStartingAt++
	}
	keySize, _ := strconv.ParseInt(fullReturn[1:contentStartingAt], 10, 64)
	contentStartingAt++
	contentStartingAt++
	contentFinishAt := contentStartingAt + keySize

	return fullReturn[contentStartingAt:contentFinishAt], nil
}

func (g *GodisTcp) Exists(key Key, keys ...Key) (int64, error) {
	respCommand := &WithArgumentRedisCommand{
		&SimpleRedisCommand{"EXISTS"},
		append([]Argument{key}, keys...),
	}

	fullReturn, err := g.Execute(respCommand)
	if err != nil {
		return -1, err
	}

	contentStartingAt := int64(1)
	for i := 1; fullReturn[i] != '\r'; i++ {
		contentStartingAt++
	}
	existing, _ := strconv.ParseInt(fullReturn[1:contentStartingAt], 10, 64)

	return existing, nil
}

func (g *GodisTcp) Del(key Key, keys ...Key) (int64, error) {
	respCommand := &WithArgumentRedisCommand{
		&SimpleRedisCommand{"DEL"},
		append([]Argument{key}, keys...),
	}

	fullReturn, err := g.Execute(respCommand)
	if err != nil {
		return -1, err
	}

	contentStartingAt := int64(1)
	for i := 1; fullReturn[i] != '\r'; i++ {
		contentStartingAt++
	}
	existing, _ := strconv.ParseInt(fullReturn[1:contentStartingAt], 10, 64)

	return existing, nil
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
