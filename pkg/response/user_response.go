package response

type (
	UserResponse struct {
		Data []struct {
			BroadcasterType string `json:"broadcaster_type"`
			DisplayName     string `json:"display_name"`
			Description     string `json:"description"`
			Id              string `json:"id"`
		} `json:"data"`
	}
)
