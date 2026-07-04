package models

type RegisterCommand struct {
	Phone    string
	Password string
	Role     string
}

type LoginCommand struct {
	Phone    string
	Password string
}
