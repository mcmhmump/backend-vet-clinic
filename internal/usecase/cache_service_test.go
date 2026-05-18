package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheService_SetGet(t *testing.T) {
	cache := NewCacheService()
	key := "test-key"
	value := []byte("hello")

	cache.Set(key, value, 5*time.Second)

	result, found := cache.Get(key)

	assert.True(t, found)
	assert.Equal(t, value, result)
}

func TestCacheService_Expire(t *testing.T) {
	cache := NewCacheService()
	key := "test-key"
	value := []byte("hello")

	cache.Set(key, value, 500*time.Millisecond)
	time.Sleep(600 * time.Millisecond)

	_, found := cache.Get(key)

	assert.False(t, found)
}

func TestCacheService_Invalidate(t *testing.T) {
	cache := NewCacheService()
	key := "test-key"
	value := []byte("hello")

	cache.Set(key, value, 5*time.Second)
	cache.Invalidate(key)

	_, found := cache.Get(key)

	assert.False(t, found)
}
