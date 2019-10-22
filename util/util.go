package util

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// InStrings returns true if item exists in items. It must be an exact string
// match.
func InStrings(item string, items []string) bool {
	for _, v := range items {
		if item == v {
			return true
		}
	}

	return false
}

// Ucfirst returns the word with the first letter uppercased; none of the other
// letters in the word are modified. For example "fooBar" would return "FooBar".
func Ucfirst(word string) string {
	if word == "" {
		return ""
	}

	if len(word) == 1 {
		return strings.ToUpper(word)
	}

	return strings.ToUpper(string(word[0])) + word[1:]
}

// Atoi converts a string to an integer in cases where we are sure that s will
// be a valid integer, otherwise it will panic.
func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	PanicOnError(err, "bad integer")

	return i
}

// GetExportedName returns a deterministic and Go safe name for a C type. For
// example, "*__foo[]" will return "FooSlice".
func GetExportedName(field string) string {
	if strings.Contains(field, "interface{}") ||
		strings.Contains(field, "Interface{}") {
		return "Interface"
	}

	// Convert "[]byte" into "byteSlice". This also works with multiple slices,
	// like "[][]byte" to "byteSliceSlice".
	for len(field) > 2 && field[:2] == "[]" {
		field = field[2:] + "Slice"
	}

	// NotFunc(int)()
	field = strings.Replace(field, "(", "_", -1)
	field = strings.Replace(field, ")", "_", -1)

	return Ucfirst(strings.TrimLeft(field, "*_"))
}

// IsFunction - return true if string is function like "void (*)(void)"
func IsFunction(s string) bool {
	s = strings.Replace(s, "(*)", "", -1)
	return strings.Contains(s, "(")
}

// IsLastArray - check type have array '[]'
func IsLastArray(s string) bool {
	for _, b := range s {
		switch b {
		case '[':
			return true
		case '*':
			break
		}
	}
	return false
}

