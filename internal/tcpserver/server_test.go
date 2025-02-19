package tcpserver

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"
	"wowtcp/internal/tcpserver/mocks"
	challengerMocks "wowtcp/pkg/challenger/mocks"
	"wowtcp/pkg/tcpio"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *mockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
	}
}

func (m *mockConn) Close() error {
	return nil
}

func TestHandleConnectionQuit(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockChallenger := new(mocks.Challenger)

	server := &Server{
		repository: mockRepo,
		challenger: mockChallenger,
	}

	ctx := context.Background()
	readBuffer := bytes.NewBufferString("quit!\n")
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}

	server.handleConnection(ctx, conn)

	assert.Contains(t, writeBuffer.String(), "")
}

func TestHandleConnectionQuote(t *testing.T) {
	challengeStr := "challenge: version=1.0, resourceType=quote, timestamp=1234567890, difficulty=2"
	quote := "This is a WoW quote"
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetWoWQuote").Return(quote)

	mockChallenger := new(mocks.Challenger)
	mockChallenge := new(challengerMocks.Challenge)
	mockChallenge.On("GetChallengeMessage").Return(challengeStr)
	mockChallenge.On("VerifyPoW", mock.Anything).Return(true)
	mockChallenger.On("NewChallenge", mock.Anything).Return(mockChallenge)

	server := &Server{
		repository: mockRepo,
		challenger: mockChallenger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Hack for testing. We need to cancel the context after a delay because the handleConnection function have infinite loop
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	readBuffer := bytes.NewBufferString("quote!\nnonce: 000000\n")
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}

	server.handleConnection(ctx, conn)

	assert.Contains(t, writeBuffer.String(), challengeStr)
	assert.Contains(t, writeBuffer.String(), quote)
}

func TestHandleConnectionUnknownMessage(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockChallenger := new(mocks.Challenger)

	server := &Server{
		repository: mockRepo,
		challenger: mockChallenger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Hack for testing. We need to cancel the context after a delay because the handleConnection function have infinite loop
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	readBuffer := bytes.NewBufferString("unknown\n")
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}

	server.handleConnection(ctx, conn)

	assert.Contains(t, writeBuffer.String(), "")
}

func TestHandleConnectionErrorReadingMessage(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockChallenger := new(mocks.Challenger)

	server := &Server{
		repository: mockRepo,
		challenger: mockChallenger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Hack for testing. We need to cancel the context after a delay because the handleConnection function have infinite loop
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	readBuffer := new(bytes.Buffer)
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}

	server.handleConnection(ctx, conn)

	assert.Contains(t, writeBuffer.String(), "")
}

func TestHandleQuote(t *testing.T) {
	expectedChallenge := "challenge: version=1.0, resourceType=quote, timestamp=1234567890, difficulty=2"
	expectedQuote := "quote: This is a WoW quote"

	mockRepo := new(mocks.Repository)
	mockRepo.On("GetWoWQuote").Return(expectedQuote)

	mockChallenge := new(challengerMocks.Challenge)
	mockChallenge.On("GetChallengeMessage").Return(expectedChallenge)
	mockChallenge.On("VerifyPoW", mock.Anything).Return(true)

	mockChallenger := new(mocks.Challenger)
	mockChallenger.On("NewChallenge", mock.Anything).Return(mockChallenge)

	server := &Server{
		repository: mockRepo,
		challenger: mockChallenger,
	}

	ctx := context.Background()
	readBuffer := bytes.NewBufferString("nonce: 000000\n")
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}
	messages := tcpio.NewTCPReadWriter(conn)

	err := server.handleQuote(ctx, messages)
	require.NoError(t, err)

	assert.Contains(t, writeBuffer.String(), expectedChallenge)
	assert.Contains(t, writeBuffer.String(), expectedQuote)
}

func TestHandleQuoteInvalidNonce(t *testing.T) {
	expectedChallenge := "challenge: version=1.0, resourceType=quote, timestamp=1234567890, difficulty=2\n"
	expectedError := "Invalid nonce\n"

	mockRepo := new(mocks.Repository)

	mockChallenger := new(mocks.Challenger)
	mockChallenge := new(challengerMocks.Challenge)
	mockChallenge.On("GetChallengeMessage").Return(expectedChallenge)
	mockChallenge.On("VerifyPoW", "invalid").Return(false)
	mockChallenger.On("NewChallenge", "quote").Return(mockChallenge)

	server := &Server{
		repository: mockRepo,
		challenger: mockChallenger,
	}

	ctx := context.Background()
	readBuffer := bytes.NewBufferString("nonce: invalid\n")
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}
	messages := tcpio.NewTCPReadWriter(conn)

	err := server.handleQuote(ctx, messages)
	require.NoError(t, err)

	assert.Contains(t, writeBuffer.String(), expectedChallenge)
	assert.Contains(t, writeBuffer.String(), expectedError)
}
