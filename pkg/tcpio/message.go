package tcpio

import (
	"bufio"
	"bytes"
	"io"
	"net"
)

//go:generate mockery --name=Messages --output=./mocks --outpkg=mocks

type ReadWriter interface {
	Read() (string, error)
	Write(content string) (int, error)
}

type TCPReadWriter struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewTCPReadWriter(conn net.Conn) ReadWriter {
	return &TCPReadWriter{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (t *TCPReadWriter) Read() (string, error) {
	var buffer bytes.Buffer
	for {
		ba, isPrefix, err := t.reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		buffer.Write(ba)
		if !isPrefix {
			break
		}
	}
	return buffer.String(), nil
}

func (t *TCPReadWriter) Write(content string) (int, error) {
	number, err := t.writer.WriteString(content)
	if err == nil {
		err = t.writer.Flush()
	}
	return number, err
}
