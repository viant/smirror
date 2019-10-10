package event

import "time"

// PubsubBucketNotification
type PubsubBucketNotification struct {
	Attributes *Attributes `json:"attributes"`
}

//Attributes represents attributes
type Attributes struct {
	NotificationConfig string     `json:"notificationConfig"`
	ObjectId           string     `json:"objectId"`
	BucketId           string     `json:"bucketId"`
	EventTime          *time.Time `json:"eventTime"`
	EventType          string     `json:"eventTime"`
}

//StorageEvent returns a storage event
func (e *PubsubBucketNotification) StorageEvent() *StorageEvent {
	if e.Attributes == nil {
		return nil
	}
	return &StorageEvent{
		Bucket: e.Attributes.BucketId,
		Name:   e.Attributes.ObjectId,
	}
}
