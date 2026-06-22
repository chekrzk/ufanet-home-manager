package domain

type CreateRequest struct {
	Category    string
	Description string
}

type UpdateRequestStatus struct {
	Status string
}

type AddRequestComment struct {
	Text string
}
