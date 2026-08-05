package errkit

import (
	"strings"
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestFrame_String(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		f   Frame
		str string
	}{
		"empty": {
			f:   Frame{},
			str: "Pkg: File: Func: Line:0",
		},
		"non empty": {
			f:   Frame{Pkg: "foo/bar", File: "test.go", Func: "testFunc", Line: 100},
			str: "Pkg:foo/bar File:test.go Func:testFunc Line:100",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tester.AssertEqual(t, tc.str, tc.f.String())
		})
	}
}

func TestCallerFrame(t *testing.T) {
	t.Parallel()
	t.Run("skip=0", func(t *testing.T) {
		fs := callerFrame(0)
		s := fs.String()
		t.Log(s)
		tester.AssertEqual(t, true, strings.Contains(s, "Pkg:github.com/aileron-projects/go-errkit"))
		tester.AssertEqual(t, true, strings.Contains(s, "Func:TestCallerFrame"))
	})
	t.Run("skip=-1", func(t *testing.T) {
		fs := callerFrame(-1)
		s := fs.String()
		t.Log(s)
		tester.AssertEqual(t, true, strings.Contains(s, "Pkg:github.com/aileron-projects/go-errkit"))
		tester.AssertEqual(t, true, strings.Contains(s, "Func:callerFrame"))
	})
	t.Run("skip=1", func(t *testing.T) {
		fs := callerFrame(1)
		s := fs.String()
		t.Log(s)
		tester.AssertEqual(t, true, strings.Contains(s, "Pkg:testing"))
		tester.AssertEqual(t, true, strings.Contains(s, "Func:tRunner"))
	})
	t.Run("skip=-999", func(t *testing.T) {
		fs := callerFrame(-999)
		s := fs.String()
		t.Log(s)
		tester.AssertEqual(t, true, strings.Contains(s, "Pkg:runtime"))
		tester.AssertEqual(t, true, strings.Contains(s, "Func:Callers"))
	})
	t.Run("skip=999", func(t *testing.T) {
		fs := callerFrame(999)
		s := fs.String()
		t.Log(s)
		tester.AssertEqual(t, "Pkg: File: Func: Line:0", s)
	})
}

func TestCallerFrames(t *testing.T) {
	t.Parallel()
	t.Run("skip=0", func(t *testing.T) {
		fs := callerFrames(0)
		t.Log(fs)
		tester.AssertEqual(t, true, len(fs) > 0)
	})
	t.Run("skip=-1", func(t *testing.T) {
		fs := callerFrames(-1)
		t.Log(fs)
		tester.AssertEqual(t, true, len(fs) > 0)
	})
	t.Run("skip=1", func(t *testing.T) {
		fs := callerFrames(1)
		t.Log(fs)
		tester.AssertEqual(t, true, len(fs) > 0)
	})
	t.Run("skip=-999", func(t *testing.T) {
		fs := callerFrames(-999)
		t.Log(fs)
		tester.AssertEqual(t, true, len(fs) > 0)
	})
	t.Run("skip=999", func(t *testing.T) {
		fs := callerFrames(999)
		t.Log(fs)
		tester.AssertEqual(t, 0, len(fs))
	})
}
