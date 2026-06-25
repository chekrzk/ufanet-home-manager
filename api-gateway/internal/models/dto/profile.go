package dto

type UpdateProfileRequest struct {
	FullName  string `json:"full_name"`
	Apartment string `json:"apartment"`
}
