package auth

//Credentials represents credentials
type Credentials struct {
	Secret
	Auth []byte
}
