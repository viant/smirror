package mon

var singleton Service
var singletonEnvKey string

//NewFromEnv returns singleton service for env key
func NewFromEnv(envKey string) (Service, error) {
	if singleton != nil && envKey == singletonEnvKey {
		return singleton, nil
	}
	config := &Config{}
	config.Init()
	service := New(config)
	singletonEnvKey = envKey
	singleton = service
	return singleton, nil
}
