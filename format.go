package betterr

import (
	"runtime"
	"strconv"
	"strings"
)

type ErrorFormatter interface {
	Format(err error) string
}

var DefaultFortmatter ErrorFormatter = &JavaStyleFormatter{}

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
			// Iterates on Frame to output the File, Line, and Function.
			frames := runtime.CallersFrames(betterr.Stack)
			for {
				frame, ok := frames.Next()
				if !ok {
					break
				}
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
