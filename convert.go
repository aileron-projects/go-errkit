package errkit

import (
	"log/slog"
	"strconv"
)

// ToMap converts error into map.
// If the err implements `Map() map[string]any`, it calls Map().
// ToMap returns nil when the err is nil.
func ToMap(err error) map[string]any {
	if err == nil {
		return nil
	}
	if a, ok := err.(interface{ Map() map[string]any }); ok {
		return a.Map()
	}
	m := map[string]any{"message": err.Error()}
	if e := UnwrapErr(err); e != nil {
		m["causes"] = []map[string]any{ToMap(e)}
		return m
	}
	if errs := UnwrapErrs(err); len(errs) > 0 {
		s := make([]map[string]any, len(errs))
		for i, e := range errs {
			s[i] = ToMap(e)
		}
		m["causes"] = s
	}
	return m
}

// ToSlogAttrs converts error into [slog.Attr].
// If the err implements `SlogAttrs() []slog.Attr`, it calls SlogAttrs().
// ToSlogAttrs returns nil when the err is nil.
func ToSlogAttrs(err error) []slog.Attr {
	if err == nil {
		return nil
	}
	if a, ok := err.(interface{ SlogAttrs() []slog.Attr }); ok {
		return a.SlogAttrs()
	}
	s := []slog.Attr{slog.String("message", err.Error())}
	if e := UnwrapErr(err); e != nil {
		s = append(s, slog.GroupAttrs("cause", ToSlogAttrs(e)...))
		return s
	}
	if errs := UnwrapErrs(err); len(errs) > 0 {
		for i, e := range errs {
			s = append(s, slog.GroupAttrs("cause."+strconv.Itoa(i+1), ToSlogAttrs(e)...))
		}
	}
	return s
}

// ToSlogAttr converts error into [slog.Attr].
// ToSlogAttr is the alias for slog.GroupAttrs("error", ToSlogAttrs(err)...)
func ToSlogAttr(err error) slog.Attr {
	return slog.GroupAttrs("error", ToSlogAttrs(err)...)
}
