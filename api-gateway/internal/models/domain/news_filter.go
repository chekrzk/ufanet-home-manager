package domain

type NewsFilter struct {
	Pagination Pagination
	DateFrom   string
	DateTo     string
}

type CreateNews struct {
	Title   string
	Body    string
	HouseID string
}
