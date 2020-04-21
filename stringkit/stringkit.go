package stringkit

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Reverse reverses a string.
func Reverse(s string) string {
	rr := []rune(s)
	for from, to := 0, len(rr)-1; from < to; from, to = from+1, to-1 {
		rr[from], rr[to] = rr[to], rr[from]
	}
	return string(rr)
}

// Int64 parse string to int64
func Int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// Int parse string to int
func Int(str string) (int, error) {
	return strconv.Atoi(str)
}

// Float64 parse string to float64
func Float64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

// MaskPhone replace last 4 number by ****
func MaskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:7] + "****"
}

// MaskName replace second name by *
func MaskName(name string) string {
	count := utf8.RuneCountInString(name)
	if count == 0 {
		return ""
	}
	prefix, _ := utf8.DecodeRuneInString(name)
	suffix := strings.Repeat("*", count-1)
	return string(prefix) + suffix
}

// FormatAnswer format unstandardized answer
func FormatAnswer(answer string) string {
	trimedSpaceString := strings.TrimSpace(answer)
	space := regexp.MustCompile(`\s+`)
	result := space.ReplaceAllString(trimedSpaceString, " ")
	return result
}
