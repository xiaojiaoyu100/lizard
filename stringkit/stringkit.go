package stringkit

// Reverse reverses a string.
func Reverse(s string) string {
	rr := []rune(s)
	for from, to := 0, len(rr)-1; from < to; from, to = from+1, to-1 {
		rr[from], rr[to] = rr[to], rr[from]
	}
	return string(rr)
}
