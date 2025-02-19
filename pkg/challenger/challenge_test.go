package challenger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHashcashChallenger(t *testing.T) {
	difficulty := 5
	challenger := NewHashcashChallenger(&Config{Difficulty: difficulty})

	assert.Equal(t, version, challenger.version)
	assert.Equal(t, difficulty, challenger.difficulty)
}

func TestNewChallenge(t *testing.T) {
	resourceType := "resource123"
	difficulty := 5
	challenger := NewHashcashChallenger(&Config{Difficulty: difficulty})
	chall, ok := challenger.NewChallenge(resourceType).(*HashcashChallenge)
	if !ok {
		t.Fatalf("expected *HashcashChallenge, got %T", chall)
	}

	assert.Equal(t, version, chall.version)
	assert.Equal(t, resourceType, chall.resourceType)
	assert.Equal(t, difficulty, chall.difficulty)
	assert.NotEmpty(t, chall.timestamp)
}

func TestGetChallengeMessage(t *testing.T) {
	resourceType := "resource123"
	difficulty := 1
	challenger := NewHashcashChallenger(&Config{Difficulty: difficulty})
	chall, ok := challenger.NewChallenge(resourceType).(*HashcashChallenge)
	if !ok {
		t.Fatalf("expected *HashcashChallenge, got %T", chall)
	}

	expectedMessage := "challenge: version=1.0, resourceType=resource123, timestamp=" + chall.timestamp + ", difficulty=1"
	assert.Equal(t, expectedMessage, chall.GetChallengeMessage())
}

func TestParseChallengeMessage(t *testing.T) {
	message := "challenge: version=1.0, resourceType=resource123, timestamp=1234567890, difficulty=5"
	chall := &HashcashChallenge{}

	err := chall.ParseChallengeMessage(message)
	require.NoError(t, err)
	assert.Equal(t, "1.0", chall.version)
	assert.Equal(t, "resource123", chall.resourceType)
	assert.Equal(t, "1234567890", chall.timestamp)
	assert.Equal(t, 5, chall.difficulty)
}

func TestVerifyPoW(t *testing.T) {
	resourceType := "resource123"
	difficulty := 1
	challenger := NewHashcashChallenger(&Config{Difficulty: difficulty})
	chall := challenger.NewChallenge(resourceType)

	nonce := chall.SolvePoW()
	assert.True(t, chall.VerifyPoW(nonce))
}

func TestSolvePoW(t *testing.T) {
	resourceType := "resource123"
	difficulty := 1
	challenger := NewHashcashChallenger(&Config{Difficulty: difficulty})
	chall := challenger.NewChallenge(resourceType)

	start := time.Now()
	nonce := chall.SolvePoW()
	duration := time.Since(start)

	assert.True(t, chall.VerifyPoW(nonce))
	t.Logf("Solved PoW in %s with nonce %s", duration, nonce)
}
