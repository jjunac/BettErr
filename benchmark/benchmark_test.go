package benchmark

import (
	"errors"
	"runtime"
	"testing"

	"github.com/jjunac/betterr"
	"github.com/joomcode/errorx"
	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
)

func recursiveError(n int, errorFunc func() error) error {
	// We return the error also if n == 1 because since the errorFunc is a function too, it adds an additional frame
	if n <= 1 {
		return errorFunc()
	}
	return recursiveError(n-1, errorFunc)
}

var ErrorFrameworks = []struct {
	Name string
	Func func() error
}{
	{
		Name: "Betterr",
		Func: func() error {
			return betterr.New("A BetterError error")
		},
	},
	{
		Name: "Eris",
		Func: func() error {
			return eris.New("An Eris error")
		},
	},
	{
		Name: "Errorx",
		Func: func() error {
			return errorx.IllegalState.New("An Errorx error")
		},
	},
	{
		Name: "Errors",
		Func: func() error {
			return errors.New("A plain Go error")
		},
	},
}

func Benchmark_Stack10(b *testing.B) {
	for _, ef := range ErrorFrameworks {
		b.Run(ef.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				recursiveError(10, ef.Func)
			}
		})
	}
}

func Benchmark_Stack100(b *testing.B) {
	for _, ef := range ErrorFrameworks {
		b.Run(ef.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				recursiveError(100, ef.Func)
			}
		})
	}
}

func TestRecursiveError(t *testing.T) {
	err := recursiveError(10, func() error {
		return betterr.New("A BetterError error")
	})
	assert.Error(t, err)
	assert.IsType(t, &betterr.BetterError{}, err)
	betterror := err.(*betterr.BetterError)
	assert.Equal(t, "A BetterError error", betterror.Msg)

	// We want to check that the stack is indeed 10 frames long, so we have to remove the frames from the test function
	nb_frame_test := runtime.Callers(0, make([]uintptr, 64))
	assert.Equal(t, 10, betterror.Stack.FramesLen()-nb_frame_test)
}
