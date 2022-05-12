package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	input := []rune(str)
	result := make([]rune, 0, len(str)*2)
	var prev rune
	for i, v := range input {
		if unicode.IsLetter(v) {
			result = append(result, v)
		} else if unicode.IsDigit(v) {
			if i == 0 || unicode.IsDigit(prev) {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(v))
			if count == 0 {
				last := len(result) - 1
				result = result[0:last]
			} else if count > 0 {
				for i := 0; i < count-1; i++ {
					result = append(result, prev)
				}
			}
		}
		prev = v
	}
	return string(result), nil
}
