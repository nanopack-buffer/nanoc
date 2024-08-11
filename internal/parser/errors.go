package parser

import (
	"errors"
	"fmt"
	"strings"
)

type SyntaxError struct {
	Msg   string
	Stack []string
}

func (s *SyntaxError) Error() string {
	var sb strings.Builder
	for _, v := range s.Stack {
		sb.WriteString(fmt.Sprintf("    in: %v\n", v))
	}
	return s.Msg + ".\n" + sb.String()
}

func (s *SyntaxError) PrependToStack(entry string) {
	s.Stack = append([]string{entry}, s.Stack...)
}

func NewSyntaxError(msg string, stack ...string) *SyntaxError {
	return &SyntaxError{
		Msg:   msg,
		Stack: stack,
	}
}

// wrapSyntaxErr prepends the given item to the syntax error stack if err is SyntaxError.
// Otherwise, it simply returns err as is.
func wrapSyntaxErr(err error, expr string) error {
	var syntaxErr *SyntaxError
	if errors.As(err, &syntaxErr) {
		syntaxErr.PrependToStack(expr)
		return syntaxErr
	}
	return err
}
