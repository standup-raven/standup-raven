package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasHave(t *testing.T) {
	assert.Equal(t, "have", HasHave(0), "0 is plural")
	assert.Equal(t, "has", HasHave(1), "1 is singular")
	assert.Equal(t, "have", HasHave(2), "2 is plural")
	assert.Equal(t, "has", HasHave(-1), "-1 is singular")
	assert.Equal(t, "have", HasHave(-2), "-2 is plural")
}

func TestSingularPlural(t *testing.T) {
	assert.Equal(t, "", SingularPlural(0), "0 is plural")
	assert.Equal(t, "", SingularPlural(1), "1 is singular")
	assert.Equal(t, "s", SingularPlural(2), "2 is plural")
	assert.Equal(t, "", SingularPlural(-1), "-1 is singular")
	assert.Equal(t, "s", SingularPlural(-2), "-2 is plural")
}
