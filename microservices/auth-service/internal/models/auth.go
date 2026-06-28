package models

type RegisterCommand struct {
	Phone    string
	Password string
}

type LoginCommand struct {
	Phone    string
	Password string
}
