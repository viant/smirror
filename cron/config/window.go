package config

import "time"

//TimeWindow represents resource asset last modification loopback time window
type TimeWindow struct {
	Duration      time.Duration
	DurationInSec int
}
