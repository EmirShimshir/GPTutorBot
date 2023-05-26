package domain

type Url struct {
	Utm           string
	CountRequests int64
}

func NewUrl(utm string) *Url {
	return &Url{
		Utm:           utm,
		CountRequests: 0,
	}
}
