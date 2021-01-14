package models

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Home struct {
	ID      int64  `json:"id"`
	Address string `json:"address"`
	Number  byte   `json:"number"`
}
