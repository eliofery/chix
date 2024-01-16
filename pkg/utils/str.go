package utils

import "unicode"

// FirstToUpper переводит первый символ в верхний регистр
func FirstToUpper(s string) string {
	str := []rune(s)

	return string(append([]rune{unicode.ToUpper(str[0])}, str[1:]...))
}
