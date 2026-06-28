package domain

type CreateRequest struct {
	Category    string
	Description string
}

type UpdateRequestStatus struct {
	Status     string
	AssignedTo string
}

type AddRequestComment struct {
	Text string
}
