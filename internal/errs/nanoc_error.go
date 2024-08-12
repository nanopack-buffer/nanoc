package errs

import (
	"errors"
	"fmt"
	"strings"
)

type NanocError struct {
	Msg   string
	Stack []string
}

func (s *NanocError) Error() string {
	var sb strings.Builder
	for _, v := range s.Stack {
		sb.WriteString(fmt.Sprintf("    in: %v\n", v))
	}
	return s.Msg + ".\n" + sb.String()
}

func (s *NanocError) PrependToStack(entry string) {
	s.Stack = append([]string{entry}, s.Stack...)
}

func NewNanocError(msg string, stack ...string) *NanocError {
	return &NanocError{
		Msg:   msg,
		Stack: stack,
	}
}

// WrapNanocErr prepends the given item to the syntax error stack if err is NanocError.
// Otherwise, it simply returns err as is.
func WrapNanocErr(err error, expr string) error {
	var nanocErr *NanocError
	if errors.As(err, &nanocErr) {
		nanocErr.PrependToStack(expr)
		return nanocErr
	}
	return err
}
