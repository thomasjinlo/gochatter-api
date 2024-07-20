package users

type localRepository struct {
	users map[string]*User
}

func NewLocalRepository() *localRepository {
	return &localRepository{
		users: make(map[string]*User),
	}
}

func (r *localRepository) Login(username string) error {
	_, ok := r.users[username]
	if !ok {
		r.users[username] = &User{username: username}
	}
	return nil
}

func (r *localRepository) GetAll() []*User {
	var users []*User
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}
