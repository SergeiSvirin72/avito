package dto

type AuthOutput struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_in"`
}
