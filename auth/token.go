package auth

//OAuthToken represents OAuth slack token
type OAuthToken struct {
	Token string
	RawToken string `json:"-"`
}
