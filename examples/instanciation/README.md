# Error instanciation example

## About the example

This example shows how to instanciate an error from error definitions
and how to output errors in log.

## Run the example

```sh
go run ./main.go
```

This shows the error info in various format.

```txt
E123 KindXXX 03145 :example error. foo=FOO (type=server)
E123 KindXXX 03145 :example error. foo=FOO (type=server)
{"time":"2026-08-16T14:50:18.3397+09:00","level":"INFO","msg":"JSON logger.","error":{"code":"E123","kind":"KindXXX","message":"example error. foo=FOO","attrs":{"type":"server"},"instance":"03145"}}
time=2026-08-16T14:50:18.339+09:00 level=INFO msg="Text logger." error.code=E123 error.kind=KindXXX error.message="example error. foo=FOO" error.attrs.type=server error.instance=03145
```

## Tracing instanciation

Instanciation infos can be traced for debugging with environmental variable.

Use `GO_ERRKIT_TRACE_ENABLED=1`.

```sh
export GO_ERRKIT_TRACE_ENABLED=1
go run ./main.go
```

That shows

```txt
2026-08-16 14:50:18 [ERRKIT][TRACE] Pkg:main File:go-errkit/examples/instanciation/main.go Func:main Line:24
  | map[string]interface {}{
  |     "attrs":    map[string]string{"type":"server"},
  |     "code":     "E123",
  |     "instance": "03145",
  |     "kind":     "KindXXX",
  |     "message":  "example error. foo=FOO",
  | }
E123 KindXXX 03145 :example error. foo=FOO (type=server)
E123 KindXXX 03145 :example error. foo=FOO (type=server)
{"time":"2026-08-16T14:50:18.3397+09:00","level":"INFO","msg":"JSON logger.","error":{"code":"E123","kind":"KindXXX","message":"example error. foo=FOO","attrs":{"type":"server"},"instance":"03145"}}
time=2026-08-16T14:50:18.339+09:00 level=INFO msg="Text logger." error.code=E123 error.kind=KindXXX error.message="example error. foo=FOO" error.attrs.type=server error.instance=03145
```
