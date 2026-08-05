package errkit

import (
	"runtime"
	"strconv"
	"strings"
)

// Frame holds stack frame location.
// See also [runtime.Frame].
type Frame struct {
	// Pkg is go package name of the caller.
	Pkg string `json:"pkg" msgpack:"pkg" toml:"pkg" xml:"pkg" yaml:"pkg"`
	// File is the file name of the caller.
	File string `json:"file" msgpack:"file" toml:"file" xml:"file" yaml:"file"`
	// Func is the function name of the caller.
	Func string `json:"func" msgpack:"func" toml:"func" xml:"func" yaml:"func"`
	// Line is the line number of the caller.
	Line int `json:"line" msgpack:"line" toml:"line" xml:"line" yaml:"line"`
}

func (f *Frame) String() string {
	var builder strings.Builder
	builder.Grow(80)
	builder.WriteString("Pkg:")
	builder.WriteString(f.Pkg)
	builder.WriteString(" File:")
	builder.WriteString(f.File)
	builder.WriteString(" Func:")
	builder.WriteString(f.Func)
	builder.WriteString(" Line:")
	builder.WriteString(strconv.Itoa(f.Line))
	return builder.String()
}

// callerFrame returns single caller frame.
func callerFrame(skip int) Frame {
	pc := make([]uintptr, 1)
	n := runtime.Callers(skip+2, pc)
	if n < 1 {
		return Frame{} // No frame to report.
	}
	f, _ := runtime.CallersFrames(pc).Next()
	return frame(f)
}

// callerFrames returns a slice of caller frames.
func callerFrames(skip int) []Frame {
	pcs := make([]uintptr, 64) // Max 64 frames.
	n := runtime.Callers(skip+2, pcs)
	if n < 1 {
		return nil // No frames to report.
	}
	frames := runtime.CallersFrames(pcs[:n])
	fs := make([]Frame, n)
	for i := range n {
		f, _ := frames.Next()
		fs[i] = frame(f)
	}
	return fs
}

func frame(f runtime.Frame) Frame {
	file := f.File
	pkg, fn, pkgfn := "", f.Function, f.Function // pkgfn is "<Pkg>.<Func>"
	slash := max(0, strings.LastIndexByte(pkgfn, '/'))
	if dot := strings.IndexByte(pkgfn[slash:], '.'); dot >= 0 {
		pkg, fn = pkgfn[:slash+dot], pkgfn[slash+dot+1:]
		file = strings.TrimPrefix(strings.TrimPrefix(file, pkg), "/")
	}
	return Frame{
		Pkg:  pkg,
		File: file,
		Line: f.Line,
		Func: fn,
	}
}
