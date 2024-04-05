package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/gookit/color"
)

type LogType uint8

const (
	LogNormal LogType = iota
	LogSuccess
	LogWarning
	LogError
)

var (
	renderSuccess = color.Green.Render
	renderWarning = color.Yellow.Render
	renderError   = color.Red.Render
)

func Printf(logType LogType, format string, v ...any) {
	rendered := fmt.Sprintf(format, v...)

	switch logType {
	case LogSuccess:
		rendered = renderSuccess(rendered)
	case LogWarning:
		rendered = renderWarning(rendered)
	case LogError:
		rendered = renderError(rendered)
	}

	rendered = fmt.Sprintf("[%s] %s", time.Now().Format("Mon Jan _2 15:04:05 2006"), rendered)
	fmt.Print(rendered)
}

func Fatalf(format string, v ...any) {
	Printf(LogError, format, v...)
	os.Exit(1)
}
