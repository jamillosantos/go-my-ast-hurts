package myasthurts

type (
	ListenerBeforeFile interface {
		BeforeFile(*ParsePackageContext, string) error
	}

	ListenerAfterFile interface {
		AfterFile(*ParsePackageContext, string, error) error
	}
)
