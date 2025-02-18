package challenger

import (
	"strconv"
	"time"
)

const (
	version = "1.0"
)

type HashcashChallenger struct {
	version    string
	difficulty int
}

func NewHashcashChallenger(cfg *Config) *HashcashChallenger {
	return &HashcashChallenger{
		version:    version,
		difficulty: cfg.Difficulty,
	}
}

func (c *HashcashChallenger) NewChallenge(resousceType string) Challenge {
	return &HashcashChallenge{
		version:      c.version,
		resourceType: resousceType,
		difficulty:   c.difficulty,
		timestamp:    strconv.FormatInt(time.Now().Unix(), 10),
	}
}

type EmptyChallenger struct {
}

func NewEmptyChallenger(_ string, _ int) *EmptyChallenger {
	return &EmptyChallenger{}
}

func (c *EmptyChallenger) NewChallenge() Challenge {
	return &EmptyChallenge{}
}
