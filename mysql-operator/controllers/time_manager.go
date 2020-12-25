package controllers

import (
	"time"
)

// TimeManager manages a set of time durations
type TimeManager []time.Duration

// NewTimeManager creates a set of time durations
func NewTimeManager() *TimeManager {
	tm := TimeManager([]time.Duration{
		1 * time.Second,
		2 * time.Second,
		2 * time.Second,
		5 * time.Second,
		5 * time.Second,
		5 * time.Second,
		10 * time.Second,
		10 * time.Second,
		10 * time.Second,
		20 * time.Second,
		30 * time.Second,
		30 * time.Second,
		60 * time.Second,
		60 * time.Second,
		120 * time.Second,
		300 * time.Second,
	})
	return &tm
}

// Next returns the next duration based on the first time
func (t *TimeManager) Next(c time.Time) time.Duration {
	d := time.Now().Sub(c)
	overall := 1 * time.Second
	result := 1 * time.Second
	if t != nil {
		for _, i := range *t {
			overall += i
			result = i
			if d >= overall {
				continue
			}
			break
		}
	}
	return result
}
