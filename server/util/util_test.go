package util

import (
	"testing"
)

import "github.com/stretchr/testify/assert"

func TestSplitArgs(t *testing.T) {
	actual, err := SplitArgs("foo bar baz")
	assert.Nil(t, err)
	assert.Equal(t, actual, []string{"foo", "bar", "baz"}, "space delimited params are extracted")

	actual, err = SplitArgs("foo \"bar baz\"")
	assert.Nil(t, err)
	assert.Equal(t, actual, []string{"foo", "bar baz"}, "quoted parameters are not splitted")

	actual, err = SplitArgs("")
	assert.Nil(t, err)
	assert.Equal(t, actual, []string{}, "0 params are fine")

	actual, err = SplitArgs("foo")
	assert.Nil(t, err)
	assert.Equal(t, actual, []string{"foo"}, "single param is fine")

	actual, err = SplitArgs("foo \"bar baz")
	assert.NotNil(t, err, "error should be produced ad quote is not closed")

	actual, err = SplitArgs("   foo   ")
	assert.Nil(t, err)
	assert.Equal(t, actual, []string{"foo"}, "leading and trailing spaces are stripped")

	actual, err = SplitArgs("foo     bar")
	assert.Nil(t, err)
	assert.Equal(t, actual, []string{"foo", "bar"}, "intermediate spaces are stripped")
}

func TestMin(t *testing.T) {
	assert.Equal(t, 1, Min(1, 2))
	assert.Equal(t, 1, Min(2, 1))
	assert.Equal(t, -1, Min(-1, 0))
	assert.Equal(t, -2, Min(-1, -2))
}

func TestMax(t *testing.T) {
	assert.Equal(t, 2, Max(1, 2))
	assert.Equal(t, 2, Max(2, 1))
	assert.Equal(t, 0, Max(-1, 0))
	assert.Equal(t, -1, Max(-1, -2))
}
