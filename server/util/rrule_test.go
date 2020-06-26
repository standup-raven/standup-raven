package util

import (
	"github.com/stretchr/testify/assert"
	"github.com/teambition/rrule-go"
	"testing"
	"time"
)

func TestParseRRuleFromString(t *testing.T) {
	rruleString := "FREQ=WEEKLY;INTERVAL=10;BYDAY=MO,TU,FR;COUNT=5"
	startDate := time.Date(2020, time.July, 10, 8, 28, 0, 0, time.Now().Location())
	
	rule, err := ParseRRuleFromString(rruleString, startDate)
	
	assert.Nil(t, err)
	assert.Equal(t, rrule.WEEKLY, rule.Freq)
	assert.Equal(t, 3, len(rule.Byweekday))
	assert.Equal(t, 10, rule.Interval)
	assert.Equal(t, 5, rule.Count)
	assert.Equal(t, 5, len(rule.All()))
	assert.Equal(t, 2020, rule.DateStart.Year() )
	assert.Equal(t, time.July, rule.DateStart.Month() )
	assert.Equal(t, 10, rule.DateStart.Day() )
	assert.Equal(t, 8, rule.DateStart.Hour() )
	assert.Equal(t, 28, rule.DateStart.Minute() )
	assert.Equal(t, 0, rule.DateStart.Second() )
	assert.Equal(t, 0, rule.DateStart.Nanosecond() )
	assert.Equal(t, time.Now().Location(), rule.DateStart.Location() )
	
	// invalid rrule string
	rruleString = "some-invalid-rrule-string"
	rule, err = ParseRRuleFromString(rruleString, startDate)
	assert.NotNil(t, err)
	assert.Nil(t, rule)
}
