package util

import (
	"bou.ke/monkey"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/standup-raven/standup-raven/server/otime"
	"net/http"
	"net/url"
	"testing"
	"time"
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

func TestSendEphemeralText(t *testing.T) {
	response, err := SendEphemeralText("my message")
	
	assert.Nil(t, err)
	assert.Equal(t, model.COMMAND_RESPONSE_TYPE_EPHEMERAL, response.Type)
	assert.Equal(t, "my message", response.Text)
}

func TestDifference(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, Difference([]string{"a", "b", "c", "d"}, []string{"c", "d", "e"}))
	assert.Equal(t, []string{}, Difference([]string{"a"}, []string{"a"}))
	assert.Equal(t, []string{}, Difference([]string{"a", "b", "c", "d"}, []string{"a", "b", "c", "d"}))
}

func TestGetCurrentDateString(t *testing.T) {
	monkey.Patch(otime.Now, func(timezone string) otime.OTime {
		t, _ := time.Parse("02-Jan-06", "02-Jan-06")
		return otime.OTime{t}
	})
	defer monkey.Unpatch(otime.Now)
	
	assert.Equal(t, "20060102", GetCurrentDateString("Asia/Kolkata"))
}

func TestGetKeyHash(t *testing.T) {
	assert.Equal(t, "uEPTMBNRbhuJGbFdJXMDYF1g7v9gs+Mw7y42YBpxwQA=", GetKeyHash("dummy_key"))
}

func TestDumpRequest(t *testing.T) {
	url, _ := url.Parse("https://www.example.com")
	
	assert.True(t, len(DumpRequest(&http.Request{
		URL:  url,
	})) > 0)
}

func TestContainsDuplicates(t *testing.T) {
	duplicateItem, hasDuplicate := ContainsDuplicates(&[]string{"a", "a"})
	assert.Equal(t, "a", duplicateItem)
	assert.True(t, hasDuplicate)

	duplicateItem, hasDuplicate = ContainsDuplicates(&[]string{"a", "b"})
	assert.Equal(t, "", duplicateItem)
	assert.False(t, hasDuplicate)

	duplicateItem, hasDuplicate = ContainsDuplicates(&[]string{})
	assert.Equal(t, "", duplicateItem)
	assert.False(t, hasDuplicate)
}
