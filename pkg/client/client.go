package client

import (
	"fmt"
	"io"
	"net"
	"strings"
	"wowtcp/pkg/challenger"
	"wowtcp/pkg/tcpio"

	"github.com/pkg/errors"
)

var (
	ErrConncections     = errors.New("error connecting to server: ")
	ErrSendingMessage   = errors.New("error sending message: ")
	ErrReadingMessage   = errors.New("error reading message: ")
	ErrParsingMessage   = errors.New("error parsing message: ")
	ErrConncectionClose = errors.New("error closing connection")
)

const (
	QuoteCommand    = "quote!"
	QuitCommand     = "quit!"
	NoncePrefix     = "nonce: "
	QuotePrefix     = "quote: "
	ChallengePrefix = "challenge: "
)

type MessageType int

const (
	EmptyMessage MessageType = iota + 1
	ChallengeMessage
	QuoteMessage
	QuitMessage
)

type Client struct {
	host     string
	port     string
	conn     net.Conn
	messages tcpio.ReadWriter
}

func NewClient(host, port string) *Client {
	return &Client{
		host: host,
		port: port,
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.host, c.port))
	if err != nil {
		err = errors.Wrap(err, ErrConncections.Error())
		return err
	}
	c.conn = conn
	c.messages = tcpio.NewTCPReadWriter(conn)
	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) RequestQuote() (string, error) {
	if _, err := c.messages.Write(QuoteCommand); err != nil {
		err = errors.Wrap(err, ErrSendingMessage.Error())
		return "", err
	}
	var quote string
	for {
		message, err := c.messages.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return "", ErrConncectionClose
			}
			err = errors.Wrap(err, ErrReadingMessage.Error())
			return "", err
		}
		var messageType MessageType
		messageType, quote, err = c.handleMessage(message)
		if err != nil {
			return "", err
		}
		if messageType == QuoteMessage {
			break
		}
		if messageType == QuitMessage {
			return "", ErrConncectionClose
		}
	}
	return quote, nil
}

func (c *Client) handleMessage(message string) (MessageType, string, error) {
	switch {
	case message == QuitCommand:
		c.Close()
		return QuitMessage, "", nil
	case strings.HasPrefix(message, ChallengePrefix):
		return ChallengeMessage, "", c.handleChallenge(message)
	case strings.HasPrefix(message, QuotePrefix):
		return QuoteMessage, strings.TrimPrefix(message, QuotePrefix), nil
	default:
	}
	return EmptyMessage, "", nil
}

func (c *Client) handleChallenge(message string) error {
	chall := challenger.HashcashChallenge{}

	if err := chall.ParseChallengeMessage(message); err != nil {
		err = errors.Wrap(err, ErrParsingMessage.Error())
		return err
	}

	nonce := chall.SolvePoW()

	if _, err := c.messages.Write(NoncePrefix + nonce); err != nil {
		err = errors.Wrap(err, "error on sending nonce")
		err = errors.Wrap(err, ErrSendingMessage.Error())
		return err
	}
	return nil
}
