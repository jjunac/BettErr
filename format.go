package betterr

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Interface to format error in string.
// Implement this interface to create custom error formatters.
// The library provides the following formatters:
// - [GoStyleFormatter]
// - [JavaStyleFormatter]
// - [JsonFormatter]
type ErrorFormatter interface {
	Format(err error) string
}

// DefaultFortmatter is the default formatter used by BetterError.Error().
// You can change the default formatter by setting this variable.
// By default, it uses [JavaStyleFormatter].
var DefaultFortmatter ErrorFormatter = &JavaStyleFormatter{}

// Formats the error in Java style.
// Example:
//   failed to process: something went wrong
type GoStyleFormatter struct {
}
var _ ErrorFormatter = (*GoStyleFormatter)(nil)
func (f *GoStyleFormatter) Format(err error) string {
	sb := strings.Builder{}
	sb.Len()
	curr := err
	for curr != nil {
		if sb.Len() > 0 {
			sb.WriteString(": ")
		}
		if betterr, ok := curr.(*BetterError); ok {
			sb.WriteString(betterr.Msg)
			curr = betterr.Wrapped
		} else {
			sb.WriteString(curr.Error())
			break
		}
	}
	return sb.String()
}

// Formats the error in Java style.
// Example:
//   failed to process
//       at github.com/myapp.MyFunction (file.go:123)
//       at github.com/myapp.main (main.go:45)
//   Caused by: something went wrong
//       at github.com/myapp.OtherFunction (file.go:100)
type JavaStyleFormatter struct {
}
var _ ErrorFormatter = (*JavaStyleFormatter)(nil)
func (f *JavaStyleFormatter) Format(err error) string {
	sb := strings.Builder{}
	sb.Len()
	curr := err
	for curr != nil {
		if sb.Len() > 0 {
			sb.WriteString("Caused by: ")
		}
		if betterr, ok := curr.(*BetterError); ok {
			sb.WriteString(betterr.Msg)
			sb.WriteByte('\n')
			for _, frame := range betterr.Stack.GetFrames() {
				sb.WriteString("    at ")
				sb.WriteString(frame.Function)
				sb.WriteString(" (")
				sb.WriteString(frame.File)
				sb.WriteByte(':')
				sb.WriteString(strconv.Itoa(frame.Line))
				sb.WriteString(")\n")
			}
			curr = betterr.Wrapped
		} else {
			sb.WriteString(curr.Error())
			break
		}
	}
	return sb.String()
}



// Formats the error in JSON.
// Example:
//   {
//       "message": "failed to process",
//       "stack": [
//           {
//               "file": "file.go",
//               "function": "github.com/myapp.MyFunction",
//               "line": 123
//           },
//           {
//               "file": "main.go",
//               "function": "github.com/myapp.main",
//               "line": 45
//           }
//       ],
//       "cause": {
//           "message": "something went wrong",
//           "stack": [
//               {
//                   "file": "file.go",
//                   "function": "github.com/myapp.OtherFunction",
//                   "line": 42
//               }
//           ]
//       }
//   }
type JsonFormatter struct {
}
var _ ErrorFormatter = (*JavaStyleFormatter)(nil)
func (f *JsonFormatter) Format(err error) string {
	type jsonError struct {
		Message string      `json:"message"`
		Stack  []StackFrames `json:"stack,omitempty"`
		Cause   *jsonError  `json:"cause,omitempty"`
	}

	curr := err
	root := &jsonError{}
	current := root

	for curr != nil {
		if betterr, ok := curr.(*BetterError); ok {
			current.Message = betterr.Msg
			current.Stack = betterr.Stack.GetFrames()

			if betterr.Wrapped != nil {
				current.Cause = &jsonError{}
				current = current.Cause
			}
			curr = betterr.Wrapped
		} else {
			current.Message = curr.Error()
			break
		}
	}

	result, _ := json.Marshal(root)
	return string(result)
}
