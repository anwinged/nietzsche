package nietzsche

import "strings"

func Shorten(text string, length int) string {
	noNewLine := strings.Replace(text, "\n", "", -1)
	l := len(noNewLine)
	if l > length {
		return noNewLine[:length] + "..."
	}
	return noNewLine[:l]
}
