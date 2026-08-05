# go-errkit

**Go library for errors and exceptions.**

<div align="center">

[![GoDoc](https://godoc.org/github.com/aileron-projects/go-errkit?status.svg)](http://godoc.org/github.com/aileron-projects/go-errkit)
[![Test](https://github.com/aileron-projects/go-errkit/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/aileron-projects/go-errkit/actions/workflows/test.yaml?query=branch%3Amain)
[![License](https://img.shields.io/badge/License-Apache%202.0-yellow.svg)](./LICENSE)

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/aileron-projects/go-errkit)
[![OpenSourceInsight](https://badgen.net/badge/open%2Fsource%2F/insight/cyan)](https://deps.dev/go/github.com%2Faileron-projects%2Fgo-errkit)
[![OSS Insight](https://badgen.net/badge/OSS/Insight/orange)](https://ossinsight.io/analyze/aileron-projects/go-errkit)

</div>

## Features

- Simple and structural error (`Error`).
- Error definition and instanciation features(`ErrDefinition`, `ErrInstance`).
- SLogger integration.
- Easy stack frames manipulation.
- Application-wide error handlers.

## Tested Environments

Operating System:

- `Linux`: [ubuntu-latest](https://github.com/actions/runner-images)
- `Windows`: [windows-latest](https://github.com/actions/runner-images)
- `macOS`: [macos-latest](https://github.com/actions/runner-images)

Architecture (Using QEMU on linux):

- x86: `amd64`, `386`
- arm: `arm/v5`, `arm/v6`, `arm/v7`, `arm64`
- risc: `riscv64`, `loong64`
- ppc: `ppc64`, `ppc64le`
- mips: `mips`, `mips64`, `mips64le`, `mipsle`
- ibm: `s390x`

## Release Cycle

- Releases are made as needed.
- [Semantic Versioning](https://semver.org/) `vX.Y.Z` is used.

## License

[Apache-2.0](LICENSE)

## Usage

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

## Build Tags

No build tags defined for this library.

## Enviromental Variables

- `GO_ERRKIT_TRACE_ENABLED`: enables ErrDefinition to output tracing info of instanciation.