// ParseFunction - parsing elements of C function
func ParseFunction(s string) (prefix string, funcname string, f []string, r []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot parse function '%s' : %v", s, err)
		} else {
			prefix = strings.TrimSpace(prefix)
			funcname = strings.TrimSpace(funcname)
			for i := range r {
				r[i] = strings.TrimSpace(r[i])
			}
			for i := range f {
				f[i] = strings.TrimSpace(f[i])
			}
		}
	}()

	// remove specific attribute for function longjmp
	s = strings.Replace(s, "__attribute__((noreturn))", "", -1)

	s = strings.TrimSpace(s)
	if !IsFunction(s) {
		err = fmt.Errorf("is not function : %s", s)
		return
	}
	var returns string
	var arguments string
	{
		// Example of function types :
		// int (*)(int, float)
		// int (int, float)
		// int (*)(int (*)(int))
		// void (*(*)(int *, void *, const char *))(void)
		if s[len(s)-1] != ')' {
			err = fmt.Errorf("function type |%s| haven't last symbol ')'", s)
			return
		}
		counter := 1
		var pos int
		for i := len(s) - 2; i >= 0; i-- {
			if i == 0 {
				err = fmt.Errorf("don't found '(' in type : %s", s)
				return
			}
			if s[i] == ')' {
				counter++
			}
			if s[i] == '(' {
				counter--
			}
			if counter == 0 {
				pos = i
				break
			}
		}
		// s[:pos] = `speed_t cfgetospeed`
		if unicode.IsNumber(rune(s[pos-1])) || unicode.IsLetter(rune(s[pos-1])) {
			for i := pos - 1; i >= 0; i-- {
				if s[i] == ' ' {
					funcname = s[i+1 : pos]
					returns = strings.TrimSpace(s[:i])
					break
				}
			}
		} else {
			returns = strings.TrimSpace(s[:pos])
		}
		arguments = strings.TrimSpace(s[pos:])
	}
	if arguments == "" {
		err = fmt.Errorf("cannot parse (right part is nil) : %v", s)
		return
	}
	// separate fields of arguments
	{
		pos := 1
		counter := 0
		for i := 1; i < len(arguments)-1; i++ {
			if arguments[i] == '(' {
				counter++
			}
			if arguments[i] == ')' {
				counter--
			}
			if counter == 0 && arguments[i] == ',' {
				f = append(f, strings.TrimSpace(arguments[pos:i]))
				pos = i + 1
			}
		}
		f = append(f, strings.TrimSpace(arguments[pos:len(arguments)-1]))
	}

	// returns
	// Example:  __ssize_t
	if returns[len(returns)-1] != ')' {
		r = append(r, returns)
		return
	}

	// Example: void  ( *(*)(int *, void *, char *))
	//          -------  --------------------------- return type
	//                 ==                            prefix
	//                ++++++++++++++++++++++++++++++ block
	// return type : void (*)(int *, void *, char *)
	// prefix      : *
	// Find the block
	var counter int
	var position int
	for i := len(returns) - 1; i >= 0; i-- {
		if returns[i] == ')' {
			counter++
		}
		if returns[i] == '(' {
			counter--
		}
		if counter == 0 {
			position = i
			break
		}
	}
	block := string([]byte(returns[position:]))
	returns = returns[:position]

	// Examples returns:
	// int   (*)
	// char *(*)
	// block is : (*)
	if block == "(*)" {
		r = append(r, returns)
		return
	}

	index := strings.Index(block, "(*)")
	if index < 0 {
		if strings.Count(block, "(") == 1 {
			// Examples returns:
			// int   ( * [2])
			// ------         return type
			//        ======  prefix
			//       ++++++++ block
			bBlock := []byte(block)
			for i := 0; i < len(bBlock); i++ {
				switch bBlock[i] {
				case '(', ')':
					bBlock[i] = ' '
				}
			}
			bBlock = bytes.Replace(bBlock, []byte("*"), []byte(""), 1)
			prefix = string(bBlock)
			r = append(r, returns)
			return
		}
		// void (*(int *, void *, const char *))
		//      ++++++++++++++++++++++++++++++++ block
		block = block[1 : len(block)-1]
		index := strings.Index(block, "(")
		if index < 0 {
			err = fmt.Errorf("cannot found '(' in block")
			return
		}
		returns = returns + block[index:]
		prefix = block[:index]
		if strings.Contains(prefix, "*") {
			prefix = strings.Replace(prefix, "*", "", 1)
		} else {
			err = fmt.Errorf("undefined situation")
			return
		}
		r = append(r, returns)
		return
	}
	if len(block)-1 > index+3 && block[index+3] == '(' {
		// Examples returns:
		// void  ( *(*)(int *, void *, char *))
		//       ++++++++++++++++++++++++++++++ block
		//            ^^                        check this
		block = strings.Replace(block, "(*)", "", 1)
		block = block[1 : len(block)-1]
		index := strings.Index(block, "(")
		if index < 0 {
			err = fmt.Errorf("cannot found '(' in block")
			return
		}

		returns = returns + block[index:]
		// example of block[:index]
		// `*signal`
		// `* signal`
		if pr := strings.TrimSpace(block[:index]); unicode.IsLetter(rune(pr[len(pr)-1])) ||
			unicode.IsNumber(rune(pr[len(pr)-1])) {
			pr = strings.Replace(pr, "*", " * ", -1)
			for i := len(pr) - 1; i >= 0; i-- {
				if unicode.IsLetter(rune(pr[i])) {
					continue
				}
				if unicode.IsNumber(rune(pr[i])) {
					continue
				}
				prefix = pr[:i]
				funcname = pr[i:]
				break
			}
		} else {
			prefix = block[:index]
		}

		r = append(r, returns)
		return
	}

	// Examples returns:
	// int   ( *( *(*)))
	// -----              return type
	//        =========   prefix
	//       +++++++++++  block
	bBlock := []byte(block)
	for i := 0; i < len(bBlock); i++ {
		switch bBlock[i] {
		case '(', ')':
			bBlock[i] = ' '
		}
	}
	bBlock = bytes.Replace(bBlock, []byte("*"), []byte(""), 1)
	prefix = string(bBlock)
	r = append(r, returns)

	return
}

