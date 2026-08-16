package errkit

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/kr/pretty"
)

// traceEnabled if true, error tracing is enabled mainly for debugging.
var traceEnabled = os.Getenv("GO_ERRKIT_TRACE_ENABLED") == "1"

// traceTo prints error trace information to w.
func traceTo(w io.Writer, e error) {
	if !traceEnabled || e == nil {
		return
	}

	f := callerFrame(2)
	loc := (&f).String() // Frame location.

	var buf bytes.Buffer
	_, _ = buf.WriteString(time.Now().Format(time.DateTime) + " [ERRKIT][TRACE] " + loc + "\n")
	_, _ = pretty.Fprintf(&buf, "%# v", ToMap(e))
	trace := bytes.ReplaceAll(buf.Bytes(), []byte("\n"), []byte("\n  | "))
	trace = append(trace, '\n')
	_, _ = w.Write(trace)
}
