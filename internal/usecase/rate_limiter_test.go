package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(2, 2*time.Second)
	ip := "127.0.0.1"

	assert.True(t, rl.Allow(ip))
	assert.True(t, rl.Allow(ip))
	assert.False(t, rl.Allow(ip))
}

func TestRateLimiter_ResetAfterWindow(t *testing.T) {
	rl := NewRateLimiter(1, 1*time.Second)
	ip := "127.0.0.1"

	assert.True(t, rl.Allow(ip))
	assert.False(t, rl.Allow(ip))

	time.Sleep(1100 * time.Millisecond)

	assert.True(t, rl.Allow(ip))
}
