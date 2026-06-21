package domain

type Device struct {
	UserID   string `json:"user_id"`
	Token    string `json:"token"`
	Platform string `json:"platform"`
}
