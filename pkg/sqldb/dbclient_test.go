package sqldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("")

	if err == nil {
		assert.NotNil(t, client)
	} else {
		assert.Error(t, err)
	}
}
