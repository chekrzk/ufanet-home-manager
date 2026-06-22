package domain

type RegisterDevice struct {
	Token    string
	Platform string
}

type UnregisterDevice struct {
	Token string
}
