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

func (g *GodisTcp) Set(key, value string) error {
	resp := &RespArray{
		[]RespType{
			&RespBulkString{"SET"},
			&RespBulkString{key},
			&RespBulkString{value},
		},
	}
	command := resp.Encode().ToString()

	if _, err := g.conn.Write([]byte(command)); err != nil {
		return err
	}

	size := g.reader.Size()
	content := make([]byte, size)
	g.reader.Read(content)
	return nil
}

func (g *GodisTcp) Get(key string) (string, error) {
	resp := &RespArray{
		[]RespType{
			&RespBulkString{"GET"},
			&RespBulkString{key},
		},
	}
	command := resp.Encode().ToString()

	if _, err := g.conn.Write([]byte(command)); err != nil {
		return "", err
	} else {
		size := g.reader.Size()
		byteContent := make([]byte, size)
		g.reader.Read(byteContent)
		fullReturn := string(byteContent)

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
}

func (g *GodisTcp) Close() {
	g.conn.Close()
}
