package betterr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoStyleFormatter(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		err          error
		expectedGo   string
		expectedJava string
	}{
		{
			err:          errors.New("A plain Go error"),
			expectedGo:   "A plain Go error",
			expectedJava: "A plain Go error",
		},
		{
			err:        New("A BetterError error"),
			expectedGo: "A BetterError error",
			expectedJava: "A BetterError error\n" +
				"    at github.com/jjunac/betterr.TestGoStyleFormatter (/Users/jjunac/dev/betterr/betterr_test.go:24)\n" +
				"    at testing.tRunner (/usr/local/go/src/testing/testing.go:1576)\n",
		},
		{
			err:        Wrap(errors.New("A wrapped plain Go error")),
			expectedGo: "A wrapped plain Go error",
			expectedJava: "A wrapped plain Go error\n" +
				"    at github.com/jjunac/betterr.TestGoStyleFormatter (/Users/jjunac/dev/betterr/betterr_test.go:31)\n" +
				"    at testing.tRunner (/usr/local/go/src/testing/testing.go:1576)\n",
		},
		{
			err:        Wrap(New("A wrapped BetterError error")),
			expectedGo: "A wrapped BetterError error",
			expectedJava: "A wrapped BetterError error\n" +
				"    at github.com/jjunac/betterr.TestGoStyleFormatter (/Users/jjunac/dev/betterr/betterr_test.go:38)\n" +
				"    at testing.tRunner (/usr/local/go/src/testing/testing.go:1576)\n",
		},
		{
			err:        Decorate(errors.New("A plain Go error"), "Decorated"),
			expectedGo: "Decorated: A plain Go error",
			expectedJava: "Decorated\n" +
				"    at github.com/jjunac/betterr.TestGoStyleFormatter (/Users/jjunac/dev/betterr/betterr_test.go:45)\n" +
				"    at testing.tRunner (/usr/local/go/src/testing/testing.go:1576)\n" +
				"Caused by: A plain Go error",
		},
		{
			err:        Decorate(New("A BetterError error"), "Decorated"),
			expectedGo: "Decorated: A BetterError error",
			expectedJava: "Decorated\n" +
				"    at github.com/jjunac/betterr.TestGoStyleFormatter (/Users/jjunac/dev/betterr/betterr_test.go:53)\n" +
				"    at testing.tRunner (/usr/local/go/src/testing/testing.go:1576)\n" +
				"Caused by: A BetterError error\n" +
				"    at github.com/jjunac/betterr.TestGoStyleFormatter (/Users/jjunac/dev/betterr/betterr_test.go:53)\n" +
				"    at testing.tRunner (/usr/local/go/src/testing/testing.go:1576)\n",
		},
		{
			err:        Decoratef(Decorate(errors.New("A plain Go error"), "Decorated"), "A %s level of decoration", "second"),
			expectedGo: "A second level of decoration: Decorated: A plain Go error",
			expectedJava: "A second level of decoration\n" +
				"    at github.com/jjunac/betterr.TestGoStyleFormatter (/Users/jjunac/dev/betterr/betterr_test.go:63)\n" +
				"    at testing.tRunner (/usr/local/go/src/testing/testing.go:1576)\n" +
				"Caused by: Decorated\n" +
				"    at github.com/jjunac/betterr.TestGoStyleFormatter (/Users/jjunac/dev/betterr/betterr_test.go:63)\n" +
				"    at testing.tRunner (/usr/local/go/src/testing/testing.go:1576)\n" +
				"Caused by: A plain Go error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedGo, func(t *testing.T) {
			assert.Equal(tc.expectedGo, new(GoStyleFormatter).Format(tc.err))
			assert.Equal(tc.expectedJava, new(JavaStyleFormatter).Format(tc.err))
		})
	}
}
