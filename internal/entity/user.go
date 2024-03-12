package entity

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     `json:"role"`
}

type Role struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
