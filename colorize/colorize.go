package colorize

import (
	"strconv"
)

func colorize(s string, r int, g int, b int) string {
	red := strconv.Itoa(r)
	green := strconv.Itoa(g)
	blue := strconv.Itoa(b)

	return "\x1b[38;2;" + red + ";" + green + ";" + blue + ";10m" + s + "\x1b[0m"
}

func Error(s string) string {
	return colorize(s, 255, 0, 0)
}

func Warn(s string) string {
	return colorize(s, 255, 255, 0)
}

func Log(s string) string {
	return s
}

