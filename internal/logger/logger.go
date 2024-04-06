package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/gookit/color"
)

type LogType uint8

const (
	// Prints with default color
	LogNormal LogType = iota
	LogSuccess
	LogWarning
	LogError
)

// Terminal color renderers
var (
	renderTime    = color.Gray.Render
	renderSuccess = color.Green.Render
	renderWarning = color.Yellow.Render
	renderError   = color.Red.Render
)

var LogNoColor = false

func Printf(logType LogType, format string, v ...any) {
	rendered := fmt.Sprintf(format, v...)

	if !LogNoColor {
		switch logType {
		case LogSuccess:
			rendered = renderSuccess(rendered)
		case LogWarning:
			rendered = renderWarning(rendered)
		case LogError:
			rendered = renderError(rendered)
		}
	}

	timePrefix := fmt.Sprintf("[%s] ", time.Now().Format("Mon Jan _2 15:04:05 2006"))

	if !LogNoColor {
		timePrefix = renderTime(timePrefix)
	}

	rendered = timePrefix + rendered
	print(rendered)
}

func Fatalf(format string, v ...any) {
	Printf(LogError, format, v...)
	os.Exit(1)
}
