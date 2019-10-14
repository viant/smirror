package auth

//Credentials represents credentials
type Credentials struct {
	Secret
	OAuthToken
	Auth []byte
}
