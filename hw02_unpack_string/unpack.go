package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var sb strings.Builder
	var prevLetter bool
	var prev rune
	for i, v := range str {
		switch {
		case unicode.IsLetter(v):
			if prevLetter {
				sb.WriteRune(prev)
			}
			prevLetter = true
		case unicode.IsDigit(v):
			if i == 0 || !prevLetter {
				return "", ErrInvalidString
			}
			prevLetter = false
			count, _ := strconv.Atoi(string(v))
			if count == 0 {
				prev = v
				continue
			} else {
				sb.WriteString(strings.Repeat(string(prev), count))
			}
		}
		prev = v
	}
	if prevLetter {
		sb.WriteRune(prev)
	}
	return sb.String(), nil
}
