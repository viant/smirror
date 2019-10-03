package config

import (
	"fmt"
	"time"
)

//TimeWindow represents resource asset last modification loopback time window
type TimeWindow struct {
	Duration      time.Duration
	DurationInSec int
}

//Init initialises time window
func (t *TimeWindow) Init() {
	if t.DurationInSec != 0 {
		t.Duration = time.Duration(t.DurationInSec) * time.Second
	}
}

//Validate checks if setting is valid
func (t *TimeWindow) Validate() error {
	if t.Duration == 0 {
		return fmt.Errorf("time duration was empty")
	}
	return nil
}
