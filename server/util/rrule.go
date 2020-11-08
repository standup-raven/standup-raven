package util

import (
	"time"

	"github.com/teambition/rrule-go"
)

func ParseRRuleFromString(rruleString string, startDate time.Time) (*rrule.RRule, error) {
	// parse rrule
	rruleOptions, err := rrule.StrToROption(rruleString)
	if err != nil {
		return nil, err
	}

	rule, err := rrule.NewRRule(*rruleOptions)
	if err != nil {
		return nil, err
	}

	rule.DTStart(startDate)
	return rule, nil
}
