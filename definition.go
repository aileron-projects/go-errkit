package errkit

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"strconv"
	"strings"
)

// NewErrDefinition returns a new instance of ErrDefinition.
// See [ErrDefinition].
func NewErrDefinition(code, kind, message string, attrs map[string]string, instance func() string) *ErrDefinition {
	return &ErrDefinition{
		Code:     code,
		Kind:     kind,
		Message:  message,
		Attrs:    maps.Clone(attrs),
		Instance: instance,
	}
}

// ErrDefinition is the error definition.
// Error instances are created from the definition.
type ErrDefinition struct {
	// Code is the error code.
	Code string
	// Kind is the error kind of this error.
	Kind string
	// Message is the error message.
	// Message will be formatted by [fmt.Sprintf].
	Message string
	// Attrs is the list of attributions.
	// Attrs will be copied to all instances.
	Attrs map[string]string
	// Instance returns an instance identifier.
	Instance func() string
}

// Instanciated returns if this definition instanciated the err.
// The err will be unwrapped if possible.
func (d *ErrDefinition) Instanciated(err error) bool {
	for err != nil {
		e, ok := err.(*ErrInstance)
		if ok {
			return d.Code == e.Code && d.Kind == e.Kind
		}
		err = errors.Unwrap(err)
	}
	return false
}

// New returns a new error instance from the definition.
// Use [ErrDefinition.NewStack] when stack frames are necessary.
func (d *ErrDefinition) New(cause error, values ...any) *ErrInstance {
	err := &ErrInstance{
		Cause:   cause,
		Code:    d.Code,
		Kind:    d.Kind,
		Message: fmt.Sprintf(d.Message, values...),
		Attrs:   maps.Clone(d.Attrs),
	}
	if d.Instance != nil {
		err.Instance = d.Instance()
	}
	traceTo(os.Stdout, err) // Trace instanciation.
	return err
}

// NewStack returns a new error instance from the definition.
// Use [ErrDefinition.New] when stack frames are not necessary.
func (d *ErrDefinition) NewStack(cause error, values ...any) *ErrInstance {
	err := &ErrInstance{
		Cause:   cause,
		Code:    d.Code,
		Kind:    d.Kind,
		Message: fmt.Sprintf(d.Message, values...),
		Attrs:   maps.Clone(d.Attrs),
	}
	if d.Instance != nil {
		err.Instance = d.Instance()
	}
	e := cause
	for e != nil {
		ee, ok := e.(*ErrInstance)
		if ok && len(ee.Frames) > 0 {
			traceTo(os.Stdout, err) // Trace instanciation.
			return err
		}
		e = errors.Unwrap(e)
	}
	err.Frames = callerFrames(1)
	traceTo(os.Stdout, err) // Trace instanciation.
	return err
}

// Error is the general error type.
type ErrInstance struct {
	// Cause is the error cause, or inner error.
	Cause error `json:"cause,omitempty" msgpack:"cause,omitempty" xml:"cause,omitempty" yaml:"cause,omitempty"`
	// Code is the error code, name or alias for the error.
	// Code is compared in [Errors.Is].
	Code string `json:"code" msgpack:"code" xml:"code" yaml:"code"`
	// Kind is the error kind.
	// Kind is compared in [Errors.Is].
	Kind string `json:"kind" msgpack:"kind" xml:"kind" yaml:"kind"`
	// Instance is the instance identifier.
	// Instance is NOT compared in [Errors.Is].
	Instance string `json:"instance,omitempty" msgpack:"instance,omitempty" xml:"instance,omitempty" yaml:"instance,omitempty"`
	// Message is the error message.
	Message string `json:"message" msgpack:"message" xml:"message" yaml:"message"`
	// Attrs are the attribution, or extra information, to this error.
	Attrs map[string]string `json:"attrs" msgpack:"attrs" xml:"attrs" yaml:"attrs"`
	// Frames is the list of stack trace frames.
	Frames []Frame `json:"frames,omitempty" msgpack:"frames,omitempty" xml:"frames,omitempty" yaml:"frames,omitempty"`
}

// Error implements [error] interface.
func (e *ErrInstance) Error() string {
	var builder strings.Builder
	builder.Grow(100)
	_, _ = builder.WriteString(e.Code + " ")
	_, _ = builder.WriteString(e.Kind + " ")
	if e.Instance != "" {
		_, _ = builder.WriteString(e.Instance + " ")
	}
	_, _ = builder.WriteString(":")
	_, _ = builder.WriteString(e.Message)
	if len(e.Attrs) > 0 {
		kvs := make([]string, 0, len(e.Attrs))
		for k, v := range e.Attrs {
			kvs = append(kvs, k+"="+v)
		}
		_, _ = builder.WriteString(" (" + strings.Join(kvs, ",") + ")")
	}
	if e.Cause != nil {
		_, _ = builder.WriteString(" [")
		_, _ = builder.WriteString(e.Cause.Error())
		_, _ = builder.WriteString("]")
	}
	return builder.String()
}

// Unwrap returns the inner error if any.
func (e *ErrInstance) Unwrap() error {
	return e.Cause
}

// Is returns if this error is identical to the given error.
// This can be used with [errors.Is].
func (e *ErrInstance) Is(err error) bool {
	ee, ok := err.(*ErrInstance)
	if ok {
		return e.Code == ee.Code && e.Kind == ee.Kind
	}
	return false
}

// Map returns error information as map.
func (e *ErrInstance) Map() map[string]any {
	m := map[string]any{
		"code":    e.Code,
		"kind":    e.Kind,
		"message": e.Message,
	}
	if e.Attrs != nil {
		m["attrs"] = maps.Clone(e.Attrs)
	}
	if e.Instance != "" {
		m["instance"] = e.Instance
	}
	if len(e.Frames) > 0 {
		fs := make([]string, 0, len(e.Frames))
		for _, f := range e.Frames {
			fs = append(fs, f.Pkg+":"+f.File+":L"+strconv.Itoa(f.Line)+"("+f.Func+")")
		}
		m["frames"] = fs
	}
	if cause := ToMap(e.Cause); cause != nil {
		m["cause"] = cause
	}
	return m
}

// SlogAttrs returns error information as [slog.Attr].
func (e *ErrInstance) SlogAttrs() []slog.Attr {
	a := []slog.Attr{
		slog.String("code", e.Code),
		slog.String("kind", e.Kind),
		slog.String("message", e.Message),
	}
	if e.Attrs != nil {
		aa := make([]slog.Attr, 0, len(e.Attrs))
		for k, v := range e.Attrs {
			aa = append(aa, slog.String(k, v))
		}
		a = append(a, slog.GroupAttrs("attrs", aa...))
	}
	if e.Instance != "" {
		a = append(a, slog.String("instance", e.Instance))
	}
	if len(e.Frames) > 0 {
		fs := make([]string, 0, len(e.Frames))
		for _, f := range e.Frames {
			fs = append(fs, f.Pkg+":"+f.File+":L"+strconv.Itoa(f.Line)+"("+f.Func+")")
		}
		a = append(a, slog.Any("frames", fs))
	}
	if causes := ToSlogAttrs(e.Cause); len(causes) > 0 {
		a = append(a, slog.GroupAttrs("cause", causes...))
	}
	return a
}

// SlogAttr returns error information as [slog.Attr].
// SlogAttr is the alias for slog.GroupAttrs("error", e.SlogAttrs()...)
func (e *ErrInstance) SlogAttr() slog.Attr {
	return slog.GroupAttrs("error", e.SlogAttrs()...)
}
