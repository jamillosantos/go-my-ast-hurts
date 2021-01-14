package models

func (u *User) getName() string {
	return u.Name
}

func show(name string, age int64) {
	fmt.Println(age, name)
}

func welcome() string {
	return "Welcome"
}
