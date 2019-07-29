package secret

//Config represents a secret config
type Config struct {
	Provider     string
	TargetScheme string
	URL          string
	Parameter    string
	Key          string
}
