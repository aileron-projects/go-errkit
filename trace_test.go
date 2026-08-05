package errkit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestTraceTo(t *testing.T) {
	traceEnabled = true
	t.Cleanup(func() {
		traceEnabled = false
	})
	testCases := map[string]struct {
		err   error
		wants []string
	}{
		"empty error": {
			err: &ErrInstance{},
			wants: []string{
				`[ERRKIT][TRACE]`,
				`"code":    ""`,
				`"kind":    ""`,
				`"message": ""`,
			},
		},
		"trace error without attrs": {
			err: &ErrInstance{Code: "TestCode", Kind: "TestKind", Message: "TestMessage"},
			wants: []string{
				`[ERRKIT][TRACE]`,
				`"code":    "TestCode"`,
				`"kind":    "TestKind"`,
				`"message": "TestMessage"`,
			},
		},
		"trace error with attrs": {
			err: &ErrInstance{Code: "TestCode", Kind: "TestKind", Message: "TestMessage", Attrs: map[string]string{"foo": "bar"}},
			wants: []string{
				`[ERRKIT][TRACE]`,
				`"code":    "TestCode"`,
				`"kind":    "TestKind"`,
				`"message": "TestMessage"`,
				`"attrs":   map[string]string{"foo":"bar"}`,
			},
		},
		"trace error with instance": {
			err: &ErrInstance{Code: "TestCode", Kind: "TestKind", Message: "TestMessage", Instance: "ABC"},
			wants: []string{
				`[ERRKIT][TRACE]`,
				`"code":     "TestCode"`,
				`"kind":     "TestKind"`,
				`"message":  "TestMessage"`,
				`"instance": "ABC"`,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			traceTo(&buf, tc.err)
			result := buf.String()
			t.Log(result)
			for _, w := range tc.wants {
				t.Log("want:", w)
				tester.AssertEqual(t, true, strings.Contains(result, w))
			}
		})
	}
}
