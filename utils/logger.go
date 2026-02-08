package utils

import (
	"log"
	"os"

	"github.com/fatih/color"
)

var Logger *log.Logger

func init() {
	Logger := log.New(os.Stdout, "[crispy] ", log.LstdFlags)
	Logger.Println("Logger initialized")
}

func Info(msg string) {
	Logger.Println(msg)
}

func Warn(msg string) {
	color.Set(color.FgYellow)
	Logger.Printf("[WARN] %s", msg)
	color.Unset()
}

func Error(msg string) {
	color.Set(color.FgRed)
	Logger.Printf("[ERROR] %s", msg)
	color.Unset()
}
