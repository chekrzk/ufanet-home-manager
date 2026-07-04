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

type WorkerAvailability struct {
	ID             string `json:"id"`
	WorkerID       string `json:"worker_id,omitempty"`
	UserID         string `json:"user_id"`
	Specialization string `json:"specialization"`
	HouseID        string `json:"house_id"`
	AvailableDate  string `json:"available_date"`
	AvailableTime  string `json:"available_time"`
}

type SetWorkerAvailability struct {
	Specialization string
	HouseID        string
	AvailableDate  string
	AvailableTime  string
}

type WorkerAvailabilityFilter struct {
	HouseID        string
	Specialization string
	AvailableDate  string
}
