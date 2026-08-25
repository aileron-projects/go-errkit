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

func TestErrDefinition_Instanciated(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		def    errkit.ErrDefinition
		target error
		same   bool
	}{
		"nil": {
			def:    errkit.ErrDefinition{Code: "c", Kind: "k"},
			target: nil,
			same:   false,
		},
		"not equal": {
			def:    errkit.ErrDefinition{Code: "c", Kind: "k"},
			target: io.EOF,
			same:   false,
		},
		"code mismatch": {
			def:    errkit.ErrDefinition{Code: "c", Kind: "k"},
			target: &errkit.ErrInstance{Code: "C", Kind: "k"},
			same:   false,
		},
		"kind mismatch": {
			def:    errkit.ErrDefinition{Code: "c", Kind: "k"},
			target: &errkit.ErrInstance{Code: "c", Kind: "K"},
			same:   false,
		},
		"message mismatch": {
			def:    errkit.ErrDefinition{Code: "c", Kind: "k", Message: "m"},
			target: &errkit.ErrInstance{Code: "c", Kind: "k", Message: "M"},
			same:   true,
		},
		"instance mismatch": {
			def:    errkit.ErrDefinition{Code: "c", Kind: "k", Instance: func() string { return "i" }},
			target: &errkit.ErrInstance{Code: "c", Kind: "k", Instance: "I"},
			same:   true,
		},
		"attrs mismatch": {
			def:    errkit.ErrDefinition{Code: "c", Kind: "k", Attrs: map[string]string{"foo": "bar"}},
			target: &errkit.ErrInstance{Code: "c", Kind: "k", Attrs: map[string]string{"FOO": "Bar"}},
			same:   true,
		},
		"same after unwrap": {
			def:    errkit.ErrDefinition{Code: "c", Kind: "k"},
			target: fmt.Errorf("outer error [%w]", &errkit.ErrInstance{Code: "c", Kind: "k"}),
			same:   true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			is := tc.def.Instanciated(tc.target)
			tester.AssertEqual(t, tc.same, is)
		})
	}
}

func TestErrDefinition_New(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		instance *errkit.ErrInstance
		want     errkit.ErrInstance
	}{
		"zero value": {
			instance: (&errkit.ErrDefinition{}).New(nil),
			want:     errkit.ErrInstance{},
		},
		"non inner error": {
			instance: (&errkit.ErrDefinition{Code: "c", Kind: "k", Message: "m"}).New(nil),
			want:     errkit.ErrInstance{Code: "c", Kind: "k", Message: "m"},
		},
		"inner error": {
			instance: (&errkit.ErrDefinition{Code: "c", Kind: "k", Message: "m"}).New(io.EOF),
			want:     errkit.ErrInstance{Code: "c", Kind: "k", Message: "m", Cause: io.EOF},
		},
		"attrs": {
			instance: (&errkit.ErrDefinition{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}}).New(nil),
			want:     errkit.ErrInstance{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}},
		},
		"format message": {
			instance: (&errkit.ErrDefinition{Code: "c", Kind: "k", Message: "foo=%s"}).New(nil, "xxx"),
			want:     errkit.ErrInstance{Code: "c", Kind: "k", Message: "foo=xxx"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tester.AssertEqual(t, tc.want.Code, tc.instance.Code)
			tester.AssertEqual(t, tc.want.Kind, tc.instance.Kind)
			tester.AssertEqual(t, tc.want.Instance, tc.instance.Instance)
			tester.AssertEqual(t, tc.want.Message, tc.instance.Message)
			tester.AssertDeepEqual(t, tc.want.Attrs, tc.instance.Attrs)
			tester.AssertEqual(t, 0, len(tc.instance.Frames))
			tester.AssertEqual(t, tc.want.Cause, tc.instance.Cause)
		})
	}
}

