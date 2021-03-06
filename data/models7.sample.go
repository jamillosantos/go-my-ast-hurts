package models

func (u *User) getName() string {
	return u.Name
}

func getName_(u *User, name string) string {
	u.Name = name
	return name
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
