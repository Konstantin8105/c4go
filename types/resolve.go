package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/util"
)

// cIntegerType - slice of C integer type
var cIntegerType = []string{
	"int",
	"long long",
	"long long int",
	"long long unsigned int",
	"long unsigned int",
	"long",
	"short",
	"unsigned int",
	"unsigned long long",
	"unsigned long",
	"unsigned short",
	"unsigned short int",
	"size_t",
	"ptrdiff_t",
}

func IsSigned(p *program.Program, cType string) bool {
	if !strings.Contains(cType, "unsigned") {
		return true
	}
	if rt, ok := p.TypedefType[cType]; ok {
		return IsSigned(p, rt)
	}
	return false
}

// IsCInteger - return true is C type integer
func IsCInteger(p *program.Program, cType string) bool {
	for i := range cIntegerType {
		if cType == cIntegerType[i] {
			return true
		}
	}
	if rt, ok := p.TypedefType[cType]; ok {
		return IsCInteger(p, rt)
	}
	return false
}

var cFloatType = []string{
	"double",
	"float",
	"long double",
}

// IsCInteger - return true is C type integer
func IsCFloat(p *program.Program, cType string) bool {
	for i := range cFloatType {
		if cType == cFloatType[i] {
			return true
		}
	}
	if rt, ok := p.TypedefType[cType]; ok {
		return IsCFloat(p, rt)
	}
	return false
}

// NullPointer - is look : (double *)(nil) or (FILE *)(nil)
// created only for transpiler.CStyleCastExpr
var NullPointer = "NullPointerType *"

// ToVoid - specific type for ignore the cast
var ToVoid = "ToVoid"

