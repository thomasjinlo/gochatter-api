package users

type User struct {
	username string
}

func NewUser() *User {
	return &User{}
}

func (u *User) Username() string {
	return u.username
}
