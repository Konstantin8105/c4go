package util

import (
	"regexp"
	"sync"
)

// cachedRegex - structure for saving regexp`s
type cachedRegex struct {
	sync.RWMutex
	m map[string]*regexp.Regexp
}

// Global variable
var cr = cachedRegex{m: map[string]*regexp.Regexp{}}

// GetRegex return regexp
// added for minimaze regexp compilation
func GetRegex(rx string) *regexp.Regexp {
	cr.RLock()
	v, ok := cr.m[rx]
	cr.RUnlock()
	if ok {
		return v
	}
	// if regexp is not in map
	cr.Lock()
	cr.m[rx] = regexp.MustCompile(rx)
	cr.Unlock()
	return GetRegex(rx)
}
