package config

const (
	megaBytes              = 1024 * 1024
	defaultStreamThreshold = 1024
	defaultPartSize        = 64
)

//Streaming represents streaming option
type Streaming struct {
	ThresholdMb             int
	threshold               int
	partSize                int
	PartSizeMb              int
	ChecksumSkipThresholdMb int
	checksumSkipThreshold   int
}

//Threshold returns download/upload streaming
func (c *Streaming) Threshold() int {
	return c.threshold
}

//PartSize download part size
func (c *Streaming)  PartSize() int {
	return c.partSize
}

//ChecksumSkipThreshold upload checksum skip threshold
func (c *Streaming) ChecksumSkipThreshold() int {
	return c.checksumSkipThreshold
}

//Init initialises streaming
func (c *Streaming) Init() {
	if c.ThresholdMb == 0 {
		c.ThresholdMb = defaultStreamThreshold
	}
	if c.threshold == 0 {
		c.threshold = c.ThresholdMb * megaBytes
	}

	if c.PartSizeMb == 0 {
		c.PartSizeMb = defaultPartSize
	}

	if c.partSize == 0 {
		c.partSize = c.PartSizeMb * megaBytes
	}

	if c.ChecksumSkipThresholdMb == 0 {
		c.ChecksumSkipThresholdMb = defaultStreamThreshold
	}
	if c.checksumSkipThreshold == 0 {
		c.checksumSkipThreshold = c.ChecksumSkipThresholdMb * megaBytes
	}
}
