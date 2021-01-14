package models

import "fmt"

var (
	a string
	b byte
	c int
	d int64
	e float32
	f boolean
	g User
	h []string
)

type User struct {
	ID int64 `json:"id"`
	// Line 1
	// Line 2
	Name  string `json:"name"`
	Email string `json:"email"`
}
