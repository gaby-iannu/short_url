package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestNew(t *testing.T) {
	url := New("abcdefgh.com","user")
	assert.NotNil(t, url)
	assert.Equal(t, "user", url.User)
	assert.Equal(t, "abcdefgh.com", url.Long)
}
