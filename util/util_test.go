package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestUcfirst(t *testing.T) {
	tcs := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"a", "A"},
		{"w", "W"},
		{"wa", "Wa"},
	}

	for index, tc := range tcs {
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			a := Ucfirst(tc.in)
			if a != tc.out {
				t.Errorf("Result is not same: `%s` `%s`", a, tc.out)
			}
		})
	}
}

func TestResolveFunction(t *testing.T) {
	var tcs = []struct {
		input string

		prefix   string
		funcname string
		fields   []string
		returns  []string
	}{
		{
			input:   "__ssize_t (void *, char *, size_t)",
			prefix:  "",
			fields:  []string{"void *", "char *", "size_t"},
			returns: []string{"__ssize_t"},
		},
		{
			input:  "int (*)(sqlite3_vtab *, int, const char *, void (**)(sqlite3_context *, int, sqlite3_value **), void **)",
			prefix: "",
			fields: []string{
				"sqlite3_vtab *",
				"int",
				"const char *",
				"void (**)(sqlite3_context *, int, sqlite3_value **)",
				"void **"},
			returns: []string{"int"},
		},
		{
			input:   "void ( *(*)(int *, void *, char *))(void)",
			prefix:  "*",
			fields:  []string{"void"},
			returns: []string{"void (int *, void *, char *)"},
		},
		{
			input:   " void (*)(void)",
			prefix:  "",
			fields:  []string{"void"},
			returns: []string{"void"},
		},
		{
			input:   " int (*)(sqlite3_file *)",
			prefix:  "",
			fields:  []string{"sqlite3_file *"},
			returns: []string{"int"},
		},
		{
			input:   " int (*)(int)",
			prefix:  "",
			fields:  []string{"int"},
			returns: []string{"int"},
		},
		{
			input:   " int (*)(void *) ",
			prefix:  "",
			fields:  []string{"void *"},
			returns: []string{"int"},
		},
		{
			input:   " void (*)(sqlite3_context *, int, sqlite3_value **)",
			prefix:  "",
			fields:  []string{"sqlite3_context *", "int", "sqlite3_value **"},
			returns: []string{"void"},
		},
		{
			input:   "char *(*)( char *, ...)",
			prefix:  "",
			fields:  []string{"char *", "..."},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)( char *, struct __va_list_tag *)",
			prefix:  "",
			fields:  []string{"char *", "struct __va_list_tag *"},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)(const char *, ...)",
			prefix:  "",
			fields:  []string{"const char *", "..."},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)(ImportCtx *)",
			prefix:  "",
			fields:  []string{"ImportCtx *"},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)(int, char *, char *, ...)",
			prefix:  "",
			fields:  []string{"int", "char *", "char *", "..."},
			returns: []string{"char *"},
		},
		{
			input:   "const char *(*)(int)",
			prefix:  "",
			fields:  []string{"int"},
			returns: []string{"const char *"},
		},
		{
			input:   "const unsigned char *(*)(sqlite3_value *)",
			prefix:  "",
			fields:  []string{"sqlite3_value *"},
			returns: []string{"const unsigned char *"},
		},
		{
			input:   "int (*)(const char *, sqlite3 **)",
			prefix:  "",
			fields:  []string{"const char *", "sqlite3 **"},
			returns: []string{"int"},
		},
		{
			input:  "int (*)(fts5_api *, const char *, void *, fts5_extension_function, void (*)(void *))",
			prefix: "",
			fields: []string{"fts5_api *",
				"const char *",
				"void *",
				"fts5_extension_function",
				"void (*)(void *)"},
			returns: []string{"int"},
		},
		{
			input:  "int (*)(Fts5Context *, char *, int, void *, int (*)(void *, int, char *, int, int, int))",
			prefix: "",
			fields: []string{"Fts5Context *",
				"char *",
				"int",
				"void *",
				"int (*)(void *, int, char *, int, int, int)"},
			returns: []string{"int"},
		},
		{
			input:  "int (*)(sqlite3 *, char *, int, int, void *, void (*)(sqlite3_context *, int, sqlite3_value **), void (*)(sqlite3_context *, int, sqlite3_value **), void (*)(sqlite3_context *))",
			prefix: "",
			fields: []string{
				"sqlite3 *",
				"char *",
				"int",
				"int",
				"void *",
				"void (*)(sqlite3_context *, int, sqlite3_value **)",
				"void (*)(sqlite3_context *, int, sqlite3_value **)",
				"void (*)(sqlite3_context *)",
			},
			returns: []string{"int"},
		},
		{
			input:  "int (*)(sqlite3_vtab *, int, const char *, void (**)(sqlite3_context *, int, sqlite3_value **), void **)",
			prefix: "",
			fields: []string{
				"sqlite3_vtab *",
				"int",
				"const char *",
				"void (**)(sqlite3_context *, int, sqlite3_value **)",
				"void **"},
			returns: []string{"int"},
		},
		{
			input:   "void (*(int *, void *, const char *))(void)",
			prefix:  "",
			fields:  []string{"void"},
			returns: []string{"void (int *, void *, const char *)"},
		},
		{
			input:   "long (int, int)",
			prefix:  "",
			fields:  []string{"int", "int"},
			returns: []string{"long"},
		},
		{
			input:   "void (const char *, ...)",
			prefix:  "",
			fields:  []string{"const char *", "..."},
			returns: []string{"void"},
		},
		{
			input:   "void (void)",
			prefix:  "",
			fields:  []string{"void"},
			returns: []string{"void"},
		},
		{
			input:   "void ()",
			prefix:  "",
			fields:  []string{""},
			returns: []string{"void"},
		},
		{
			input:   "int (*)()",
			prefix:  "",
			fields:  []string{""},
			returns: []string{"int"},
		},
		{
			input:    "int ioctl(int , int , ... )",
			prefix:   "",
			funcname: "ioctl",
			fields:   []string{"int", "int", "..."},
			returns:  []string{"int"},
		},
		{
			input:    "speed_t cfgetospeed(const struct termios *)",
			prefix:   "",
			funcname: "cfgetospeed",
			fields:   []string{"const struct termios *"},
			returns:  []string{"speed_t"},
		},
		{
			input:    "void (*signal(int , void (*)(int)))(int)",
			prefix:   "*",
			funcname: "signal",
			fields:   []string{"int"},
			returns:  []string{"void(int,void(int))"},
		},
		{
			input:    "void ( * signal(int , void (*)(int)))(int)",
			prefix:   "*",
			funcname: "signal",
			fields:   []string{"int"},
			returns:  []string{"void(int,void(int))"},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Test %d : %s", i, tc.input), func(t *testing.T) {
			actualPrefix, actualName, actualField, actualReturn, err :=
				ParseFunction(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if actualPrefix != tc.prefix {
				t.Errorf("Prefix is not same.\nActual: %s\nExpected: %s\n",
					actualPrefix, tc.prefix)
			}
			if actualName != tc.funcname {
				t.Errorf("Names is not same.\nActual: %s\nExpected: %s\n",
					actualName, tc.funcname)
			}
			if len(actualField) != len(tc.fields) {
				a, _ := json.Marshal(actualField)
				f, _ := json.Marshal(tc.fields)
				t.Errorf("Size of field is not same.\nActual  : %s\nExpected: %s\n",
					string(a),
					string(f))
				return
			}
			if len(actualField) != len(tc.fields) {
				a, _ := json.Marshal(actualField)
				f, _ := json.Marshal(tc.fields)
				t.Errorf("Size of field is not same.\nActual  : %s\nExpected: %s\n",
					string(a),
					string(f))
				return
			}
			for i := range actualField {
				actualField[i] = strings.Replace(actualField[i], " ", "", -1)
				tc.fields[i] = strings.Replace(tc.fields[i], " ", "", -1)
				if actualField[i] != tc.fields[i] {
					t.Errorf("Not correct field: %v\nExpected: %v", actualField, tc.fields)
				}
			}
			if len(actualReturn) != len(tc.returns) {
				a, _ := json.Marshal(actualReturn)
				f, _ := json.Marshal(tc.returns)
				t.Errorf("Size of return field is not same.\nActual  : %s\nExpected: %s\n",
					string(a),
					string(f))
				return
			}
			if len(actualReturn) != len(tc.returns) {
				t.Errorf("Amount of return elements are different\nActual  : %v\nExpected: %v\n",
					actualReturn, tc.returns)
			}
			for i := range actualReturn {
				actualReturn[i] = strings.Replace(actualReturn[i], " ", "", -1)
				tc.returns[i] = strings.Replace(tc.returns[i], " ", "", -1)
				if actualReturn[i] != tc.returns[i] {
					t.Errorf("Not correct returns: %v\nExpected: %v", actualReturn, tc.returns)
				}
			}
		})
	}
}

// func TestGenerateCorrectType(t *testing.T) {
// tcs := []struct {
// inp string
// out string
// }{
// {
// inp: "union (anonymous union at tests/union.c:46:3)",
// out: "union __union_at_tests_union_c_46_3_",
// },
// {
// inp: " const struct (anonymous struct at /home/lepricon/go/src/github.com/tests/struct.c:282:18) [7]",
// out: "struct __struct_at__home_lepricon_go_src_github_com_tests_struct_c_282_18_ [7]",
// },
// }
//
// for i, tc := range tcs {
// t.Run(fmt.Sprintf("Test %d : %s", i, tc.inp), func(t *testing.T) {
// act := types.GenerateCorrectType(tc.inp)
// if act != tc.out {
// t.Errorf("Not correct result.\nExpected:%s\nActual:%s\n",
// tc.out, act)
// }
// })
// }
// }
