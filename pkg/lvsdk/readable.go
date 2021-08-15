package lvsdk

import (
	"fmt"
	"strings"
	"unicode"
)

func Readable(s string) string {
	b := new(strings.Builder)
	for _, c := range s {
		if unicode.IsControl(c) || unicode.IsSpace(c) {
			h := fmt.Sprintf("[%02X]", int(c))
			b.WriteString(h)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}
