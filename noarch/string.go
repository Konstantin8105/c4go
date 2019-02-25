package noarch

import (
	"bytes"
)

// Strlen returns the length of a string.
//
// The length of a C string is determined by the terminating null-character: A
// C string is as long as the number of characters between the beginning of the
// string and the terminating null character (without including the terminating
// null character itself).
func Strlen(a []byte) int {
	// TODO: The transpiler should have a syntax that means this proxy function
	// does not need to exist.

	return len(CStringToString(a))
}

// Strcpy copies the C string pointed by source into the array pointed by
// destination, including the terminating null character (and stopping at that
// point).
//
// To avoid overflows, the size of the array pointed by destination shall be
// long enough to contain the same C string as source (including the terminating
// null character), and should not overlap in memory with source.
func Strcpy(dest, src []byte) []byte {
	for i, c := range src {
		dest[i] = c

		// We only need to copy until the first NULL byte. Make sure we also
		// include that NULL byte on the end.
		if c == '\x00' {
			break
		}
	}

	return dest
}

// Strncpy copies the first num characters of source to destination. If the end
// of the source C string (which is signaled by a null-character) is found
// before num characters have been copied, destination is padded with zeros
// until a total of num characters have been written to it.
//
// No null-character is implicitly appended at the end of destination if source
// is longer than num. Thus, in this case, destination shall not be considered a
// null terminated C string (reading it as such would overflow).
//
// destination and source shall not overlap (see memmove for a safer alternative
// when overlapping).
func Strncpy(dest, src []byte, len int) []byte {
	// Copy up to the len or first NULL bytes - whichever comes first.
	i := 0
	for ; i < len && src[i] != 0; i++ {
		dest[i] = src[i]
	}

	// The rest of the dest will be padded with zeros to the len.
	for ; i < len; i++ {
		dest[i] = 0
	}

	return dest
}

// Strcat - concatenate strings
// Appends a copy of the source string to the destination string.
// The terminating null character in destination is overwritten by the first
// character of source, and a null-character is included at the end
// of the new string formed by the concatenation of both in destination.
func Strcat(dest, src []byte) []byte {
	Strcpy(dest[Strlen(dest):], src)
	return dest
}

// Strncat - concatenate strings
// Appends at most count characters of the source string to the destination string.
// The terminating null character in destination is overwritten by the first
// character of source, and a null-character is included at the end
// of the new string formed by the concatenation of both in destination.
func Strncat(dest, src []byte, len int) []byte {
	Strncpy(dest[Strlen(dest):], src, len)
	return dest
}

// Strcmp - compare two strings
// Compares the C string str1 to the C string str2.
func Strcmp(str1, str2 []byte) int {
	return bytes.Compare([]byte(CStringToString(str1)), []byte(CStringToString(str2)))
}

// Strchr - Locate first occurrence of character in string
// See: http://www.cplusplus.com/reference/cstring/strchr/
func Strchr(str []byte, ch int) []byte {
	i := 0
	for {
		if str[i] == '\x00' {
			break
		}
		if int(str[i]) == ch {
			return str[i:]
		}
		i++
	}
	return nil
}

// Strstr finds the first occurrence of the null-terminated byte string
// pointed to by substr in the null-terminated byte string pointed to by str.
//The terminating null characters are not compared.
func Strstr(str, subStr []byte) []byte {
	if subStr == nil {
		return str
	}
	if subStr[0] == '\x00' {
		return str
	}

	k := 0
	for i := range subStr {
		if subStr[i] == '\x00' {
			k = i
			break
		}
	}

	index := bytes.Index(str, subStr[:k])
	if index < 0 {
		return nil
	}
	return str[index:]
}

// Memset sets the first num bytes of the block of memory pointed by ptr to
// the specified value (interpreted as an unsigned char)
func Memset(ptr []byte, value byte, num uint32) []byte {
	for i := 0; uint32(i) < num; i++ {
		ptr[i] = value
	}
	return ptr
}

// Memmove move block of memory
func Memmove(ptr, source interface{}, num uint32) interface{} {
	p1 := ptr.([]byte)
	p2 := source.([]byte)
	for i := int(num); i >= 0; i-- {
		p1[i] = p2[i]
	}
	return p1
}

// Memcmp - compare two buffers
// Compares the first count characters of the objects pointed to be lhs and rhs
func Memcmp(lhs []byte, rhs []byte, count uint32) int {
	for i := 0; uint32(i) < count; i++ {
		if int(lhs[i]) < int(rhs[i]) {
			return -1
		} else if int(lhs[i]) > int(rhs[i]) {
			return 1
		}
	}
	return 0
}

func Memcpy(dest, source interface{}, num uint) interface{} {
	p1 := dest.([]byte)
	p2 := source.([]byte)
	for i := 0; i < int(num); i++ {
		if len(p2) <= i {
			break
		}
		p1[i] = p2[i]
	}
	return nil
}

func Strrchr(source []byte, c int) []byte {
	ch := byte(c)
	pos := len(source)
	for i := range source {
		if source[i] == '\x00' {
			pos = i
			break
		}
	}
	for i := pos; i >= 0; i-- {
		if source[i] == ch {
			return source[i:]
		}
	}
	return source
}

func Strdup(s []byte) []byte {
	return s
}

func Strerror(e int) []byte {
	return []byte("strerror")
}
