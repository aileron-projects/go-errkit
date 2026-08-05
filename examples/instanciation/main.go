package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"

	"github.com/aileron-projects/go-errkit"
)

// Define error.
var E123 = errkit.NewErrDefinition("E123", "KindXXX", "example error. foo=%s", map[string]string{"type": "server"}, instanceID)

func instanceID() string {
	id := rand.Uint32N(99999)
	return fmt.Sprintf("%5d", id)
}

func main() {
	ins123 := E123.New(nil, "FOO")

	// Println
	fmt.Println(ins123)
	fmt.Println(ins123.Error())

	// JSON log
	lgJSON := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	lgJSON.InfoContext(context.Background(), "JSON logger.", "error", ins123.SlogAttrs())

	// Text log
	lgText := slog.New(slog.NewTextHandler(os.Stdout, nil))
	lgText.InfoContext(context.Background(), "Text logger.", "error", ins123.SlogAttrs())
}
