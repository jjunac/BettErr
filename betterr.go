package betterr

import "fmt"

type BetterError struct {
	Msg     string
	Wrapped error
	Stack   []uintptr
}

func New(msg string) *BetterError {
	return &BetterError{
		Msg:   msg,
		Stack: GetCallers(1),
	}
}

// If this is a BetterError, resets the stack.
// Otherwise, wraps the error in a BetterError.
func Wrap(err error) *BetterError {
	if betterr, ok := err.(*BetterError); ok {
		betterr.Stack = GetCallers(1)
		return betterr
	} else {
		return &BetterError{
			Msg:   err.Error(),
			Stack: GetCallers(1),
		}
	}
}

func Decorate(err error, msg string) *BetterError {
	return &BetterError{
		Msg:   msg,
		Wrapped: err,
		Stack: GetCallers(1),
	}
}

func Decoratef(err error, format string, args ...any) *BetterError {
	return &BetterError{
		Msg:   fmt.Sprintf(format, args...),
		Wrapped: err,
		Stack: GetCallers(1),
	}
}

func (e *BetterError) Error() string {
	return e.Format(DefaultFortmatter)
}

func (e *BetterError) Format(errFmt ErrorFormatter) string {
	return errFmt.Format(e)
}
