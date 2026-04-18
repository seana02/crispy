//go:build dev

package main

import (
	"io"
	"os"
)

func logWriter() io.Writer {
	return os.Stdout
}
