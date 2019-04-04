package types

import (
	"fmt"
	"strings"

	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/util"
)

var sizeofStack []string

// SizeOf returns the number of bytes for a type. This the same as using the
// sizeof operator/function in C.
func SizeOf(p *program.Program, cType string) (size int, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot determine sizeof : |%s|. err = %v", cType, err)
		}
	}()

	// stack check
	for i := range sizeofStack {
		if sizeofStack[i] == cType {
			return 0, fmt.Errorf("sizeof stack loop : %v", sizeofStack)
		}
	}
	sizeofStack = append(sizeofStack, cType) // add to stack
	defer func() {
		// remove last element from stack
		if len(sizeofStack) == 0 {
			err = fmt.Errorf("cannot remove last sizeof stack element: slice is empty")
		} else {
			sizeofStack = sizeofStack[:len(sizeofStack)-1]
		}
	}()

	// Remove keywords that do not effect the size.
	cType = util.CleanCType(cType)
	cType = strings.Replace(cType, "unsigned ", "", -1)
	cType = strings.Replace(cType, "signed ", "", -1)

	// FIXME: The pointer size will be different on different platforms. We
	// should find out the correct size at runtime.
	pointerSize := 8

	if cType == "FILE *" || cType == "FILE" || cType == "struct _IO_FILE *" {
		return pointerSize, nil
	}

	// Enum with name
	if strings.HasPrefix(cType, "enum") {
		return SizeOf(p, "int")
	}

	// typedef int Integer;
	if v, ok := p.GetBaseTypeOfTypedef(cType); ok {
		return SizeOf(p, v)
	}

	// typedef Enum
	if _, ok := p.EnumTypedefName[cType]; ok {
		return SizeOf(p, "int")
	}

	// A structure will be the sum of its parts.
	var isStruct, ok bool
	var s *program.Struct
	cType = util.GenerateCorrectType(cType)
	if s, ok = p.Structs[cType]; ok {
		isStruct = true
	} else if s, ok = p.Structs["struct "+cType]; ok {
		isStruct = true
	} else if strings.HasPrefix(cType, "struct ") {
		if s, ok = p.Structs[cType[len("struct "):]]; ok {
			isStruct = true
		}
	}
	if isStruct {
		totalBytes := 0

		for _, t := range s.Fields {
			var bytes int
			var err error

			switch f := t.(type) {
			case string:
				bytes, err = SizeOf(p, f)

			case *program.Struct:
				bytes, err = SizeOf(p, f.Name)
			}

			if err != nil {
				err = fmt.Errorf(
					"Cannot calculate `struct` sizeof for `%T`. bytes = '%v'. %v",
					t, bytes, err)
				return 0, err
			}
			totalBytes += bytes
		}

		// The size of a struct is rounded up to fit the size of the pointer of
		// the OS.
		if totalBytes%pointerSize != 0 {
			totalBytes += pointerSize - (totalBytes % pointerSize)
		}

		return totalBytes, nil
	}

	// An union will be the max size of its parts.
	var isUnion bool
	cType = util.GenerateCorrectType(cType)
	if _, ok := p.Unions[cType]; ok {
		isUnion = true
	} else if _, ok := p.Unions["union "+cType]; ok {
		isUnion = true
	}
	if isUnion {
		byteCount := 0

		s, ok := p.Unions[cType]
		if !ok {
			return 0, fmt.Errorf("error in union")
		}

		for _, t := range s.Fields {
			var bytes int

			switch f := t.(type) {
			case string:
				bytes, err = SizeOf(p, f)

			case *program.Struct:
				bytes, err = SizeOf(p, f.Name)
			}

			if err != nil {
				err = fmt.Errorf("Cannot canculate `union` sizeof for `%T`. %v",
					t, err)
				return 0, err
			}

			if byteCount < bytes {
				byteCount = bytes
			}
		}

		// The size of an union is rounded up to fit the size of the pointer of
		// the OS.
		if byteCount%pointerSize != 0 {
			byteCount += pointerSize - (byteCount % pointerSize)
		}

		return byteCount, nil
	}

	// Function pointers are one byte?
	if strings.Contains(cType, "(") {
		return 1, nil
	}

	if strings.HasSuffix(cType, "*") {
		return pointerSize, nil
	}

	switch cType {
	case "char", "void":
		return 1, nil

	case "short":
		return 2, nil

	case "int", "float", "long int":
		return 4, nil

	case "long", "long long", "long long int", "double":
		return 8, nil

	case "long double", "long long unsigned int":
		return 16, nil
	}

	// definition type
	if t, ok := program.DefinitionType[cType]; ok && cType != t {
		return SizeOf(p, t)
	}

	// resolved type
	conv := func(t string) (bytes int, ok bool) {
		switch t {
		case "byte", "int8", "uint8":
			return 1, true

		case "int16", "uint16":
			return 2, true

		case "int32", "uint32", "rune", "float32":
			return 4, true

		case "int64", "uint64", "float64", "complex64", "uintptr", "int", "uint":
			return 8, true

		case "complex128":
			return 16, true
		}
		return -1, false
	}
	if t, ok := conv(cType); ok {
		return t, nil
	}
	if r, err := ResolveType(p, cType); err != nil {
		if t, ok := conv(r); ok {
			return t, nil
		}
	}

	// Get size for array types like: `base_type [count]`
	totalArraySize := 1
	arrayType, arraySize := GetArrayTypeAndSize(cType)
	if arraySize <= 0 {
		return 0, nil
	}

	for arraySize != -1 {
		totalArraySize *= arraySize
		arrayType, arraySize = GetArrayTypeAndSize(arrayType)
	}

	baseSize, err := SizeOf(p, arrayType)
	if err != nil {
		return 0, fmt.Errorf("error in sizeof baseSize for `%v`",
			arrayType)
	}

	return baseSize * totalArraySize, nil
}
