package models

type Home struct {
	ID      int64  `json:"id"`
	Address string `json:"address"`
	User    User   `json:"user"`
	Parent  *User  `json:"parent"`
}
