package betterr

import (
	"strings"
	"testing"
)

func method_1deep() *BetterError {
	return New("A BetterError error").(*BetterError)
}

func method_2deep() *BetterError {
	return method_2deep_nested()
}

func method_2deep_nested() *BetterError {
	return New("A nested error").(*BetterError)
}

func TestRuntimeStacktrace_1Deep(t *testing.T) {
	err := method_1deep()
	stack := err.Stack.GetFrames()
	assertEqual(t, 4, len(stack))
	assertEqual(t, 4, err.Stack.FramesLen())
	// Frame 1
	assertEqual(t, "github.com/jjunac/betterr.method_1deep", stack[0].Function)
	assertEqual(t, 9, stack[0].Line)
	assertTrue(t, strings.HasSuffix(stack[0].File, "/stacktrace_test.go"))
	// The rest is test framework frames
}

func TestRuntimeStacktrace_2Deep(t *testing.T) {
	err := method_2deep()
	stack := err.Stack.GetFrames()
	assertEqual(t, 5, len(stack))
	assertEqual(t, 5, err.Stack.FramesLen())

	// Frame 1
	assertEqual(t, "github.com/jjunac/betterr.method_2deep_nested", stack[0].Function)
	assertEqual(t, 17, stack[0].Line)
	assertTrue(t, strings.HasSuffix(stack[0].File, "/stacktrace_test.go"))
	// Frame 2
	assertEqual(t, "github.com/jjunac/betterr.method_2deep", stack[1].Function)
	assertEqual(t, 13, stack[1].Line)
	assertTrue(t, strings.HasSuffix(stack[1].File, "/stacktrace_test.go"))
	// The rest is test framework frames
}
