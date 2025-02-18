package tcpserver

import (
	"bytes"
	"context"
	"net"
	"testing"
	"wowtcp/internal/tcpserver/mocks"
	challengerMocks "wowtcp/pkg/challenger/mocks"

	"github.com/stretchr/testify/assert"
)

type mockConn struct {
	net.Conn
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuffer.Read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuffer.Write(b)
}

func TestHandleQuote(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetWoWQuote").Return("This is a WoW quote")

	mockChallenger := new(mocks.Challenger)
	mockChallenge := new(challengerMocks.Challenge)
	mockChallenge.On("GetChallengeMessage").Return("challenge: version=1.0, resourceType=quote, timestamp=1234567890, difficulty=2")
	mockChallenge.On("VerifyPoW", "000000").Return(true)
	mockChallenger.On("NewChallenge", "quote").Return(mockChallenge)

	server := &Server{
		repository: mockRepo,
		challenger: mockChallenger,
	}

	ctx := context.Background()
	readBuffer := bytes.NewBufferString("000000\n")
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}

	err := server.handleQuote(ctx, conn)
	assert.NoError(t, err)

	expectedChallenge := "challenge: version=1.0, resourceType=quote, timestamp=1234567890, difficulty=2\n"
	expectedQuote := "quote: This is a WoW quote\n"

	assert.Contains(t, writeBuffer.String(), expectedChallenge)
	assert.Contains(t, writeBuffer.String(), expectedQuote)
}

func TestHandleQuoteInvalidNonce(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockChallenger := new(mocks.Challenger)
	mockChallenge := new(challengerMocks.Challenge)
	mockChallenge.On("GetChallengeMessage").Return("challenge: version=1.0, resourceType=quote, timestamp=1234567890, difficulty=2")
	mockChallenge.On("VerifyPoW", "invalid").Return(false)
	mockChallenger.On("NewChallenge", "quote").Return(mockChallenge)

	server := &Server{
		repository: mockRepo,
		challenger: mockChallenger,
	}

	ctx := context.Background()
	readBuffer := bytes.NewBufferString("invalid\n")
	writeBuffer := new(bytes.Buffer)
	conn := &mockConn{
		readBuffer:  readBuffer,
		writeBuffer: writeBuffer,
	}

	err := server.handleQuote(ctx, conn)
	assert.NoError(t, err)

	expectedChallenge := "challenge: version=1.0, resourceType=quote, timestamp=1234567890, difficulty=2\n"
	expectedError := "Invalid nonce\n"

	assert.Contains(t, writeBuffer.String(), expectedChallenge)
	assert.Contains(t, writeBuffer.String(), expectedError)
}
