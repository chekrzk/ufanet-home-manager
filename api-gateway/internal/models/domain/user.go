package domain

type User struct {
	ID        string `json:"id"`
	Phone     string `json:"phone"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	HouseID   string `json:"house_id,omitempty"`
	Apartment string `json:"apartment,omitempty"`
}

type AuthContext struct {
	UserID string
	Role   string
}
