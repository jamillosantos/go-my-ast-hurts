// Package models is a test
package models

/* Testing comment
new line
*/
type Compra struct {
	// ID comment
	ID int64 `json:"id"`
	/* Comment with multilines
	Testing
	*/
	User User `json:"user"`
}

type User struct {
	ID int64 `json:"id"`
	// Line 1
	// Line 2
	Name  string `json:"name"`
	Email string `json:"email"`
	Data  interface{}
	Data2 struct {
		Name string
	}
}
