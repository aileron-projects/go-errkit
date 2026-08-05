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

func TestToMap(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		m := errkit.ToMap(nil)
		tester.AssertDeepEqual(t, nil, m)
	})
	t.Run("primitive error", func(t *testing.T) {
		m := errkit.ToMap(io.EOF)
		want := map[string]any{"message": "EOF"}
		tester.AssertDeepEqual(t, want, m)
	})
	t.Run("wrapped error", func(t *testing.T) {
		err := fmt.Errorf("outer error [%w]", io.EOF)
		m := errkit.ToMap(err)
		want := map[string]any{
			"message": "outer error [EOF]",
			"causes": []map[string]any{
				{"message": "EOF"},
			},
		}
		tester.AssertDeepEqual(t, want, m)
	})
	t.Run("wrapped errors", func(t *testing.T) {
		err := errors.Join(io.EOF, io.ErrUnexpectedEOF)
		m := errkit.ToMap(err)
		want := map[string]any{
			"message": "EOF\nunexpected EOF",
			"causes": []map[string]any{
				{"message": "EOF"},
				{"message": "unexpected EOF"},
			},
		}
		tester.AssertDeepEqual(t, want, m)
	})
	t.Run("interface", func(t *testing.T) {
		def := errkit.NewErrDefinition("c", "k", "m", nil, func() string { return "i" })
		err := def.NewStack(nil)
		m := errkit.ToMap(err)
		tester.AssertEqual(t, "c", m["code"])
		tester.AssertEqual(t, "k", m["kind"])
		tester.AssertEqual(t, "m", m["message"])
		tester.AssertEqual(t, "i", m["instance"])
		tester.AssertEqual(t, true, len(m["frames"].([]string)) > 0)
	})
}

func TestToSlogAttrs(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		m := errkit.ToSlogAttrs(nil)
		tester.AssertDeepEqual(t, nil, m)
	})
	t.Run("primitive error", func(t *testing.T) {
		m := errkit.ToSlogAttrs(io.EOF)
		want := []slog.Attr{slog.String("message", "EOF")}
		tester.AssertDeepEqual(t, want, m)
	})
	t.Run("wrapped error", func(t *testing.T) {
		err := fmt.Errorf("outer error [%w]", io.EOF)
		m := errkit.ToSlogAttrs(err)
		want := []slog.Attr{
			slog.String("message", "outer error [EOF]"),
			slog.GroupAttrs("cause", slog.String("message", "EOF")),
		}
		tester.AssertDeepEqual(t, want, m)
	})
	t.Run("interface", func(t *testing.T) {
		def := errkit.NewErrDefinition("c", "k", "m", nil, func() string { return "i" })
		err := def.New(nil)
		m := errkit.ToSlogAttrs(err)
		want := []slog.Attr{
			slog.String("code", "c"),
			slog.String("kind", "k"),
			slog.String("message", "m"),
			slog.String("instance", "i"),
		}
		tester.AssertDeepEqual(t, want, m)
	})
}
