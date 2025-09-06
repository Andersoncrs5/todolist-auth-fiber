package res

type ResponseToken struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refresh_token"`
}