package errkit_test

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/aileron-projects/go-errkit"
	"github.com/aileron-projects/go-tester"
)

func TestError_Unwrap(t *testing.T) {
	t.Parallel()
	e := &errkit.Error{Cause: io.EOF}
	u := e.Unwrap()
	tester.AssertEqual(t, io.EOF, u)
}

func TestError_Error(t *testing.T) {
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

func TestError_Is(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		err    error
		target error
		same   bool
	}{
		"nil target": {
			err:    nil,
			target: nil,
			same:   true,
		},
		"nil pointer target": {
			err:    nil,
			target: (*errkit.Error)(nil),
			same:   false,
		},
		"nil pointer err": {
			err:    (*errkit.Error)(nil),
			target: nil,
			same:   false,
		},
		"nil pointers": {
			err:    (*errkit.Error)(nil),
			target: (*errkit.Error)(nil),
			same:   true,
		},
		"nil err not same": {
			err:    nil,
			target: &errkit.Error{Message: "m", Detail: "d"},
			same:   false,
		},
		"nil target not same": {
			err:    &errkit.Error{Message: "m", Detail: "d"},
			target: nil,
			same:   false,
		},
		"eof err not same": {
			err:    io.EOF,
			target: &errkit.Error{Message: "m", Detail: "d"},
			same:   false,
		},
		"eof target not same": {
			err:    &errkit.Error{Message: "m", Detail: "d"},
			target: io.EOF,
			same:   false,
		},
		"equal": {
			err:    &errkit.Error{Message: "m", Detail: "d"},
			target: &errkit.Error{Message: "m", Detail: "d"},
			same:   true,
		},
		"message mismatch": {
			err:    &errkit.Error{Message: "m", Detail: "d"},
			target: &errkit.Error{Message: "M", Detail: "d"},
			same:   false,
		},
		"detail mismatch": {
			target: &errkit.Error{Message: "m", Detail: "d"},
			err:    &errkit.Error{Message: "m", Detail: "D"},
			same:   true,
		},
		"same after unwrap error": {
			err:    fmt.Errorf("outer error [%w]", &errkit.Error{Message: "m", Detail: "d"}),
			target: &errkit.Error{Message: "m", Detail: "d"},
			same:   true,
		},
		"same after unwrap errors": {
			err:    fmt.Errorf("outer error [%w] [%w]", io.EOF, &errkit.Error{Message: "m", Detail: "d"}),
			target: &errkit.Error{Message: "m", Detail: "d"},
			same:   true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tester.AssertEqual(t, tc.same, errors.Is(tc.err, tc.target))
		})
	}
}

func TestError_Map(t *testing.T) {
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

func TestError_SlogAttrs(t *testing.T) {
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
