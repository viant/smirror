package proxy

//Response represents proxy response
type Response struct {
	Source      string
	Destination string
	Triggered   map[string]string `json:",omitempty"`
	ProxyMethod string
	ProxyType   string
	Status      string
	Error       string
}
