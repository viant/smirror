package proxy

var singleton Proxy

//Singleton returns a proxy
func Singleton() (Proxy, error) {
	if singleton != nil {
		return singleton, nil
	}
	var err error
	singleton, err = New()
	return singleton, err
}