func TestErrDefinition_NewStack(t *testing.T) {
	t.Parallel()
	t.Run("zero value", func(t *testing.T) {
		var ed errkit.ErrDefinition
		e := ed.NewStack(nil)
		w := errkit.ErrInstance{}
		tester.AssertEqual(t, w.Code, e.Code)
		tester.AssertEqual(t, w.Kind, e.Kind)
		tester.AssertEqual(t, w.Message, e.Message)
		tester.AssertDeepEqual(t, w.Attrs, e.Attrs)
		tester.AssertEqual(t, true, len(e.Frames) > 0)
		tester.AssertEqual(t, nil, e.Cause)
	})
	t.Run("cause without stack", func(t *testing.T) {
		e := (&errkit.ErrDefinition{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}}).NewStack(io.EOF)
		w := errkit.ErrInstance{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}}
		tester.AssertEqual(t, w.Code, e.Code)
		tester.AssertEqual(t, w.Kind, e.Kind)
		tester.AssertEqual(t, w.Message, e.Message)
		tester.AssertDeepEqual(t, w.Attrs, e.Attrs)
		tester.AssertEqual(t, true, len(e.Frames) > 0)
		tester.AssertEqual(t, io.EOF, e.Cause)
	})
	t.Run("inner error with stack", func(t *testing.T) {
		inner := &errkit.ErrInstance{Frames: []errkit.Frame{{}, {}}}
		e := (&errkit.ErrDefinition{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}}).NewStack(inner)
		w := errkit.ErrInstance{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}}
		tester.AssertEqual(t, w.Code, e.Code)
		tester.AssertEqual(t, w.Kind, e.Kind)
		tester.AssertEqual(t, w.Message, e.Message)
		tester.AssertDeepEqual(t, w.Attrs, e.Attrs)
		tester.AssertEqual(t, 0, len(e.Frames))
		tester.AssertEqual(t, error(inner), e.Cause)
	})
}

func TestErrInstance_Unwrap(t *testing.T) {
	t.Parallel()
	e := &errkit.ErrInstance{Cause: io.EOF}
	u := e.Unwrap()
	tester.AssertEqual(t, io.EOF, u)
}

func TestErrInstance_Error(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		err  *errkit.ErrInstance
		want string
	}{
		"code kind": {
			err:  &errkit.ErrInstance{Code: "c", Kind: "k"},
			want: "c k :",
		},
		"code kind instance": {
			err:  &errkit.ErrInstance{Code: "c", Kind: "k", Instance: "i"},
			want: "c k i :",
		},
		"code kind message": {
			err:  &errkit.ErrInstance{Code: "c", Kind: "k", Message: "m"},
			want: "c k :m",
		},
		"code kind instance message": {
			err:  &errkit.ErrInstance{Code: "c", Kind: "k", Instance: "i", Message: "m"},
			want: "c k i :m",
		},
		"attrs": {
			err:  &errkit.ErrInstance{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}},
			want: "c k :m (foo=bar)",
		},
		"cause": {
			err:  &errkit.ErrInstance{Code: "c", Kind: "k", Message: "m", Cause: io.EOF},
			want: "c k :m [EOF]",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.err.Error()
			tester.AssertEqual(t, tc.want, got)
		})
	}
}

func TestErrInstance_Is(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		err    error
		target error
		same   bool
	}{
		"nil": {
			err:    nil,
			target: nil,
			same:   true,
		},
		"nil pointer err": {
			err:    (*errkit.ErrInstance)(nil),
			target: nil,
			same:   false,
		},
		"nil pointer target": {
			err:    nil,
			target: (*errkit.ErrInstance)(nil),
			same:   false,
		},
		"nil pointers": {
			err:    (*errkit.ErrInstance)(nil),
			target: (*errkit.ErrInstance)(nil),
			same:   true,
		},
		"nil target": {
			err:    &errkit.ErrInstance{Code: "c", Kind: "k"},
			target: nil,
			same:   false,
		},
		"equal": {
			err:    &errkit.ErrInstance{Code: "c", Kind: "k"},
			target: &errkit.ErrInstance{Code: "c", Kind: "k"},
			same:   true,
		},
		"not equal": {
			err:    &errkit.ErrInstance{Code: "c", Kind: "k"},
			target: io.EOF,
			same:   false,
		},
		"code mismatch": {
			err:    &errkit.ErrInstance{Code: "c", Kind: "k"},
			target: &errkit.ErrInstance{Code: "C", Kind: "k"},
			same:   false,
		},
		"kind mismatch": {
			err:    &errkit.ErrInstance{Code: "c", Kind: "k"},
			target: &errkit.ErrInstance{Code: "c", Kind: "K"},
			same:   false,
		},
		"instance mismatch": {
			err:    &errkit.ErrInstance{Code: "c", Kind: "k", Instance: "i"},
			target: &errkit.ErrInstance{Code: "c", Kind: "k", Instance: "I"},
			same:   true,
		},
		"message mismatch": {
			err:    &errkit.ErrInstance{Code: "c", Kind: "k", Message: "m"},
			target: &errkit.ErrInstance{Code: "c", Kind: "k", Message: "M"},
			same:   true,
		},
		"attrs mismatch": {
			err:    &errkit.ErrInstance{Code: "c", Kind: "k", Attrs: map[string]string{"foo": "bar"}},
			target: &errkit.ErrInstance{Code: "c", Kind: "k", Attrs: map[string]string{"FOO": "Bar"}},
			same:   true,
		},
		"same after unwrap error": {
			err:    fmt.Errorf("outer error [%w]", &errkit.ErrInstance{Code: "c", Kind: "k"}),
			target: &errkit.ErrInstance{Code: "c", Kind: "k"},
			same:   true,
		},
		"same after unwrap errors": {
			err:    fmt.Errorf("outer error [%w] [%w]", io.EOF, &errkit.ErrInstance{Code: "c", Kind: "k"}),
			target: &errkit.ErrInstance{Code: "c", Kind: "k"},
			same:   true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tester.AssertEqual(t, tc.same, errors.Is(tc.err, tc.target))
		})
	}
}

