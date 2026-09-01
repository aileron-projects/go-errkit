<!-- markdownlint-disable MD033 MD041 -->

<div align="center">

[![Release](https://img.shields.io/github/v/release/aileron-projects/go-errkit?sort=semver)](https://github.com/aileron-projects/go-errkit/releases)
[![Reference](https://pkg.go.dev/badge/github.com/aileron-projects/go-errkit.svg)](https://pkg.go.dev/github.com/aileron-projects/go-errkit)
[![DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/aileron-projects/go-errkit)
[![Test](https://github.com/aileron-projects/go-errkit/actions/workflows/test.yaml/badge.svg)](https://github.com/aileron-projects/go-errkit/actions/workflows/test.yaml)

[![Insights](https://badgen.net/badge/Insights/open%2Fsource%2Finsights/cyan)](https://deps.dev/go/github.com%2Faileron-projects%2Fgo-errkit)
[![Insights](https://badgen.net/badge/Insights/OSS%2FInsight/orange)](https://ossinsight.io/analyze/aileron-projects/go-errkit)

</div>

# go-errkit

**Go library for errors and exceptions.**

## Features

- Simple and structural error (`Error`).
- Error definition and instanciation features(`ErrDefinition`, `ErrInstance`).
- SLogger integration.
- Easy stack frames manipulation.
- Application-wide error handlers.

## Usages

### Error definition and instanciation

Use `NewErrDefinition` to create a new definition.
`New` and `NewStack` creates a new error instance from the defition.

```go
// Define error.
// Definition has code, kind, message, attributes and instance id generation func.
E123 := errkit.NewErrDefinition("E123", "KindXXX", "example error. foo=%s", map[string]string{"type": "server"}, instanceID)

// Instanciate an error from definition.
// "FOO" will be the args for message. 
//   i.e. fmt.Sprintf("example error. foo=%s", "FOO)
ins123 := E123.New(nil, "FOO")
```

The `ins123` can be used for logging.

```go
fmt.Println(ins123)
// E123 KindXXX :example error. foo=FOO (type=server)

fmt.Println(ins123.Error())
// E123 KindXXX :example error. foo=FOO (type=server)

lgJSON := slog.New(slog.NewJSONHandler(os.Stdout, nil))
lgJSON.InfoContext(context.Background(), "JSON logger.",ins123.SlogAttr())
// {"time":"2026-08-06T08:31:03.6992616+09:00","level":"INFO","msg":"JSON logger.","error":{"code":"E123","kind":"KindXXX","message":"example error. foo=FOO","attrs":{"type":"server"}}}

lgText := slog.New(slog.NewTextHandler(os.Stdout, nil))
lgText.InfoContext(context.Background(), "Text logger.", ins123.SlogAttr())
// time=2026-08-06T08:31:03.732+09:00 level=INFO msg="Text logger." error.code=E123 error.kind=KindXXX error.message="example error. foo=FOO" error.attrs.type=server
```

### Enviromental Variables

- `GO_ERRKIT_TRACE_ENABLED`: enables ErrDefinition to output tracing info of instanciation.

## Docs & Examples

- GoDoc: <https://pkg.go.dev/github.com/aileron-projects/go-errkit>
- Examples:
  - [example_test.go](./example_test.go)
  - Error definitions and instanciation: [examples/instanciation/](./examples/instanciation/)

## References