// ResolveType determines the Go type from a C type.
//
// Some basic examples are obvious, such as "float" in C would be "float32" in
// Go. But there are also much more complicated examples, such as compound types
// (structs and unions) and function pointers.
//
// Some general rules:
//
// 1. The Go type must be deterministic. The same C type will ALWAYS return the
//    same Go type, in any condition. This is extremely important since the
//    nature of C is that is may not have certain information available about the
//    rest of the program or libraries when it is being compiled.
//
// 2. Many C type modifiers and properties are lost as they have no sensible or
//    valid translation to Go. Some example of those would be "const" and
//    "volatile". It is left be up to the clang (or other compiler) to warn if
//    types are being abused against the standards in which they are being
//    compiled under. Go will make no assumptions about how you expect it act,
//    only how it is used.
//
// 3. New types are registered (discovered) throughout the transpiling of the
//    program, so not all types are know at any given time. This works exactly
//    the same way in a C compiler that will not let you use a type before it
//    has been defined.
//
// 4. If all else fails an error is returned. However, a type (which is almost
//    certainly incorrect) "interface{}" is also returned. This is to allow the
//    transpiler to step over type errors and put something as a placeholder
//    until a more suitable solution is found for those cases.
func ResolveType(p *program.Program, s string) (resolveResult string, err error) {
	defer func() {
		resolveResult = strings.TrimSpace(resolveResult)
		if err != nil {
			err = fmt.Errorf("Cannot resolve type '%s' : %v", s, err)
		}
	}()

	if strings.Contains(s, ":") {
		return "interface{}", errors.New("probably an incorrect type translation 0")
	}

	s = util.CleanCType(s)

	// FIXME: This is a hack to avoid casting in some situations.
	if s == "" {
		return "interface{}", errors.New("probably an incorrect type translation 1")
	}

	// FIXME: I have no idea, how to solve.
	// See : /issues/628
	if strings.Contains(s, "__locale_data") {
		s = strings.Replace(s, "struct __locale_data", "int", -1)
		s = strings.Replace(s, "__locale_data", "int", -1)
	}
	if strings.Contains(s, "__locale_struct") {
		return "int", nil
	}
	if strings.Contains(s, "__sFILEX") {
		s = strings.Replace(s, "__sFILEX", "int", -1)
	}

	// The simple resolve types are the types that we know there is an exact Go
	// equivalent. For example float, int, etc.
	if v, ok := program.DefinitionType[s]; ok {
		return p.ImportType(v), nil
	}

	// function type is pointer in Go by default
	if len(s) > 2 {
		base := s[:len(s)-2]
		if ff, ok := p.TypedefType[base]; ok {
			if util.IsFunction(ff) {
				return base, nil
			}
		}
	}

	// No need resolve typedef types
	if _, ok := p.TypedefType[s]; ok {
		if tt, ok := program.DefinitionType[s]; ok {
			// "div_t":   "github.com/Konstantin8105/c4go/noarch.DivT",
			ii := p.ImportType(tt)
			return ii, nil
		}
		return s, nil
	}
	if tt, ok := program.DefinitionType[s]; ok {
		// "div_t":   "github.com/Konstantin8105/c4go/noarch.DivT",
		ii := p.ImportType(tt)
		return ii, nil
	}

	// For function
	if util.IsFunction(s) {
		g, e := resolveFunction(p, s)
		return g, e
	}

	// Example of cases:
	// "int []"
	// "int [6]"
	// "struct s [2]"
	// "int [2][3]"
	// "unsigned short [512]"
	//
	// Example of resolving:
	// int [2][3] -> [][]int
	// int [2][3][4] -> [][][]int
	if s[len(s)-1] == ']' {
		index := strings.LastIndex(s, "[")
		if index < 0 {
			err = fmt.Errorf("Cannot found [ in type : %v", s)
			return
		}
		r := strings.TrimSpace(s[:index])
		r, err = ResolveType(p, r)
		if err != nil {
			err = fmt.Errorf("Cannot []: %v", err)
			return
		}
		return "[]" + r, nil
	}

	// Example :
	// "int (*)"
	if strings.Contains(s, "(*)") {
		return ResolveType(p, strings.Replace(s, "(*)", "*", -1))
	}

	// Check is it typedef enum
	if _, ok := p.EnumTypedefName[s]; ok {
		return ResolveType(p, "int")
	}

	if v, ok := p.TypedefType[s]; ok {
		if util.IsFunction(v) {
			// typedef function
			return s, nil
		}
		return ResolveType(p, v)
	}

	// If the type is already defined we can proceed with the same name.
	if p.IsTypeAlreadyDefined(s) {
		if strings.HasPrefix(s, "struct ") {
			s = s[len("struct "):]
		}
		if strings.HasPrefix(s, "union ") {
			s = s[len("union "):]
		}
		if strings.HasPrefix(s, "class ") {
			s = s[len("class "):]
		}
		// TODO : Why here ImportType???
		return p.ImportType(s), nil
	}

	// It may be a pointer of a simple type. For example, float *, int *,
	// etc.
	if s[len(s)-1] == '*' {
		r := strings.TrimSpace(s[:len(s)-1])
		r, err = ResolveType(p, r)
		if err != nil {
			err = fmt.Errorf("Cannot resolve star `*` for %v : %v", s, err)
			return s, err
		}
		prefix := "[]"
		if strings.Contains(r, "noarch.File") {
			prefix = "*"
		}
		return prefix + r, err
	}

	// Structures are by name.
	var isUnion bool
	{
		isUnion = strings.HasPrefix(s, "struct ") ||
			strings.HasPrefix(s, "class ") ||
			strings.HasPrefix(s, "union ")
		if str := p.GetStruct(s); str != nil {
			isUnion = true
		}
	}
	if isUnion {
		if str := p.GetStruct("c4go_" + s); str != nil {
			s = str.Name
		} else {
			if strings.HasPrefix(s, "struct ") {
				s = s[len("struct "):]
			}
			if strings.HasPrefix(s, "union ") {
				s = s[len("union "):]
			}
			if strings.HasPrefix(s, "class ") {
				s = s[len("class "):]
			}
		}
		return p.ImportType(s), nil
	}

	// Enums are by name.
	if strings.HasPrefix(s, "enum ") {
		s = s[len("enum "):]
		return s, nil
	}

	// I have no idea how to handle this yet.
	if strings.Contains(s, "anonymous union") {
		return "interface{}", errors.New("probably an incorrect type translation 3")
	}

	// Function pointers are not yet supported. In the mean time they will be
	// replaced with a type that certainly wont work until we can fix this
	// properly.
	search := util.GetRegex("[\\w ]+\\(\\*.*?\\)\\(.*\\)").MatchString(s)
	if search {
		return "interface{}",
			fmt.Errorf("function pointers are not supported [1] : '%s'", s)
	}

	search = util.GetRegex("[\\w ]+ \\(.*\\)").MatchString(s)
	if search {
		return "interface{}",
			fmt.Errorf("function pointers are not supported [2] : '%s'", s)
	}

	errMsg := fmt.Sprintf(
		"I couldn't find an appropriate Go type for the C type '%s'.", s)
	return "interface{}", errors.New(errMsg)
}

