package noarch

import (
	"reflect"
	"unsafe"
)

// CStringToString returns a string that contains all the bytes in the
// provided C string up until the first NULL character.
func CStringToString(s *byte) string {
	if s == nil {
		return ""
	}

	end := -1
	for i := 0; ; i++ {
		if *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + uintptr(i))) == 0 {
			end = i
			break
		}
	}

	if end == -1 {
		return ""
	}

	return string(toByteSlice(s, end))
}

// StringToCString returns the C string (also known as a null terminated string)
// to be as used as a string in C.
func StringToCString(s string) *byte {
	cString := make([]byte, len(s)+1)
	copy(cString, []byte(s))
	cString[len(s)] = 0

	return &cString[0]
}

// CStringIsNull will test if a C string is NULL. This is equivalent to:
//
//    s == NULL
func CStringIsNull(s *byte) bool {
	if s == nil {
		return true
	}

	return *s == 0
}

// CPointerToGoPointer converts a C-style pointer into a Go-style pointer.
//
// C pointers are represented as slices that have one element pointing to where
// the original C pointer would be referencing. This isn't useful if the pointed
// value needs to be passed to another Go function in these libraries.
//
// See also GoPointerToCPointer.
func CPointerToGoPointer(a interface{}) interface{} {
	t := reflect.TypeOf(a).Elem()

	return reflect.New(t).Elem().Addr().Interface()
}

// toByteSlice returns a byte slice to a with the given length.
func toByteSlice(a *byte, length int) []byte {
	header := reflect.SliceHeader{
		uintptr(unsafe.Pointer(a)),
		int(length),
		int(length),
	}
	return (*(*[]byte)(unsafe.Pointer(&header)))[:]
}

// GoPointerToCPointer does the opposite of CPointerToGoPointer.
//
// A Go pointer (simply a pointer) is converted back into the original slice
// structure (of the original slice reference) so that the calling functions
// will be able to see the new data of that pointer.
func GoPointerToCPointer(destination interface{}, value interface{}) {
	v := reflect.ValueOf(destination).Elem()
	reflect.ValueOf(value).Index(0).Set(v)
}
