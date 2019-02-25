package types

import (
	"bytes"
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

	s = CleanCType(s)

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
			if IsFunction(ff) {
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
	if IsFunction(s) {
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
		if IsFunction(v) {
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
	pr, f, r, err := ParseFunction(s)
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

// IsFunction - return true if string is function like "void (*)(void)"
func IsFunction(s string) bool {
	s = strings.Replace(s, "(*)", "", -1)
	return strings.Contains(s, "(")
}

//IsCPointer - check C type is pointer
func IsCPointer(s string) bool {
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
	return false
}

// IsCArray - check C type is array
func IsCArray(s string) bool {
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
	return false
}

// IsPointer - check type is pointer
func IsPointer(s string) bool {
	return strings.ContainsAny(s, "*[]")
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

// IsTypedefFunction - return true if that type is typedef of function.
func IsTypedefFunction(p *program.Program, s string) bool {
	if v, ok := p.TypedefType[s]; ok && IsFunction(v) {
		return true
	}
	s = string(s[0 : len(s)-len(" *")])
	if v, ok := p.TypedefType[s]; ok && IsFunction(v) {
		return true
	}
	return false
}

// ParseFunction - parsing elements of C function
func ParseFunction(s string) (prefix string, f []string, r []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot parse function '%s' : %v", s, err)
		} else {
			prefix = strings.TrimSpace(prefix)
			for i := range r {
				r[i] = strings.TrimSpace(r[i])
			}
			for i := range f {
				f[i] = strings.TrimSpace(f[i])
			}
		}
	}()

	s = strings.TrimSpace(s)
	if !IsFunction(s) {
		err = fmt.Errorf("Is not function : %s", s)
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
				err = fmt.Errorf("Don't found '(' in type : %s", s)
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
		returns = strings.TrimSpace(s[:pos])
		arguments = strings.TrimSpace(s[pos:])
	}
	if arguments == "" {
		err = fmt.Errorf("Cannot parse (right part is nil) : %v", s)
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
			err = fmt.Errorf("Cannot found '(' in block")
			return
		}
		returns = returns + block[index:]
		prefix = block[:index]
		if strings.Contains(prefix, "*") {
			prefix = strings.Replace(prefix, "*", "", 1)
		} else {
			err = fmt.Errorf("Undefined situation")
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
			err = fmt.Errorf("Cannot found '(' in block")
			return
		}
		prefix = block[:index]
		returns = returns + block[index:]

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
	out = strings.Replace(out, "const", "", -1)
	out = strings.Replace(out, "volatile", "", -1)
	out = strings.Replace(out, "__restrict", "", -1)
	out = strings.Replace(out, "restrict", "", -1)
	out = strings.Replace(out, "_Nullable", "", -1)
	out = strings.Replace(out, "\t", "", -1)
	out = strings.Replace(out, "\n", "", -1)
	out = strings.Replace(out, "\r", "", -1)

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
func GenerateCorrectType(name string) string {
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

// GetAmountArraySize - return amount array size
// Example :
// In  : 'char [40]'
// Out : 40
func GetAmountArraySize(cType string) (size int, err error) {
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
	s = CleanCType(s)
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
