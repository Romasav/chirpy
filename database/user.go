package database

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func NewUser(id int, email string) (*User, error) {
	newUser := User{
		ID:    id,
		Email: email,
	}
	return &newUser, nil
}
