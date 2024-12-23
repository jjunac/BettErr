package betterr

import (
	"errors"
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
func New(msg string) error {
	return &BetterError{
		Msg:   msg,
		Stack: GetStacktrace(1),
	}
}

// Wraps the error in a BetterError.
// Usually used to wrap errors that are not BetterError.
// The stack trace will start from the caller of this function.
// If the error is already a BetterError, it will return the error as is.
//  Wrapping a nil error will return nil.
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	if betterr, ok := err.(*BetterError); ok {
		return betterr
	} else {
		return &BetterError{
			Msg:   err.Error(),
			Stack: GetStacktrace(1),
		}
	}
}

// Decorates the error in a BetterError and adds a message.
// Usually used to add a message to an existing error, to provide more context.
// This would be the equivalent of Go's fmt.Errorf("%s: %w", msg, err), or Java's new Exception(msg, err).
// Decorating a nil error will return nil.
func Decorate(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &BetterError{
		Msg:   msg,
		Wrapped: err,
		Stack: GetStacktrace(1),
	}
}

// Decorates the error in a BetterError and adds a formatted message.
// See [Decorate] for more information.
func Decoratef(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return &BetterError{
		Msg:   fmt.Sprintf(format, args...),
		Wrapped: err,
		Stack: GetStacktrace(1),
	}
}

func Is(err, target error) bool {
    return errors.Is(err, target)
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

// Is reports whether any error in err's tree matches target, as defined by the [errors.Is] interface.
// For now, it checks if the target error message is equal to one of the error messages in the tree.
func (e *BetterError) Is(target error) bool {
	if e == target {
		return true
	}
	if betterrTarget, ok := target.(*BetterError); ok {
		if e.Msg == betterrTarget.Msg {
			return true
		}
		if e.Wrapped != nil {
			return Is(e.Wrapped, target)
		}
	}
	if e.Msg == target.Error() {
		return true
	}
	if e.Wrapped != nil {
		return Is(e.Wrapped, target)
	}
	return false
}


// Returns the wrapped error, if any, as defined by the [errors.Unwrap] interface.
func (e *BetterError) Unwrap() error {
	return e.Wrapped
}
