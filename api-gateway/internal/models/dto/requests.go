package dto

type CreateRequestRequest struct {
	Category         string `json:"category"`
	Description      string `json:"description"`
	PreferredDate    string `json:"preferred_date"`
	AssignedWorkerID string `json:"assigned_worker_id"`
	Address          string `json:"address"`
	Apartment        string `json:"apartment"`
	Phone            string `json:"phone"`
}

func (r CreateRequestRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.Category == "" {
		errs["category"] = "required"
	}
	if r.Description == "" {
		errs["description"] = "required"
	}
	return errs
}

type UpdateRequestStatusRequest struct {
	Status     string `json:"status"`
	AssignedTo string `json:"assigned_to"`
}

func (r UpdateRequestStatusRequest) Validate() map[string]string {
	if r.Status == "" {
		return map[string]string{"status": "required"}
	}
	return nil
}

type AddRequestCommentRequest struct {
	Text string `json:"text"`
}

func (r AddRequestCommentRequest) Validate() map[string]string {
	if r.Text == "" {
		return map[string]string{"text": "required"}
	}
	return nil
}
