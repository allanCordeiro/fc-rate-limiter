package checker

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("RATE_LIMITER_LIMIT", "5")
	os.Setenv("RATE_LIMITER_IP_BLOCK_TIME", "60")

	os.Exit(m.Run())
}

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter()

	assert.Equal(t, 5, rl.ipLimit)
	assert.Equal(t, 60, rl.timeLimit)
}
