//go:build !dev

package main

import (
	"io"
	"os"
)

func logWriter() io.Writer {
	filepath, _ := getConfigFilePath("crispy", "log")
	f, _ := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return f
}
