package domain

type UpdateProfile struct {
	FullName  string
	HouseID   string
	Apartment string
}

type Worker struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id,omitempty"`
	FullName       string `json:"full_name"`
	Specialization string `json:"specialization"`
	Phone          string `json:"phone,omitempty"`
	HouseID        string `json:"house_id,omitempty"`
}

type AddWorker struct {
	UserID         string
	FullName       string
	Specialization string
	Phone          string
	HouseID        string
}
