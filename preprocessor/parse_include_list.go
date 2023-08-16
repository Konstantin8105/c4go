package preprocessor

import (
	"fmt"
	"strings"
)

// parseIncludeList - parse list of includes
// Example :
//
//	exit.o: exit.c /usr/include/stdlib.h /usr/include/features.h \
//	   /usr/include/stdc-predef.h /usr/include/x86_64-linux-gnu/sys/cdefs.h
func parseIncludeList(line string) (lines []string, err error) {
	line = strings.Replace(line, "\t", " ", -1)
	line = strings.Replace(line, "\r", " ", -1) // Added for Mac endline symbol
	line = strings.Replace(line, "\xFF", " ", -1)
	line = strings.Replace(line, "\u0100", " ", -1)

	sepLines := strings.Split(line, "\n")
	var indexes []int
	for i := range sepLines {
		if !(strings.Index(sepLines[i], ":") < 0) {
			indexes = append(indexes, i)
		}
	}
	if len(indexes) > 1 {
		for i := 0; i < len(indexes); i++ {
			var partLines []string
			var block string
			if i+1 == len(indexes) {
				block = strings.Join(sepLines[indexes[i]:], "\n")
			} else {
				block = strings.Join(sepLines[indexes[i]:indexes[i+1]], "\n")
			}
			partLines, err = parseIncludeList(block)
			if err != nil {
				return lines, fmt.Errorf("part of lines : %v. %v", i, err)
			}
			lines = append(lines, partLines...)
		}
		return
	}

	index := strings.Index(line, ":")
	if index < 0 {
		err = fmt.Errorf("cannot find `:` in line : %v", line)
		return
	}
	line = line[index+1:]
	parts := strings.Split(line, "\\\n")

	for _, p := range parts {
		p = strings.TrimSpace(p)
		begin := 0
		for i := 0; i <= len(p)-2; i++ {
			if p[i] == '\\' && p[i+1] == ' ' {
				i++
				continue
			}
			if p[i] == ' ' {
				lines = append(lines, p[begin:i])
				begin = i + 1
			}
			if i == len(p)-2 {
				lines = append(lines, p[begin:])
			}
		}
	}
again:
	for i := range lines {
		if lines[i] == "" {
			lines = append(lines[:i], lines[i+1:]...)
			goto again
		}
		lines[i] = strings.Replace(lines[i], "\\", "", -1)
	}

	return
}
