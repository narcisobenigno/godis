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

func (g *GodisTcp) Execute(command RedisCommand) (*RespReply, error) {
	_, err := g.conn.Write([]byte(command.ToResp().Encode().ToString()))
	if err != nil {
		return nil, err
	}

	size := g.reader.Size()
	byteContent := make([]byte, size)
	g.reader.Read(byteContent)

	return &RespReply{string(byteContent)}, nil
}

func (g *GodisTcp) Close() {
	g.conn.Close()
}
