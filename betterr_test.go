package betterr

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestErrorFormatter(t *testing.T) {
	testCases := []struct {
		err               error
		expectedRegexGo   string
		expectedRegexJava string
		expectedRegexJson map[string]any
	}{
		{
			err:               errors.New("A plain Go error"),
			expectedRegexGo:   "A plain Go error",
			expectedRegexJava: "A plain Go error",
			expectedRegexJson: map[string]any{
				"message": "A plain Go error",
			},
		},
		{
			err:             New("A BetterError error"),
			expectedRegexGo: "A BetterError error",
			expectedRegexJava: "A BetterError error\n" +
				"    at github\\.com/jjunac/betterr\\.TestErrorFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\n" +
				"    at runtime\\.goexit \\(.*\\)\n",
		},
		{
			err:             Wrap(errors.New("A wrapped plain Go error")),
			expectedRegexGo: "A wrapped plain Go error",
			expectedRegexJava: "A wrapped plain Go error\n" +
				"    at github\\.com/jjunac/betterr\\.TestErrorFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\n" +
				"    at runtime\\.goexit \\(.*\\)\n",
		},
		{
			err:             Wrap(New("A wrapped BetterError error")),
			expectedRegexGo: "A wrapped BetterError error",
			expectedRegexJava: "A wrapped BetterError error\n" +
				"    at github\\.com/jjunac/betterr\\.TestErrorFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\n" +
				"    at runtime\\.goexit \\(.*\\)\n",
		},
		{
			err:             Decorate(errors.New("A plain Go error"), "Decorated"),
			expectedRegexGo: "Decorated: A plain Go error",
			expectedRegexJava: "Decorated\n" +
				"    at github\\.com/jjunac/betterr\\.TestErrorFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\\\n" +
				"    at runtime\\.goexit \\(.*\\)\n" +
				"Caused by: A plain Go error",
		},
		{
			err:             Decorate(New("A BetterError error"), "Decorated"),
			expectedRegexGo: "Decorated: A BetterError error",
			expectedRegexJava: "Decorated\n" +
				"    at github\\.com/jjunac/betterr\\.TestErrorFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\\\n" +
				"    at runtime\\.goexit \\(.*\\)\n" +
				"Caused by: A BetterError error\n" +
				"    at github\\.com/jjunac/betterr\\.TestErrorFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\n" +
				"    at runtime\\.goexit \\(.*\\)\n",
		},
		{
			err:             Decoratef(Decorate(errors.New("A plain Go error"), "Decorated"), "A %s level of decoration", "second"),
			expectedRegexGo: "A second level of decoration: Decorated: A plain Go error",
			expectedRegexJava: "A second level of decoration\n" +
				"    at github\\.com/jjunac/betterr\\.TestErrorFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\\\n" +
				"    at runtime\\.goexit \\(.*\\)\n" +
				"Caused by: Decorated\n" +
				"    at github\\.com/jjunac/betterr\\.TestErrorFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\\\n" +
				"    at runtime\\.goexit \\(.*\\)\n" +
				"Caused by: A plain Go error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedRegexGo, func(t *testing.T) {
			assertRegexp(t, tc.expectedRegexGo, new(GoStyleFormatter).Format(tc.err))
			assertRegexp(t, tc.expectedRegexJava, new(JavaStyleFormatter).Format(tc.err))
			// Testing JSON is shitty so we don't do it every time
			if tc.expectedRegexJson != nil {
				expectedJson, err := json.Marshal(tc.expectedRegexJson)
				assertNoError(t, err)
				assertJSONEq(t, string(expectedJson), new(JsonFormatter).Format(tc.err))
			}
		})
	}
}

type mockedStacktrace struct {
	frames []StackFrames
}

func (s *mockedStacktrace) GetFrames() []StackFrames {
	return s.frames
}

func (s *mockedStacktrace) FramesLen() int {
	return len(s.frames)
}

