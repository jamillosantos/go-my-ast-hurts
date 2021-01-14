package models

type Test struct {
	ID   int64 `json:"id"`
	User User  `json:"user"`
}

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
