package betterr

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestGoStyleFormatter(t *testing.T) {
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
				"    at github\\.com/jjunac/betterr\\.TestGoStyleFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\n",
		},
		{
			err:             Wrap(errors.New("A wrapped plain Go error")),
			expectedRegexGo: "A wrapped plain Go error",
			expectedRegexJava: "A wrapped plain Go error\n" +
				"    at github\\.com/jjunac/betterr\\.TestGoStyleFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\n",
		},
		{
			err:             Wrap(New("A wrapped BetterError error")),
			expectedRegexGo: "A wrapped BetterError error",
			expectedRegexJava: "A wrapped BetterError error\n" +
				"    at github\\.com/jjunac/betterr\\.TestGoStyleFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\n",
		},
		{
			err:             Decorate(errors.New("A plain Go error"), "Decorated"),
			expectedRegexGo: "Decorated: A plain Go error",
			expectedRegexJava: "Decorated\n" +
				"    at github\\.com/jjunac/betterr\\.TestGoStyleFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\\\n" +
				"Caused by: A plain Go error",
		},
		{
			err:             Decorate(New("A BetterError error"), "Decorated"),
			expectedRegexGo: "Decorated: A BetterError error",
			expectedRegexJava: "Decorated\n" +
				"    at github\\.com/jjunac/betterr\\.TestGoStyleFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\\\n" +
				"Caused by: A BetterError error\n" +
				"    at github\\.com/jjunac/betterr\\.TestGoStyleFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\n",
		},
		{
			err:             Decoratef(Decorate(errors.New("A plain Go error"), "Decorated"), "A %s level of decoration", "second"),
			expectedRegexGo: "A second level of decoration: Decorated: A plain Go error",
			expectedRegexJava: "A second level of decoration\n" +
				"    at github\\.com/jjunac/betterr\\.TestGoStyleFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\\\n" +
				"Caused by: Decorated\n" +
				"    at github\\.com/jjunac/betterr\\.TestGoStyleFormatter \\(.*/betterr_test.go:\\d+\\)\n" +
				"    at testing\\.tRunner \\(.*/go/src/testing/testing.go:\\d+\\)\\\n" +
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

func TestGoStyleFormatter_MockedStacktrace(t *testing.T) {
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
