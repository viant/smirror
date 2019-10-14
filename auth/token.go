package auth

//OAuthToken represents OAuth slack token
type OAuthToken struct {
	Token string `json:",omitempty"`
	RawToken string `json:"-"`
}
