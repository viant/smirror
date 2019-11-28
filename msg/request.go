package msg


//Request represents a service request
type Request struct {
	EventID string
	Data []byte `json:"data"`
}
