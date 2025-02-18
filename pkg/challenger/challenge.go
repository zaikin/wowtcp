package challenger

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	challengeRegex = regexp.MustCompile(`challenge: version=(\S+), resourceType=(\S+), timestamp=(\S+), difficulty=(\d+)`)
)

const (
	callengeStrFormat = "challenge: version=%s, resourceType=%s, timestamp=%s, difficulty=%d"
)

//go:generate mockery --name=Challenge --output=./mocks --outpkg=mocks

type Challenge interface {
	GetChallengeMessage() string
	ParseChallengeMessage(message string) error
	VerifyPoW(nonce string) bool
	SolvePoW() string
}

type HashcashChallenge struct {
	version      string
	resourceType string
	timestamp    string
	difficulty   int
}

func (c *HashcashChallenge) GetChallengeMessage() string {
	return fmt.Sprintf(callengeStrFormat,
		c.version, c.resourceType, c.timestamp, c.difficulty)
}

func (c *HashcashChallenge) ParseChallengeMessage(message string) error {
	matches := challengeRegex.FindStringSubmatch(message)
	//nolint:mnd
	if len(matches) != 5 {
		err := errors.New("invalid challenge message format")
		log.Error().Err(err).Str("message", message).Msg("Error parsing challenge message")
		return err
	}

	log.Info().Interface("matches", matches).Msgf("Matches: %v", matches)

	c.version = matches[1]
	c.resourceType = matches[2]
	c.timestamp = matches[3]
	difficulty, err := strconv.Atoi(matches[4])
	if err != nil {
		log.Error().Err(err).Str("message", message).Msg("Error parsing difficulty")
		return err
	}
	c.difficulty = difficulty

	return nil
}

func (c *HashcashChallenge) VerifyPoW(nonce string) bool {
	input := fmt.Sprintf("%s:%s:%s:%s", c.version, c.resourceType, c.timestamp, nonce)
	hash := sha256.New()
	hash.Write([]byte(input))
	hashResult := hash.Sum(nil)
	hashString := hex.EncodeToString(hashResult)
	return strings.HasPrefix(hashString, strings.Repeat("0", c.difficulty))
}

func (c *HashcashChallenge) SolvePoW() string {
	for i := 0; ; i++ {
		nonce := strconv.FormatInt(int64(i), 10)
		if c.VerifyPoW(nonce) {
			return nonce
		}
	}
}

type EmptyChallenge struct {
}

func (c *EmptyChallenge) GetChallengeMessage() string {
	return "challenge"
}

func (c *EmptyChallenge) ParseChallengeMessage(_ string) error {
	return nil
}

func (c *EmptyChallenge) VerifyPoW(_ string) bool {
	return true
}

func (c *EmptyChallenge) SolvePoW() string {
	return "0"
}
