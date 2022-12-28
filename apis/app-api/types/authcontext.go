package types

type AuthContext struct {
	PK        string `json:"pk"`
	TokenType string `json:"tokenType"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}
