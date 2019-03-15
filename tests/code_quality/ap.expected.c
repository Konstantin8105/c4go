//
//	Package - transpiled by c4go
//
//	If you have found any issues, please raise an issue at:
//	https://github.com/Konstantin8105/c4go/
//

package code_quality

import "unsafe"
import "github.com/Konstantin8105/c4go/noarch"

// a - transpiled function from  C4GO/tests/code_quality/ap.c:4
func a(v1 []int32) {
	// input argument - C-pointer
	noarch.Printf([]byte("a: %d\n\x00"), v1[0])
}

// b - transpiled function from  C4GO/tests/code_quality/ap.c:7
func b(v1 []int32, size int32) {
	{
		// input argument - C-array
		for size -= 1; size >= 0; size-- {
			noarch.Printf([]byte("b: %d %d\n\x00"), size, v1[size])
		}
	}
}

// get - transpiled function from  C4GO/tests/code_quality/ap.c:11
func get() int32 {
	return (0)
}

// main - transpiled function from  C4GO/tests/code_quality/ap.c:15
func main() {
	var i1 int32 = 42
	// value
	a((*[100000000]int32)(unsafe.Pointer(&i1))[:])
	b((*[100000000]int32)(unsafe.Pointer(&i1))[:], 1)
	var i2 []int32 = []int32{11, 22}
	// C-array
	a(i2)
	b(i2, 2)
	var i3 []int32 = (*[100000000]int32)(unsafe.Pointer(&i1))[:]
	// C-pointer from value
	a(i3)
	b(i3, 1)
	var i4 []int32 = i2
	// C-pointer from array
	a(i4)
	b(i4, 2)
	var i5 []int32 = i2[1:]
	// C-pointer from array
	a(i5)
	b(i5, 1)
	var i6 []int32 = i5[0+1:]
	// pointer arithmetic
	a(i6)
	b(i6, 1)
	var val int32 = 2 - 2
	var i7 []int32 = (*(*[1000000000]int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&(*(*[1000000000]int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&(*(*[1000000000]int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&i5[0])) + (uintptr)(1+(1-1)+val+0*(100-2))*unsafe.Sizeof(i5[0]))))[:][0])) + (uintptr)(0)*unsafe.Sizeof((*(*[1000000000]int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&i5[0])) + (uintptr)(1+(1-1)+val+0*(100-2))*unsafe.Sizeof(i5[0]))))[:][0]))))[:][0])) - (uintptr)(0*0)*unsafe.Sizeof((*(*[1000000000]int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&(*(*[1000000000]int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&i5[0])) + (uintptr)(1+(1-1)+val+0*(100-2))*unsafe.Sizeof(i5[0]))))[:][0])) + (uintptr)(0)*unsafe.Sizeof((*(*[1000000000]int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&i5[0])) + (uintptr)(1+(1-1)+val+0*(100-2))*unsafe.Sizeof(i5[0]))))[:][0]))))[:][0]))))[:]
	// pointer arithmetic
	a(i7)
	b(i7, 1)
	var i8 []int32 = i5[0+1+0:]
	// pointer arithmetic
	a(i8)
	b(i8, 1)
	var i9 []int32 = []int32{i3[0], i3[(0 + 1)]}
	// pointer arithmetic
	a(i9)
	b(i9, 1)
	var i10 []int32 = i5[1+0+0+5*get()+get()+(12+3)*get():]
	// pointer arithmetic
	a(i10)
	b(i10, 1)
	var i11 []int32 = (*(*[1000000000]int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&i5[1+0+0+5*get()+get():][0])) - (uintptr)((12+3)*get())*unsafe.Sizeof(i5[1+0+0+5*get()+get():][0]))))[:]
	// pointer arithmetic
	a(i11)
	b(i11, 1)
	return
}
