package proxy

import (
	"os"
	"smirror/base"
)

var singleton Proxy

//Singleton returns a proxy
func Singleton() (Proxy, error) {
	if singleton != nil {
		return singleton, nil
	}
	method := os.Getenv(base.ProxyMethod)
	var err error
	singleton, err = New(method)
	return singleton, err
}
