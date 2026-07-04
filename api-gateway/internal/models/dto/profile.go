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

type SetWorkerAvailabilityRequest struct {
	Specialization string `json:"specialization"`
	HouseID        string `json:"house_id"`
	AvailableDate  string `json:"available_date"`
	AvailableTime  string `json:"available_time"`
}

func (r SetWorkerAvailabilityRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.Specialization == "" {
		errs["specialization"] = "required"
	}
	if r.HouseID == "" {
		errs["house_id"] = "required"
	}
	if r.AvailableDate == "" {
		errs["available_date"] = "required"
	}
	if r.AvailableTime == "" {
		errs["available_time"] = "required"
	}
	return errs
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