func TestErrInstance_Map(t *testing.T) {
	t.Parallel()
	t.Run("minimum", func(t *testing.T) {
		e := &errkit.ErrInstance{Code: "c", Kind: "k", Message: "m"}
		got := e.Map()
		want := map[string]any{
			"code":    "c",
			"kind":    "k",
			"message": "m",
		}
		tester.AssertDeepEqual(t, want, got)
	})
	t.Run("cause", func(t *testing.T) {
		e := &errkit.ErrInstance{Cause: io.EOF}
		got := e.Map()
		want := map[string]any{
			"message": "EOF",
		}
		tester.AssertDeepEqual(t, want, got["cause"].(map[string]any))
	})
	t.Run("attrs", func(t *testing.T) {
		e := &errkit.ErrInstance{Attrs: map[string]string{"foo": "bar"}}
		got := e.Map()
		want := map[string]string{
			"foo": "bar",
		}
		tester.AssertDeepEqual(t, want, got["attrs"].(map[string]string))
	})
	t.Run("instance", func(t *testing.T) {
		e := &errkit.ErrInstance{Instance: "ABC"}
		got := e.Map()
		want := map[string]any{
			"code":     "",
			"kind":     "",
			"message":  "",
			"instance": "ABC",
		}
		tester.AssertDeepEqual(t, want, got)
	})
	t.Run("frames", func(t *testing.T) {
		def := errkit.NewErrDefinition("c", "k", "m", nil, nil)
		e := def.NewStack(nil)
		got := e.Map()
		tester.AssertEqual(t, true, len(got["frames"].([]string)) > 0)
	})
}

func TestErrInstance_SlogAttrs(t *testing.T) {
	t.Parallel()
	t.Run("minimum", func(t *testing.T) {
		e := &errkit.ErrInstance{Code: "c", Kind: "k", Message: "m"}
		got := e.SlogAttrs()
		want := []slog.Attr{
			slog.String("code", "c"),
			slog.String("kind", "k"),
			slog.String("message", "m"),
		}
		tester.AssertDeepEqual(t, want, got)
	})
	t.Run("cause", func(t *testing.T) {
		e := &errkit.ErrInstance{Cause: io.EOF}
		got := e.SlogAttrs()
		want := []slog.Attr{
			slog.String("code", ""),
			slog.String("kind", ""),
			slog.String("message", ""),
			slog.GroupAttrs("cause", slog.String("message", "EOF")),
		}
		tester.AssertDeepEqual(t, want, got)
	})
	t.Run("attrs", func(t *testing.T) {
		e := &errkit.ErrInstance{Attrs: map[string]string{"foo": "bar"}}
		got := e.SlogAttrs()
		want := []slog.Attr{
			slog.String("code", ""),
			slog.String("kind", ""),
			slog.String("message", ""),
			slog.GroupAttrs("attrs", slog.String("foo", "bar")),
		}
		tester.AssertDeepEqual(t, want, got)
	})
	t.Run("instance", func(t *testing.T) {
		e := &errkit.ErrInstance{Instance: "ABC"}
		got := e.SlogAttrs()
		want := []slog.Attr{
			slog.String("code", ""),
			slog.String("kind", ""),
			slog.String("message", ""),
			slog.String("instance", "ABC"),
		}
		tester.AssertDeepEqual(t, want, got)
	})
	t.Run("frames", func(t *testing.T) {
		def := errkit.NewErrDefinition("c", "k", "m", nil, nil)
		e := def.NewStack(nil)
		got := e.SlogAttrs()
		for _, a := range got {
			if a.Key == "frames" {
				tester.AssertEqual(t, true, len(a.Value.Any().([]string)) > 0)
				return
			}
		}
		t.Error("frame does not exist.")
	})
}
