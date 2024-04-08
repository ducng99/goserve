package themes

const (
	ThemeBasic  = "basic"
	ThemePretty = "pretty"
)

// Checks whether the input theme exists
func Exists(inputTheme string) bool {
	switch inputTheme {
	case ThemePretty, ThemeBasic:
		return true
	default:
		return false
	}
}
