package logger

import "fmt"

type Color uint8

const (
	Red Color = iota + 31
	Green
	Yellow
	Gray Color = 90
)

func (c Color) Render(text string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, text)
}
