package port

type Token struct {
	Value string
	Err   error
	CntReq int64
}

func NewToken(value string, err error) *Token {
	return &Token{
		Value: value,
		Err:   err,
		CntReq: 0,
	}
}

type Queue interface {
	Get() *Token
	Add(v *Token)
	Next() *Token
	GetAll() []*Token
	Remove(token string) error
}
