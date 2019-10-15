package proxy

//Response represents proxy response
type Response struct {
	Source      string
	Destination string
	Copy        map[string]string `json:",omitempty"`
	ProxyType   string
	Status      string
	Error       string
}
