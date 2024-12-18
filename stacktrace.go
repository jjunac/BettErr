package betterr

import "runtime"

type Stacktrace interface {
	GetFrames() []StackFrames
	FramesLen() int
}

type StackFrames struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

var _ Stacktrace = (*RuntimeStacktrace)(nil)

type RuntimeStacktrace struct {
	Stack []uintptr
}

func NewRuntimeStacktrace(skip int) Stacktrace {
	var pcs [32]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	return &RuntimeStacktrace{
		Stack: pcs[:n],
	}
}

func (s RuntimeStacktrace) GetFrames() []StackFrames {
	frames := runtime.CallersFrames(s.Stack)
	var frameList []StackFrames
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		frameList = append(frameList, StackFrames{
			File:     frame.File,
			Function: frame.Function,
			Line:     frame.Line,
		})
	}
	return frameList
}

func (s RuntimeStacktrace) FramesLen() int {
	return len(s.Stack)
}
