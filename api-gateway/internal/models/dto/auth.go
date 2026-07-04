package dto

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.Phone == "" {
		errs["phone"] = "required"
	}
	if r.Password == "" {
		errs["password"] = "required"
	}
	return errs
}

type RegisterRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func (r RegisterRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.Phone == "" {
		errs["phone"] = "required"
	}
	if r.Password == "" {
		errs["password"] = "required"
	}
	if r.FullName == "" {
		errs["full_name"] = "required"
	}
	return errs
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshRequest) Validate() map[string]string {
	if r.RefreshToken == "" {
		return map[string]string{"refresh_token": "required"}
	}
	return nil
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
