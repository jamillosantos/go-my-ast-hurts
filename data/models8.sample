package models

import (
	"fmt"
	t "time"
)

type User struct {
	ID        int64  `json:"id,uuidTest"`
	Name      string `json:"name" bson:""`
	CreatedAt t.Time `json:"created_at"`
}

func (u *User) getName() string {
	fmt.Println(u.Name)
	return u.Name
}
