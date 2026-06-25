package dto

type ListNewsRequest struct {
	Pagination
	DateFrom string `query:"date_from" json:"date_from,omitempty"`
	DateTo   string `query:"date_to" json:"date_to,omitempty"`
}

type CreateNewsRequest struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	HouseID string `json:"house_id,omitempty"`
}

func (r CreateNewsRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.Title == "" {
		errs["title"] = "required"
	}
	if r.Body == "" {
		errs["body"] = "required"
	}
	return errs
}
