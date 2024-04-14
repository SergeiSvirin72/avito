package model

const (
	UserRoleAdmin = "admin"
	UserRoleUser  = "user"
)

type User struct {
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
	ID       int64  `json:"id"`
}

func (m *User) IsAdmin() bool {
	return m.Role == UserRoleAdmin
}
