package util

var (
	alphanum [256]bool
)

func init() {
	for _, b := range "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		alphanum[b] = true
	}
}

func IsAlphaNum(b byte) bool {
	return alphanum[b]
}
