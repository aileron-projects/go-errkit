package errkit

import (
	"fmt"
	"log/slog"
	"strings"
)

// NewError returns a new instance of [Error].
func NewError(cause error, message, detail string, a ...any) *Error {
	return &Error{
		Cause:   cause,
		Message: message,
		Detail:  fmt.Sprintf(detail, a...),
	}
}

// Error is the basic error struct.
// Use [NewError] to create an instance.
type Error struct {
	// Cause is the cause of this error, or inner error.
	Cause error
	// Message is the error message.
	// Message is compared in the [Error.Is].
	Message string
	// Detail is the error detail.
	// Detail is NOT compared in the [Error.Is].
	Detail string
}

// Error implements [error].
func (e *Error) Error() string {
	var b strings.Builder
	b.Grow(len(e.Message) + len(e.Detail) + 1)
	_, _ = b.WriteString(e.Message)
	if e.Detail != "" {
		_, _ = b.WriteString(" ")
		_, _ = b.WriteString(e.Detail)
	}
	if e.Cause != nil {
		_, _ = b.WriteString(" [")
		_, _ = b.WriteString(e.Cause.Error())
		_, _ = b.WriteString("]")
	}
	return b.String()
}

// Unwrap returns the inner error if any.
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is returns if this error is identical to the given error.
// This can be used with [errors.Is].
func (e *Error) Is(target error) bool {
	ee, ok := target.(*Error)
	if ok {
		return e.Message == ee.Message
	}
	return false
}

// Map returns error information in map.
func (e *Error) Map() map[string]any {
	m := map[string]any{
		"message": e.Message,
		"detail":  e.Detail,
	}
	if cause := ToMap(e.Cause); cause != nil {
		m["cause"] = cause
	}
	return m
}

// SlogAttrs returns error information as [slog.Attr].
func (e *Error) SlogAttrs() []slog.Attr {
	a := []slog.Attr{
		slog.String("message", e.Message),
		slog.String("detail", e.Detail),
	}
	if causes := ToSlogAttrs(e.Cause); len(causes) > 0 {
		a = append(a, slog.GroupAttrs("cause", causes...))
	}
	return a
}

// SlogAttr returns error information as [slog.Attr].
// SlogAttr is the alias for slog.GroupAttrs("error", e.SlogAttrs()...)
func (e *Error) SlogAttr() slog.Attr {
	return slog.GroupAttrs("error", e.SlogAttrs()...)
}
