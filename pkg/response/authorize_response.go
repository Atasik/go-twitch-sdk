package response

type (
	AuthorizeResponse struct {
		AccessToken    string `json:"access_token"`
		RefreshToken   string `json:"refresh_token"`
		TokenType      string `json:"token_type"`
		ExpirationDate string `json:"expires_in"`
	}
)
