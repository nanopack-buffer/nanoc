package parser

type SyntaxError struct {
	Msg           string
	OffendingCode string
}

func (s SyntaxError) Error() string {
	return s.Msg + ". At: " + s.OffendingCode
}