func TestErrorFormatter_MockedStacktrace(t *testing.T) {
	GetStacktrace = func(skip int) Stacktrace {
		return &mockedStacktrace{
			frames: []StackFrames{
				{
					File:     "file.go",
					Function: "github.com/myapp.OtherFunction",
					Line:     42,
				},
			},
		}
	}

	defer func() {
		GetStacktrace = NewRuntimeStacktrace
	}()

	baseErr := New("something went wrong")

	GetStacktrace = func(skip int) Stacktrace {
		return &mockedStacktrace{
			frames: []StackFrames{
				{
					File:     "file.go",
					Function: "github.com/myapp.MyFunction",
					Line:     123,
				},
				{
					File:     "main.go",
					Function: "github.com/myapp.main",
					Line:     45,
				},
			},
		}
	}

	decoratedErr := Decorate(baseErr, "process failed")

	assertEqual(t,
		"process failed: something went wrong",
		new(GoStyleFormatter).Format(decoratedErr))

	assertEqual(t,
		"process failed\n"+
			"    at github.com/myapp.MyFunction (file.go:123)\n"+
			"    at github.com/myapp.main (main.go:45)\n"+
			"Caused by: something went wrong\n"+
			"    at github.com/myapp.OtherFunction (file.go:42)\n",
		new(JavaStyleFormatter).Format(decoratedErr))

	expectedJson := map[string]any{
		"message": "process failed",
		"stack": []map[string]any{
			{
				"function": "github.com/myapp.MyFunction",
				"file":     "file.go",
				"line":     123,
			},
			{
				"function": "github.com/myapp.main",
				"file":     "main.go",
				"line":     45,
			},
		},
		"cause": map[string]any{
			"message": "something went wrong",
			"stack": []map[string]any{
				{
					"function": "github.com/myapp.OtherFunction",
					"file":     "file.go",
					"line":     42,
				},
			},
		},
	}
	expectedJsonBytes, err := json.Marshal(expectedJson)
	assertNoError(t, err)
	assertJSONEq(t, string(expectedJsonBytes), new(JsonFormatter).Format(decoratedErr))
}

func TestIs(t *testing.T) {
	testCases := []struct {
		name      string
		targetErr error
	}{
		{
			"A plain Go error",
			errors.New("A plain Go error"),
		},
		{
			"A better error",
			New("A better error"),
		},
		{
			"A wrapped plain Go error",
			Wrap(errors.New("A wrapped plain Go error")),
		},
		{
			"A wrapped better error",
			Wrap(New("A wrapped better error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// We should not match nil or other errors
			assertFalse(t, Is(nil, tc.targetErr))
			assertFalse(t, Is(tc.targetErr, nil))
			assertFalse(t, Is(New("Another error"), tc.targetErr))
			assertFalse(t, Is(errors.New("A plain other error"), tc.targetErr))

			// We should match all the variations of the error itself, New with same message, Wrap, Decorate and Decoratef
			assertTrue(t, Is(tc.targetErr, tc.targetErr))
			assertTrue(t, Is(New(tc.name), tc.targetErr))
			assertTrue(t, Is(Wrap(tc.targetErr), tc.targetErr))
			assertTrue(t, Is(Decorate(tc.targetErr, "Decorated"), tc.targetErr))
			assertTrue(t, Is(Decoratef(tc.targetErr, "Decorated %s", "yolo"), tc.targetErr))

			// We should match as well when we call the member Is instead of the package one
			assertTrue(t, New(tc.name).(*BetterError).Is(tc.targetErr))
			assertTrue(t, Wrap(tc.targetErr).(*BetterError).Is(tc.targetErr))
			assertTrue(t, Decorate(tc.targetErr, "Decorated").(*BetterError).Is(tc.targetErr))
			assertTrue(t, Decoratef(tc.targetErr, "Decorated %s", "yolo").(*BetterError).Is(tc.targetErr))
		})
	}

}

func TestIs_WhenSubpartOfTheError(t *testing.T) {
	targetErr := New("table not found")
	// For now, we don't support this feature. Maybe we'll do eventually.
	assertFalse(t, Is(New("table not found 'test'"), targetErr))
}

func TestWrap_ShouldNotWrapNil(t *testing.T) {
	assertEqual(t, nil, Wrap(nil))
	assertTrue(t, Wrap(nil) == nil)
}

func TestDecorate_ShouldNotDecorateNil(t *testing.T) {
	assertEqual(t, nil, Decorate(nil, "message"))
	assertEqual(t, nil, Decoratef(nil, "%s", "message"))
	assertTrue(t, Decorate(nil, "message") == nil)
}
