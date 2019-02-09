package noarch

import "strings"

type WcharT = rune

func Wcscmp(w1, w2 []WcharT) int {
	if len(w1) > len(w2) {
		return -1
	}
	if len(w1) < len(w2) {
		return 1
	}
	for i := range w1 {
		if r := strings.Compare(string(w1[i]), string(w2[i])); r != 0 {
			return r
		}
	}
	return 0
}
