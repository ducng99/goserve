package logger

import (
	"fmt"
	"os"
	"time"
)

type LogType uint8

const (
	// Prints with default color
	LogNormal LogType = iota
	LogSuccess
	LogWarn
	LogError
)

// Terminal color renderers
var (
	renderTime    = Gray.Render
	renderSuccess = Green.Render
	renderWarning = Yellow.Render
	renderError   = Red.Render
)

var LogWithColor = true

func Printf(logType LogType, format string, v ...any) {
	rendered := fmt.Sprintf(format, v...)

	if LogWithColor {
		switch logType {
		case LogSuccess:
			rendered = renderSuccess(rendered)
		case LogWarn:
			rendered = renderWarning(rendered)
		case LogError:
			rendered = renderError(rendered)
		}
	}

	timePrefix := fmt.Sprintf("[%s] ", time.Now().Format("Mon Jan _2 15:04:05 2006"))

	if LogWithColor {
		timePrefix = renderTime(timePrefix)
	}

	rendered = timePrefix + rendered
	print(rendered)
}

func Fatalf(format string, v ...any) {
	Printf(LogError, format, v...)
	os.Exit(1)
}
