package program

// CStdStructType - conversion map from C standart library structures to
// c4go structures
var CStdStructType = map[string]string{
	// time.h
	"time_t": "github.com/Konstantin8105/c4go/noarch.TimeT",

	"fpos_t": "int",

	// built-in
	"bool":                   "bool",
	"char *":                 "[]byte",
	"char":                   "byte",
	"char*":                  "[]byte",
	"double":                 "float64",
	"float":                  "float32",
	"int":                    "int",
	"long double":            "float64",
	"long int":               "int32",
	"long long":              "int64",
	"long long int":          "int64",
	"long long unsigned int": "uint64",
	"long unsigned int":      "uint32",
	"long":                   "int32",
	"short":                  "int16",
	"signed char":            "int8",
	"unsigned char":          "uint8",
	"unsigned int":           "uint32",
	"unsigned long long":     "uint64",
	"unsigned long":          "uint32",
	"unsigned short":         "uint16",
	"unsigned short int":     "uint16",
	"void":                   "",
	"_Bool":                  "int",
	"size_t":                 "uint",
	"ptrdiff_t":              "github.com/Konstantin8105/c4go/noarch.PtrdiffT",

	// void*
	"void*":  "interface{}",
	"void *": "interface{}",

	// null is a special case (it should probably have a less ambiguos name)
	// when using the NULL macro.
	"null": "null",

	// Non platform-specific types.
	"uint32":     "uint32",
	"uint64":     "uint64",
	"__uint16_t": "uint16",
	"__uint32_t": "uint32",
	"__uint64_t": "uint64",

	// These are special cases that almost certainly don't work. I've put
	// them here because for whatever reason there is no suitable type or we
	// don't need these platform specific things to be implemented yet.
	"__builtin_va_list": "int64",
	"unsigned __int128": "uint64",
	"__int128":          "int64",
	"__mbstate_t":       "int64",
	"__sbuf":            "int64",
	"__sFILEX":          "interface{}",
	"FILE":              "github.com/Konstantin8105/c4go/noarch.File",
}

func (p *Program) initializationStructs() {
	// stdlib.h
	p.Unions["div_t"] = &Struct{
		Name: "div_t",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "int",
			"rem":  "int",
		},
	}
	p.Unions["c4go_div_t"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.DivT",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "Quot",
			"rem":  "Rem",
		},
	}

	// stdlib.h
	p.Unions["ldiv_t"] = &Struct{
		Name: "ldiv_t",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "long int",
			"rem":  "long int",
		},
	}
	p.Unions["c4go_ldiv_t"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.LdivT",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "Quot",
			"rem":  "Rem",
		},
	}

	// stdlib.h
	p.Unions["lldiv_t"] = &Struct{
		Name: "lldiv_t",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "long long",
			"rem":  "long long",
		},
	}
	p.Unions["c4go_lldiv_t"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.LldivT",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "Quot",
			"rem":  "Rem",
		},
	}

	// time.h
	p.Unions["struct tm"] = &Struct{
		Name: "struct tm",
		Type: StructType,
		Fields: map[string]interface{}{
			"tm_sec":   "int",
			"tm_min":   "int",
			"tm_hour":  "int",
			"tm_mday":  "int",
			"tm_mon":   "int",
			"tm_year":  "int",
			"tm_wday":  "int",
			"tm_yday":  "int",
			"tm_isdst": "int",
		},
	}
	p.Unions["c4go_struct tm"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.Tm",
		Type: StructType,
		Fields: map[string]interface{}{
			"tm_sec":   "TmSec",
			"tm_min":   "TmMin",
			"tm_hour":  "TmHour",
			"tm_mday":  "TmMday",
			"tm_mon":   "TmMon",
			"tm_year":  "TmYear",
			"tm_wday":  "TmWday",
			"tm_yday":  "TmYday",
			"tm_isdst": "TmIsdst",
		},
	}
}
