package tcpio

import (
	"bytes"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConn struct {
	net.Conn
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
}

func (m *mockConn) Read(b []byte) (int, error) {
	return m.readBuffer.Read(b)
}

func (m *mockConn) Write(b []byte) (int, error) {
	return m.writeBuffer.Write(b)
}

func (m *mockConn) Close() error {
	return nil
}

func TestTcpReadWriter_Read(t *testing.T) {
	expectedMessage := "Hello, World!"
	readBuffer := bytes.NewBufferString(expectedMessage + "\n")
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}
	rw := NewTCPReadWriter(conn)

	message, err := rw.Read()
	require.NoError(t, err)
	assert.Equal(t, expectedMessage, message)
}

func TestTcpReadWriter_ReadManyLines(t *testing.T) {
	expectedMessage := "Hello, World!"
	expectedMessage2 := "Hello, World, Again!"
	readBuffer := bytes.NewBufferString(expectedMessage + "\n" + expectedMessage2 + "\n")
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}
	rw := NewTCPReadWriter(conn)

	message, err := rw.Read()
	require.NoError(t, err)
	assert.Equal(t, expectedMessage, message)
	message, err = rw.Read()
	require.NoError(t, err)
	assert.Equal(t, expectedMessage2, message)
}

func TestTcpReadWriter_Write(t *testing.T) {
	message := "Hello, World!"
	readBuffer := new(bytes.Buffer)
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}
	rw := NewTCPReadWriter(conn)

	n, err := rw.Write(message + "\n")
	require.NoError(t, err)
	assert.Equal(t, len(message)+1, n)
	assert.Equal(t, message+"\n", writeBuffer.String())
}

func TestTcpReadWriter_ReadWrite(t *testing.T) {
	message := "Hello, World!"
	readBuffer := bytes.NewBufferString(message + "\n")
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}
	rw := NewTCPReadWriter(conn)

	readMessage, err := rw.Read()
	require.NoError(t, err)
	assert.Equal(t, message, readMessage)

	n, err := rw.Write(readMessage + "\n")
	require.NoError(t, err)
	assert.Equal(t, len(readMessage)+1, n)
	assert.Equal(t, message+"\n", writeBuffer.String())
}