// resolveType determines the Go type from a C type.
func resolveFunction(p *program.Program, s string) (goType string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("resolveFunction error for `%s` : %v", s, err)
		}
	}()
	var prefix string
	var f, r []string
	prefix, f, r, err = SeparateFunction(p, s)
	if err != nil {
		return
	}
	goType = strings.Replace(prefix, "*", "[]", -1)
	goType += "func("
	for i := range f {
		goType += fmt.Sprintf("%s", f[i])
		if i < len(f)-1 {
			goType += " , "
		}
	}
	goType += ")("
	for i := range r {
		goType += fmt.Sprintf("%s", r[i])
		if i < len(r)-1 {
			goType += " , "
		}
	}
	goType += ")"
	return
}

// SeparateFunction separate a function C type to Go types parts.
func SeparateFunction(p *program.Program, s string) (
	prefix string, fields []string, returns []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot separate function '%s' : %v", s, err)
		}
	}()
	pr, _, f, r, err := util.ParseFunction(s)
	if err != nil {
		return
	}
	for i := range f {
		if f[i] == "" {
			continue
		}
		var t string
		t, err = ResolveType(p, f[i])
		if err != nil {
			err = fmt.Errorf("Error in field %s. err = %v", t, err)
			return
		}
		fields = append(fields, t)
	}
	for i := range r {
		var t string
		t, err = ResolveType(p, r[i])
		if err != nil {
			err = fmt.Errorf("Error in return field %s. err = %v", t, err)
			return
		}
		returns = append(returns, t)
	}
	prefix = pr
	return
}

// IsTypedefFunction - return true if that type is typedef of function.
func IsTypedefFunction(p *program.Program, s string) bool {
	if v, ok := p.TypedefType[s]; ok && util.IsFunction(v) {
		return true
	}
	s = string(s[0 : len(s)-len(" *")])
	if v, ok := p.TypedefType[s]; ok && util.IsFunction(v) {
		return true
	}
	return false
}

// GetAmountArraySize - return amount array size
// Example :
// In  : 'char [40]'
// Out : 40
func GetAmountArraySize(cType string, p *program.Program) (size int, err error) {
	if !IsCArray(cType, p) {
		err = fmt.Errorf("Is not array: `%s`", cType)
		return
	}

	reg := util.GetRegex("\\[(?P<size>\\d+)\\]")
	match := reg.FindStringSubmatch(cType)

	if reg.NumSubexp() != 1 {
		err = fmt.Errorf("Cannot found size of array in type : %s", cType)
		return
	}

	result := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return strconv.Atoi(result["size"])
}

// GetBaseType - return base type without pointera, array symbols
// Input:
// s =  struct BSstructSatSShomeSlepriconSgoSsrcSgithubPcomSD260D18E [7]
func GetBaseType(s string) string {
	s = strings.TrimSpace(s)
	s = util.CleanCType(s)
	if len(s) < 1 {
		return s
	}
	if s[len(s)-1] == ']' {
		for i := len(s) - 1; i >= 0; i-- {
			if s[i] == '[' {
				s = s[:i]
				return GetBaseType(s)
			}
		}
	}
	if s[len(s)-1] == '*' {
		return GetBaseType(s[:len(s)-1])
	}
	if strings.Contains(s, "(*)") {
		return GetBaseType(strings.TrimSpace(
			strings.Replace(s, "(*)", "*", -1)))
	}
	return s
}

// IsCArray - check C type is array
func IsCArray(s string, p *program.Program) bool {
	if len(s) == 0 {
		return false
	}
	for i := len(s) - 1; i >= 0; i-- {
		switch s[i] {
		case ' ':
			continue
		case ']':
			return true
		default:
			break
		}
	}
	if p != nil {
		if t, ok := p.TypedefType[s]; ok {
			return IsCArray(t, p)
		}
	}
	return false
}

// IsPointer - check type is pointer
func IsPointer(s string, p *program.Program) bool {
	return IsCPointer(s, p) || IsCArray(s, p)
}

//IsCPointer - check C type is pointer
func IsCPointer(s string, p *program.Program) bool {
	if len(s) == 0 {
		return false
	}
	for i := len(s) - 1; i >= 0; i-- {
		switch s[i] {
		case ' ':
			continue
		case '*':
			return true
		default:
			break
		}
	}
	if p != nil {
		if t, ok := p.TypedefType[s]; ok {
			return IsCPointer(t, p)
		}
	}
	return false
}