// CleanCType - remove from C type not Go type
func CleanCType(s string) (out string) {
	out = s

	// remove space from pointer symbols
	out = strings.Replace(out, "* *", "**", -1)

	// add space for simplification redactoring
	out = strings.Replace(out, "*", " *", -1)

	out = strings.Replace(out, "( *)", "(*)", -1)

	// Remove any whitespace or attributes that are not relevant to Go.
	out = strings.Replace(out, "\t", "", -1)
	out = strings.Replace(out, "\n", "", -1)
	out = strings.Replace(out, "\r", "", -1)
	list := []string{"const", "volatile", "__restrict", "restrict", "_Nullable"}
	for _, word := range list {
		// example :
		// `const`
		if out == word {
			out = ""
			continue
		}

		// examples :
		// `const char  * *`
		// `const struct parg_option`
		// `void (*)(int  *, void  *, const char  *)`
		out = strings.Replace(out, " "+word+" ", "", -1)
		out = strings.Replace(out, " "+word+"*", "*", -1)
		out = strings.Replace(out, "*"+word+" ", "*", -1)
		out = strings.Replace(out, "*"+word+"*", "* *", -1)

		if pr := word + " "; strings.HasPrefix(out, pr) {
			out = out[len(word):]
		}
		if po := " " + word; strings.HasSuffix(out, po) {
			out = out[:len(out)-len(word)]
		}

		if pr := word + "*"; strings.HasPrefix(out, pr) {
			out = out[len(word):]
		}
		if po := "*" + word; strings.HasSuffix(out, po) {
			out = out[:len(out)-len(word)]
		}

		out = strings.TrimSpace(out)
	}

	// remove space from pointer symbols
	out = strings.Replace(out, "* *", "**", -1)
	out = strings.Replace(out, "[", " [", -1)
	out = strings.Replace(out, "] [", "][", -1)

	// remove addition spaces
	out = strings.Replace(out, "  ", " ", -1)

	// remove spaces around
	out = strings.TrimSpace(out)

	if out != s {
		return CleanCType(out)
	}

	return out
}

// GenerateCorrectType - generate correct type
// Example: 'union (anonymous union at tests/union.c:46:3)'
func GenerateCorrectType(name string) (result string) {
	if !strings.Contains(name, "anonymous") {
		return CleanCType(name)
	}
	index := strings.Index(name, "(anonymous")
	if index < 0 {
		return name
	}
	name = strings.Replace(name, "anonymous", "", 1)
	var last int
	for last = index; last < len(name); last++ {
		if name[last] == ')' {
			break
		}
	}

	// Create a string, for example:
	// Input (name)   : 'union (anonymous union at tests/union.c:46:3)'
	// Output(inside) : '(anonymous union at tests/union.c:46:3)'
	inside := string(([]byte(name))[index : last+1])

	// change unacceptable C name letters
	inside = strings.Replace(inside, "(", "_", -1)
	inside = strings.Replace(inside, ")", "_", -1)
	inside = strings.Replace(inside, " ", "_", -1)
	inside = strings.Replace(inside, ":", "_", -1)
	inside = strings.Replace(inside, "/", "_", -1)
	inside = strings.Replace(inside, "-", "_", -1)
	inside = strings.Replace(inside, "\\", "_", -1)
	inside = strings.Replace(inside, ".", "_", -1)
	out := string(([]byte(name))[0:index]) + inside + string(([]byte(name))[last+1:])

	// For case:
	// struct siginfo_t::(anonymous at /usr/include/x86_64-linux-gnu/bits/siginfo.h:119:2)
	// we see '::' before 'anonymous' word
	out = strings.Replace(out, ":", "D", -1)

	return CleanCType(out)
}
