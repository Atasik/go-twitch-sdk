package request

type (
	AuthorizeRequest struct {
		ClientId     string `url:"client_id"`
		ClientSecret string `url:"client_secret"`
		GrantType    string `url:"grant_type"`
		Code         string `url:"code,omitempty"`
		RedirectUri  string `url:"redirect_uri,omitempty"`
	}

	SubRequest struct {
		Type      string    `json:"type,omitempty"`
		Version   string    `json:"version,omitempty"`
		Condition Condition `json:"condition,omitempty"`
		Transport Transport `json:"transport,omitempty"`
	}

	Condition struct {
		Id string `json:"broadcaster_user_id"`
	}

	Transport struct {
		Method   string `json:"method"`
		Callback string `json:"callback"`
		Secret   string `json:"secret"`
	}
)
