package myasthurts

type (
	ParseFileListener interface {
		BeforeFile(string) error
	}
)
