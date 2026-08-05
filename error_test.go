package errkit_test

import (
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/aileron-projects/go-errkit"
	"github.com/aileron-projects/go-tester"
)

func TestErr_Unwrap(t *testing.T) {
	t.Parallel()
	e := &errkit.Error{Cause: io.EOF}
	u := e.Unwrap()
	tester.AssertEqual(t, io.EOF, u)
}

func TestErr_Error(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		err  *errkit.Error
		want string
	}{
		"message": {
			err:  &errkit.Error{Message: "m"},
			want: "m",
		},
		"detail": {
			err:  &errkit.Error{Detail: "d"},
			want: " d",
		},
		"cause": {
			err:  &errkit.Error{Cause: io.EOF},
			want: " [EOF]",
		},
		"all": {
			err:  &errkit.Error{Cause: io.EOF, Message: "m", Detail: "d"},
			want: "m d [EOF]",
		},
		"message detail": {
			err:  &errkit.Error{Message: "m", Detail: "d"},
			want: "m d",
		},
		"message cause": {
			err:  &errkit.Error{Message: "m", Cause: io.EOF},
			want: "m [EOF]",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.err.Error()
			tester.AssertEqual(t, tc.want, got)
		})
	}
}

func TestErr_Is(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		use    *errkit.Error
		target error
		same   bool
	}{
		"nil": {
			use:    nil,
			target: nil,
			same:   false,
		},
		"nil pointer": {
			use:    nil,
			target: (*errkit.Error)(nil),
			same:   true,
		},
		"nil target": {
			use:    &errkit.Error{Cause: io.EOF, Message: "m", Detail: "d"},
			target: nil,
			same:   false,
		},
		"equal": {
			use:    &errkit.Error{Message: "m", Detail: "d"},
			target: &errkit.Error{Message: "m", Detail: "d"},
			same:   true,
		},
		"not equal": {
			use:    &errkit.Error{Cause: io.EOF, Message: "m", Detail: "d"},
			target: io.EOF,
			same:   false,
		},
		"message mismatch": {
			use:    &errkit.Error{Message: "m", Detail: "d"},
			target: &errkit.Error{Message: "M", Detail: "d"},
			same:   false,
		},
		"detail mismatch": {
			use:    &errkit.Error{Message: "m", Detail: "d"},
			target: &errkit.Error{Message: "m", Detail: "D"},
			same:   true,
		},
		"same after unwrap": {
			use:    &errkit.Error{Message: "m", Detail: "d"},
			target: fmt.Errorf("outer error [%w]", &errkit.Error{Message: "m", Detail: "d"}),
			same:   true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			is := tc.use.Is(tc.target)
			tester.AssertEqual(t, tc.same, is)
		})
	}
}

func TestErr_Map(t *testing.T) {
	t.Parallel()
	t.Run("minimum", func(t *testing.T) {
		e := &errkit.Error{Message: "m", Detail: "d"}
		got := e.Map()
		want := map[string]any{
			"message": "m",
			"detail":  "d",
		}
		tester.AssertDeepEqual(t, want, got)
	})
	t.Run("cause", func(t *testing.T) {
		e := &errkit.Error{Cause: io.EOF}
		got := e.Map()
		want := map[string]any{
			"message": "EOF",
		}
		tester.AssertDeepEqual(t, want, got["cause"].(map[string]any))
	})
}

func TestErr_SlogAttrs(t *testing.T) {
	t.Parallel()
	t.Run("minimum", func(t *testing.T) {
		e := &errkit.Error{Message: "m", Detail: "d"}
		got := e.SlogAttrs()
		want := []slog.Attr{
			slog.String("message", "m"),
			slog.String("detail", "d"),
		}
		tester.AssertDeepEqual(t, want, got)
	})
	t.Run("cause", func(t *testing.T) {
		e := &errkit.Error{Cause: io.EOF}
		got := e.SlogAttrs()
		want := []slog.Attr{
			slog.String("message", ""),
			slog.String("detail", ""),
			slog.GroupAttrs("cause", slog.String("message", "EOF")),
		}
		tester.AssertDeepEqual(t, want, got)
	})
}
