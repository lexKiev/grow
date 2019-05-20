package models


type User struct {
	Id string
	Login string
	Email string
	Password string
}

func NewUser(id, login, email, password string) *User {
	return &User{id,login, email, password}
}