package util

import (
	"testing"
)

import "github.com/stretchr/testify/assert"

func TestHasHave(t *testing.T) {
	assert.Equal(t, "have", HasHave(0), "0 is plural")
	assert.Equal(t, "has", HasHave(1), "1 is singular")
	assert.Equal(t, "have", HasHave(2), "2 is plural")
	assert.Equal(t, "has", HasHave(-1), "-1 is singular")
	assert.Equal(t, "have", HasHave(-2), "-2 is plural")
}
