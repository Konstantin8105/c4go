package noarch

import (
	"strings"
)

type WcharT = rune

func Wcscmp(w1, w2 []WcharT) int {
	var len1, len2 int
	for len1 = range w1 {
		if rune(w1[len1]) == '\x00' {
			break
		}
	}
	for len2 = range w2 {
		if rune(w2[len2]) == '\x00' {
			break
		}
	}
	if len1 > len2 {
		return -1
	}
	if len1 < len2 {
		return 1
	}
	for i := range w1 {
		if i == len1 {
			break
		}
		if r := strings.Compare(string(w1[i]), string(w2[i])); r != 0 {
			return r
		}
	}
	return 0
}

func Wcscpy(w1, w2 []WcharT) []WcharT {
	for i, c := range w2 {
		w1[i] = WcharT(rune(c))
		if rune(w2[i]) == '\x00' {
			break
		}
	}
	return w2
}

func Wcslen(w []WcharT) int {
	return len(w)
}
