package betterr

import (
	"fmt"
)

type BetterError struct {
	Msg     string
	Wrapped error
	Stack   Stacktrace
}

var GetStacktrace func(skip int)Stacktrace = NewRuntimeStacktrace

// Creates a new BetterError with the provided message.
// The stack trace will start from the caller of this function.
// This would be the equivalent of Go's errors.New(msg) or Java's new Exception(msg).
func New(msg string) *BetterError {
	return &BetterError{
		Msg:   msg,
		Stack: GetStacktrace(1),
	}
}

// Wraps the error in a BetterError.
// Usually used to wrap errors that are not BetterError.
// The stack trace will start from the caller of this function.
// If the error is already a BetterError, it will return the error as is.
func Wrap(err error) *BetterError {
	if betterr, ok := err.(*BetterError); ok {
		return betterr
	} else {
		return &BetterError{
			Msg:   err.Error(),
			Stack: GetStacktrace(1),
		}
	}
}

// Wraps the error in a BetterError and adds a message.
// Usually used to add a message to an existing error, to provide more context.
// This would be the equivalent of Go's fmt.Errorf("%s: %w", msg, err), or Java's new Exception(msg, err).
func Decorate(err error, msg string) *BetterError {
	return &BetterError{
		Msg:   msg,
		Wrapped: err,
		Stack: GetStacktrace(1),
	}
}

// Wraps the error in a BetterError and adds a formatted message.
// See [Decorate] for more information.
func Decoratef(err error, format string, args ...any) *BetterError {
	return &BetterError{
		Msg:   fmt.Sprintf(format, args...),
		Wrapped: err,
		Stack: GetStacktrace(1),
	}
}

// Formats the error using the default formatter (JavaStyleFormatter by default).
// You can change the default formatter by setting [DefaultFortmatter]
func (e *BetterError) Error() string {
	return e.Format(DefaultFortmatter)
}

// Formats the error using the provided [ErrorFormatter]
func (e *BetterError) Format(errFmt ErrorFormatter) string {
	return errFmt.Format(e)
}
