package domain

type LoginCredentials struct {
	Phone    string
	Password string
}

type RegisterUser struct {
	Phone     string
	Password  string
	FullName  string
	HouseID   string
	Apartment string
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
