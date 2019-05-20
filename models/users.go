package models


type User struct {
	Login string
	Email string
	Password string
}

func NewUser(login, email, password string) *User {
	return &User{login, email, password}
}