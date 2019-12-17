package config

const (
	megaBytes              = 1024 * 1024
	defaultStreamThreshold = 1024
	defaultPartSize        = 64
)

//Streaming represents streaming option
type Streaming struct {
	ThresholdMb             int
	Threshold               int
	PartSize                int
	PartSizeMb              int
	ChecksumSkipThresholdMb int
	ChecksumSkipThreshold   int
}

func (c *Streaming) Init() {
	if c.ThresholdMb == 0 {
		c.ThresholdMb = defaultStreamThreshold
	}
	if c.Threshold == 0 {
		c.Threshold = c.ThresholdMb * megaBytes
	}

	if c.PartSizeMb == 0 {
		c.PartSizeMb = defaultPartSize
	}

	if c.PartSize == 0 {
		c.PartSize = c.PartSizeMb * megaBytes
	}

	if c.ChecksumSkipThresholdMb == 0 {
		c.ChecksumSkipThresholdMb = defaultStreamThreshold
	}
	if c.ChecksumSkipThreshold == 0 {
		c.ChecksumSkipThreshold = c.ChecksumSkipThresholdMb * megaBytes
	}
}
