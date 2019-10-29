package models

type HasName interface {
	Name() string
	SetName(value string)
}

type HasAge interface {
	Age() int
	SetAge(value int)
}

type HasNameWrong interface {
	Name() string
	SetName(value int)
}

type InterfaceUser struct {
	name string
}

func (u *InterfaceUser) Name() string {
	return u.name
}

func (u *InterfaceUser) SetName(value string) {
	u.name = value
}
