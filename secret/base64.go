package secret

import "encoding/base64"



func decodeBase64IfNeeded(data []byte) []byte {
	plainText := string(data)
	isBase64 := false
	if _, err := base64.StdEncoding.DecodeString(string(data)); err == nil {
		isBase64 = true
	}
	if !isBase64 {
		plainText = base64.StdEncoding.EncodeToString(data)
	}
	return []byte(plainText)
}
