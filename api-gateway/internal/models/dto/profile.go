package dto

type UpdateProfileRequest struct {
	FullName  string `json:"full_name"`
	HouseID   string `json:"house_id"`
	Apartment string `json:"apartment"`
}

type AddWorkerRequest struct {
	UserID         string `json:"user_id"`
	FullName       string `json:"full_name"`
	Specialization string `json:"specialization"`
	Phone          string `json:"phone"`
	HouseID        string `json:"house_id"`
}

func (r AddWorkerRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.FullName == "" {
		errs["full_name"] = "required"
	}
	if r.Specialization == "" {
		errs["specialization"] = "required"
	}
	return errs
}
