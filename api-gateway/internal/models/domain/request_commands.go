package domain

type CreateRequest struct {
	Category         string
	Description      string
	PreferredDate    string
	AssignedWorkerID string
	Address          string
	Apartment        string
	Phone            string
}

type UpdateRequestStatus struct {
	Status     string
	AssignedTo string
}

type AddRequestComment struct {
	Text string
}
