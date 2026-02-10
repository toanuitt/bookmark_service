package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewRedisClient("")
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
