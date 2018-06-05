// This file contains functions for transpiling declarations of variables and
// types. The usage of variables is handled in variables.go.

package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

/*
Example of AST for union without name inside struct:
-RecordDecl 0x40d41b0 <...> line:453:8 struct EmptyName definition
 |-RecordDecl 0x40d4260 <...> line:454:2 union definition
 | |-FieldDecl 0x40d4328 <...> col:8 referenced l1 'long'
 | `-FieldDecl 0x40d4388 <...> col:8 referenced l2 'long'
 |-FieldDecl 0x40d4420 <...> col:2 implicit referenced 'union EmptyName::(anonymous at struct.c:454:2)'
 |-IndirectFieldDecl 0x40d4478 <...> col:8 implicit l1 'long'
 | |-Field 0x40d4420 '' 'union EmptyName::(anonymous at /struct.c:454:2)'
 | `-Field 0x40d4328 'l1' 'long'
 `-IndirectFieldDecl 0x40d44c8 <...> col:8 implicit l2 'long'
   |-Field 0x40d4420 '' 'union EmptyName::(anonymous at /struct.c:454:2)'
   `-Field 0x40d4388 'l2' 'long'
*/

func newFunctionField(p *program.Program, name, cType string) (
	_ *goast.Field, err error) {
	if name == "" {
		err = fmt.Errorf("Name of function field cannot be empty")
		return
	}
	if !types.IsFunction(cType) {
		err = fmt.Errorf("Cannot create function field for type : %s", cType)
		return
	}

	// TODO : add err handling
	fieldType, _ := types.ResolveType(p, cType)

	return &goast.Field{
		Names: []*goast.Ident{util.NewIdent(name)},
		Type:  goast.NewIdent(fieldType),
	}, nil
}

func generateNameFieldDecl(t string) string {
	return "implicit_" + strings.Replace(t, " ", "S", -1)
}

func transpileFieldDecl(p *program.Program, n *ast.FieldDecl) (
	field *goast.Field, err error) {
	if types.IsFunction(n.Type) {
		field, err = newFunctionField(p, n.Name, n.Type)
		if err == nil {
			return
		}
	}

	if n.Name == "" {
		//&ast.FieldDecl{Addr:0x3157420, Pos:ast.Position{...}, Position2:"col:2", Name:"", Type:"union EmptyNameDD__at__home_lepricon_go_src_github_com_Konstantin8105_c4go_tests_struct_c_454_2_", Type2:"", Implicit:true, Referenced:true, ChildNodes:[]ast.Node{}}
		n.Name = generateNameFieldDecl(n.Type)
	}

	name := n.Name

	fieldType, err := types.ResolveType(p, n.Type)
	p.AddMessage(p.GenerateWarningMessage(err, n))

	// TODO: The name of a variable or field cannot be a reserved word
	// https://github.com/Konstantin8105/c4go/issues/83
	// Search for this issue in other areas of the codebase.
	if util.IsGoKeyword(name) {
		name += "_"
	}

	arrayType, arraySize := types.GetArrayTypeAndSize(n.Type)
	if arraySize != -1 {
		fieldType, err = types.ResolveType(p, arrayType)
		p.AddMessage(p.GenerateWarningMessage(err, n))
		fieldType = fmt.Sprintf("[%d]%s", arraySize, fieldType)
	}

	return &goast.Field{
		Names: []*goast.Ident{util.NewIdent(name)},
		Type:  util.NewTypeIdent(fieldType),
	}, nil
}

var ignoreRecordDecl = map[string]string{
	"struct __fsid_t":                             "/usr/include/x86_64-linux-gnu/bits/types.h",
	"union __union_at__usr_include_wchar_h_85_3_": "/usr/include/wchar.h",
	"struct __mbstate_t":                          "/usr/include/wchar.h",
	"struct _G_fpos_t":                            "/usr/include/_G_config.h",
	"struct _G_fpos64_t":                          "/usr/include/_G_config.h",
	"_IO_jump_t":                                  "/usr/include/libio.h",
	"_IO_marker":                                  "/usr/include/libio.h",
	"_IO_FILE":                                    "/usr/include/libio.h",
	"_IO_FILE_plus":                               "/usr/include/libio.h",
	"struct _IO_cookie_io_functions_t":            "/usr/include/libio.h",
	"_IO_cookie_file":                             "/usr/include/libio.h",
	"obstack":                                     "/usr/include/stdio.h",
	"timeval":                                     "/usr/include/x86_64-linux-gnu/bits/time.h",
	"timespec":                                    "/usr/include/time.h",
	"itimerspec":                                  "/usr/include/time.h",
	"sigevent":                                    "/usr/include/time.h",
	"__locale_struct":                             "/usr/include/xlocale.h",
	"exception":                                   "/usr/include/math.h",

	"wait":                    "/usr/include/x86_64-linux-gnu/bits/waitstatus.h",
	"struct __sigset_t":       "/usr/include/x86_64-linux-gnu/bits/sigset.h",
	"struct fd_set":           "/usr/include/x86_64-linux-gnu/sys/select.h",
	"pthread_attr_t":          "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"__pthread_internal_list": "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"__pthread_mutex_s":       "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_mutex_t":         "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_mutexattr_t":     "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_cond_t":          "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_condattr_t":      "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_rwlock_t":        "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_rwlockattr_t":    "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_barrier_t":       "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_barrierattr_t":   "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"random_data":             "/usr/include/stdlib.h",
	"drand48_data":            "/usr/include/stdlib.h",

	"siginfo_t":     "/usr/include/sys/signal.h",
	"__sigaction_u": "/usr/include/sys/signal.h",
	"__sigaction":   "/usr/include/sys/signal.h",
	"sigaction":     "/usr/include/sys/signal.h",
	"sig_t":         "/usr/include/sys/signal.h",
	"sigvec":        "/usr/include/sys/signal.h",
	"sigstack":      "/usr/include/sys/signal.h",
	"__siginfo":     "/usr/include/sys/signal.h",
}

func transpileRecordDecl(p *program.Program, n *ast.RecordDecl) (
	decls []goast.Decl, err error) {

	var addPackageUnsafe bool

	n.Name = types.GenerateCorrectType(n.Name)
	name := n.Name
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpileRecordDecl `%v`. %v",
				n.Name, err)
		} else {
			if _, ok := types.CStdStructType[name]; ok {
				// no need add struct for registrated C standart library
				decls = nil
				return
			}
			if !p.IncludeHeaderIsExists(n.Pos.File) {
				// no need add struct from C STD
				decls = nil
				return
			}
			if h, ok := ignoreRecordDecl[n.Name]; ok && p.IncludeHeaderIsExists(h) {
				decls = nil
				return
			}
			if addPackageUnsafe {
				p.AddImports("unsafe")
			}
			// Only for adding to ignore list
			// fmt.Printf("%40s:\t\"%s\",\n", "\""+n.Name+"\"", n.Pos.File)
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error - panic : %#v", r)
		}
	}()

	// ignore if haven`t definition
	if !n.Definition {
		return
	}

	if name == "" || p.IsTypeAlreadyDefined(name) {
		err = nil
		return
	}

	name = types.GenerateCorrectType(name)
	p.DefineType(name)
	defer func() {
		if err != nil {
			p.UndefineType(name)
		}
	}()

	// TODO: Some platform structs are ignored.
	// https://github.com/Konstantin8105/c4go/issues/85
	if name == "__locale_struct" ||
		name == "__sigaction" ||
		name == "sigaction" {
		err = nil
		return
	}

	var fields []*goast.Field

	// repair name for anonymous RecordDecl
	for pos := range n.Children() {
		if rec, ok := n.Children()[pos].(*ast.RecordDecl); ok && rec.Name == "" {
			if pos < len(n.Children()) {
				switch v := n.Children()[pos+1].(type) {
				case *ast.FieldDecl:
					rec.Name = types.GetBaseType(types.GenerateCorrectType(v.Type))
				default:
					p.AddMessage(p.GenerateWarningMessage(
						fmt.Errorf("Cannot find name for anon RecordDecl: %T",
							v), n))
					rec.Name = "UndefinedNameC2GO"
				}
			}
		}
	}

	for pos := range n.Children() {
		switch field := n.Children()[pos].(type) {
		case *ast.FieldDecl:
			field.Type = types.GenerateCorrectType(field.Type)
			field.Type2 = types.GenerateCorrectType(field.Type2)
			var f *goast.Field
			f, err = transpileFieldDecl(p, field)
			if err != nil {
				err = fmt.Errorf("cannot transpile field. %v", err)
				p.AddMessage(p.GenerateWarningMessage(err, field))
				// TODO ignore error
				// return
				err = nil
			} else {
				// ignore fields without name
				if len(f.Names) != 1 {
					p.AddMessage(p.GenerateWarningMessage(
						fmt.Errorf("Ignore FieldDecl with more then 1 names"+
							" in RecordDecl : `%v`", n.Name), n))
					continue
				}
				if f.Names[0].Name == "" {
					p.AddMessage(p.GenerateWarningMessage(
						fmt.Errorf("Ignore FieldDecl without name "+
							" in RecordDecl : `%v`", n.Name), n))
					continue
				}
				// remove dublicates of fields
				var isDublicate bool
				for i := range fields {
					if fields[i].Names[0].Name == f.Names[0].Name {
						isDublicate = true
					}
				}
				if isDublicate {
					f.Names[0].Name += strconv.Itoa(pos)
				}
				fields = append(fields, f)
			}

		case *ast.IndirectFieldDecl:
			// ignore

		case *ast.AlignedAttr:
			// ignore

		case *ast.PackedAttr:
			// ignore

		case *ast.MaxFieldAlignmentAttr:
			// ignore

		case *ast.FullComment:
			// We haven't Go ast struct for easy inject a comments.
			// All comments are added like CommentsGroup.
			// So, we can ignore that comment, because all comments
			// will be added by another way.

		case *ast.TransparentUnionAttr:
			// Don't do anythink
			// Example of AST:
			// |-RecordDecl 0x3632d78 <...> line:67:9 union definition
			// | |-TransparentUnionAttr 0x3633050 <...>
			// | |-FieldDecl 0x3632ed0 <...> col:17 __uptr 'union wait *'
			// | `-FieldDecl 0x3632f60 <...> col:10 __iptr 'int *'
			// |-TypedefDecl 0x3633000 <...> col:5 __WAIT_STATUS 'union __WAIT_STATUS':'__WAIT_STATUS'
			// | `-ElaboratedType 0x3632fb0 'union __WAIT_STATUS' sugar
			// |   `-RecordType 0x3632e00 '__WAIT_STATUS'
			// |     `-Record 0x3632d78 ''

		default:
			// For case anonymous enum:

			// |-EnumDecl 0x26c3970 <...> line:77:5
			// | `-EnumConstantDecl 0x26c3a50 <...> col:9 referenced SWE_ENUM_THREE 'int'
			// |   `-IntegerLiteral 0x26c3a30 <...> 'int' 3
			// |-FieldDecl 0x26c3af0 <...> col:7 EnumThree 'enum (anonymous enum at ...
			if eDecl, ok := field.(*ast.EnumDecl); ok && eDecl.Name == "" {
				if pos+1 <= len(n.Children())-1 {
					if f, ok := n.Children()[pos+1].(*ast.FieldDecl); ok {
						n.Children()[pos].(*ast.EnumDecl).Name = f.Type
					}
				}
			}

			// default
			var declsIn []goast.Decl
			declsIn, err = transpileToNode(field, p)
			if err != nil {
				err = fmt.Errorf("Cannot transpile %T", field)
				p.AddMessage(p.GenerateWarningMessage(err, field))
				return
			}
			decls = append(decls, declsIn...)
		}
	}

	s, err := program.NewStruct(n)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
		return
	}
	switch s.Type {
	case program.UnionType:
		if strings.HasPrefix(s.Name, "union ") {
			p.Structs[s.Name] = s
			defer func() {
				if err != nil {
					delete(p.Structs, s.Name)
					p.UndefineType(s.Name)
				}
			}()
		} else {
			p.Unions["union "+s.Name] = s
			defer func() {
				if err != nil {
					delete(p.Structs, "union "+s.Name)
					p.UndefineType("union " + s.Name)
				}
			}()
		}

	case program.StructType:
		if strings.HasPrefix(s.Name, "struct ") {
			p.Structs[s.Name] = s
			defer func() {
				if err != nil {
					delete(p.Structs, s.Name)
					p.UndefineType(s.Name)
				}
			}()
		} else {
			p.Structs["struct "+s.Name] = s
			defer func() {
				if err != nil {
					delete(p.Structs, "struct "+s.Name)
					p.UndefineType("struct " + s.Name)
				}
			}()
		}

	default:
		err = fmt.Errorf("Undefine type of struct : %v", s.Type)
		return
	}

	name = strings.TrimPrefix(name, "struct ")
	name = strings.TrimPrefix(name, "union ")

	var d []goast.Decl
	switch s.Type {
	case program.UnionType:
		// Union size
		var size int
		size, err = types.SizeOf(p, "union "+name)

		// In normal case no error is returned,
		if err != nil {
			// but if we catch one, send it as a warning
			err = fmt.Errorf("could not determine the size of type `union %s`"+
				" for that reason: %s", name, err)
			return
		}
		// So, we got size, then
		// Add imports needed
		addPackageUnsafe = true

		// Declaration for implementing union type
		d, err = transpileUnion(name, size, fields)
		if err != nil {
			return nil, err
		}

	case program.StructType:
		d = append(d, &goast.GenDecl{
			Tok: token.TYPE,
			Specs: []goast.Spec{
				&goast.TypeSpec{
					Name: util.NewIdent(name),
					Type: &goast.StructType{
						Fields: &goast.FieldList{
							List: fields,
						},
					},
				},
			},
		})

	default:
		err = fmt.Errorf("Undefine type of struct : %v", s.Type)
		return
	}

	decls = append(decls, d...)

	return
}

func transpileCXXRecordDecl(p *program.Program, n *ast.RecordDecl) (
	decls []goast.Decl, err error) {

	n.Name = types.GenerateCorrectType(n.Name)
	name := n.Name

	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpileCXXRecordDecl : `%v`. %v",
				n.Name, err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error - panic : %#v", r)
		}
	}()

	// ignore if haven`t definition
	if !n.Definition {
		return
	}

	if name == "" || p.IsTypeAlreadyDefined(name) {
		err = nil
		return
	}

	p.DefineType(n.Kind + " " + name)
	defer func() {
		if err != nil {
			p.UndefineType(n.Kind + " " + name)
		}
	}()

	var fields []*goast.Field
	for _, v := range n.Children() {
		switch v := v.(type) {
		case *ast.CXXRecordDecl:
			// ignore

		case *ast.FieldDecl:
			var f *goast.Field
			f, err = transpileFieldDecl(p, v)
			if err != nil {
				return
			}
			fields = append(fields, f)

		default:
			p.AddMessage(p.GenerateWarningMessage(
				fmt.Errorf("Cannot transpilation field in CXXRecordDecl : %T", v), n))
		}
	}

	return []goast.Decl{&goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: util.NewIdent(name),
				Type: &goast.StructType{
					Fields: &goast.FieldList{
						List: fields,
					},
				},
			},
		},
	}}, nil
}

var ignoreTypedef = map[string]string{
	"__u_char":   "bits/types.h",
	"__u_short":  "bits/types.h",
	"__u_int":    "bits/types.h",
	"__u_long":   "bits/types.h",
	"__int8_t":   "bits/types.h",
	"__uint8_t":  "bits/types.h",
	"__int16_t":  "bits/types.h",
	"__uint16_t": "bits/types.h",
	"__int32_t":  "bits/types.h",
	"__uint32_t": "bits/types.h",
	"__int64_t":  "bits/types.h",
	"__uint64_t": "bits/types.h",
	"__quad_t":   "bits/types.h",
	"__u_quad_t": "bits/types.h",
	"__dev_t":    "bits/types.h",
	"__uid_t":    "bits/types.h",
	"__gid_t":    "bits/types.h",
	"__ino_t":    "bits/types.h",
	"__ino64_t":  "bits/types.h",
	"__mode_t":   "bits/types.h",
	"__nlink_t":  "bits/types.h",
	"__off_t":    "bits/types.h",
	"__off64_t":  "bits/types.h",
	"__pid_t":    "bits/types.h",
	"__fsid_t":   "bits/types.h",
	"__clock_t":  "bits/types.h",
	"__rlim_t":   "bits/types.h",
	"__rlim64_t": "bits/types.h",
	"__id_t":     "bits/types.h",
	// "__time_t":          "bits/types.h", // need for test time.c
	"__useconds_t":      "bits/types.h",
	"__suseconds_t":     "bits/types.h",
	"__daddr_t":         "bits/types.h",
	"__key_t":           "bits/types.h",
	"__clockid_t":       "bits/types.h",
	"__timer_t":         "bits/types.h",
	"__blksize_t":       "bits/types.h",
	"__blkcnt_t":        "bits/types.h",
	"__blkcnt64_t":      "bits/types.h",
	"__fsblkcnt_t":      "bits/types.h",
	"__fsblkcnt64_t":    "bits/types.h",
	"__fsfilcnt_t":      "bits/types.h",
	"__fsfilcnt64_t":    "bits/types.h",
	"__fsword_t":        "bits/types.h",
	"__ssize_t":         "bits/types.h",
	"__syscall_slong_t": "bits/types.h",
	"__syscall_ulong_t": "bits/types.h",
	"__loff_t":          "bits/types.h",
	"__qaddr_t":         "bits/types.h",
	"__caddr_t":         "bits/types.h",
	"__intptr_t":        "bits/types.h",
	"__socklen_t":       "bits/types.h",

	"u_char":      "/usr/include/x86_64-linux-gnu/sys/types.h",
	"u_short":     "/usr/include/x86_64-linux-gnu/sys/types.h",
	"u_int":       "/usr/include/x86_64-linux-gnu/sys/types.h",
	"u_long":      "/usr/include/x86_64-linux-gnu/sys/types.h",
	"quad_t":      "/usr/include/x86_64-linux-gnu/sys/types.h",
	"u_quad_t":    "/usr/include/x86_64-linux-gnu/sys/types.h",
	"fsid_t":      "/usr/include/x86_64-linux-gnu/sys/types.h",
	"loff_t":      "/usr/include/x86_64-linux-gnu/sys/types.h",
	"ino_t":       "/usr/include/x86_64-linux-gnu/sys/types.h",
	"ino64_t":     "/usr/include/x86_64-linux-gnu/sys/types.h",
	"dev_t":       "/usr/include/x86_64-linux-gnu/sys/types.h",
	"gid_t":       "/usr/include/x86_64-linux-gnu/sys/types.h",
	"mode_t":      "/usr/include/x86_64-linux-gnu/sys/types.h",
	"nlink_t":     "/usr/include/x86_64-linux-gnu/sys/types.h",
	"uid_t":       "/usr/include/x86_64-linux-gnu/sys/types.h",
	"id_t":        "/usr/include/x86_64-linux-gnu/sys/types.h",
	"daddr_t":     "/usr/include/x86_64-linux-gnu/sys/types.h",
	"caddr_t":     "/usr/include/x86_64-linux-gnu/sys/types.h",
	"key_t":       "/usr/include/x86_64-linux-gnu/sys/types.h",
	"useconds_t":  "/usr/include/x86_64-linux-gnu/sys/types.h",
	"suseconds_t": "/usr/include/x86_64-linux-gnu/sys/types.h",
	"ulong":       "/usr/include/x86_64-linux-gnu/sys/types.h",
	"ushort":      "/usr/include/x86_64-linux-gnu/sys/types.h",
	"uint":        "/usr/include/x86_64-linux-gnu/sys/types.h",
	"int8_t":      "/usr/include/x86_64-linux-gnu/sys/types.h",
	"int16_t":     "/usr/include/x86_64-linux-gnu/sys/types.h",
	// "int32_t":     "/usr/include/x86_64-linux-gnu/sys/types.h",// need for struct
	"int64_t":    "/usr/include/x86_64-linux-gnu/sys/types.h",
	"u_int8_t":   "/usr/include/x86_64-linux-gnu/sys/types.h",
	"u_int16_t":  "/usr/include/x86_64-linux-gnu/sys/types.h",
	"u_int32_t":  "/usr/include/x86_64-linux-gnu/sys/types.h",
	"u_int64_t":  "/usr/include/x86_64-linux-gnu/sys/types.h",
	"register_t": "/usr/include/x86_64-linux-gnu/sys/types.h",

	"blksize_t":    "/usr/include/x86_64-linux-gnu/sys/types.h",
	"blkcnt_t":     "/usr/include/x86_64-linux-gnu/sys/types.h",
	"fsblkcnt_t":   "/usr/include/x86_64-linux-gnu/sys/types.h",
	"fsfilcnt_t":   "/usr/include/x86_64-linux-gnu/sys/types.h",
	"blkcnt64_t":   "/usr/include/x86_64-linux-gnu/sys/types.h",
	"fsblkcnt64_t": "/usr/include/x86_64-linux-gnu/sys/types.h",
	"fsfilcnt64_t": "/usr/include/x86_64-linux-gnu/sys/types.h",

	"clock_t":   "include/time.h",
	"time_t":    "include/time.h",
	"clockid_t": "include/time.h",
	"timer_t":   "include/time.h",
	"pid_t":     "include/time.h",

	"FILE":     "/usr/include/stdio.h",
	"__FILE":   "/usr/include/stdio.h",
	"off_t":    "/usr/include/stdio.h",
	"off64_t":  "/usr/include/stdio.h",
	"ssize_t":  "/usr/include/stdio.h",
	"fpos_t":   "/usr/include/stdio.h",
	"fpos64_t": "/usr/include/stdio.h",

	"_IO_lock_t":                "/usr/include/libio.h",
	"_IO_FILE":                  "/usr/include/libio.h",
	"__io_read_fn":              "/usr/include/libio.h",
	"__io_write_fn":             "/usr/include/libio.h",
	"__io_seek_fn":              "/usr/include/libio.h",
	"__io_close_fn":             "/usr/include/libio.h",
	"cookie_read_function_t":    "/usr/include/libio.h",
	"cookie_write_function_t":   "/usr/include/libio.h",
	"cookie_seek_function_t":    "/usr/include/libio.h",
	"cookie_close_function_t":   "/usr/include/libio.h",
	"_IO_cookie_io_functions_t": "/usr/include/libio.h",
	"cookie_io_functions_t":     "/usr/include/libio.h",

	"__mbstate_t":           "/usr/include/wchar.h",
	"_G_fpos_t":             "/usr/include/_G_config.h",
	"_G_fpos64_t":           "/usr/include/_G_config.h",
	"va_list":               "/usr/lib/llvm-4.0/bin/../lib/clang/4.0.0/include/stdarg.h",
	"__gnuc_va_list":        "/usr/lib/llvm-4.0/bin/../lib/clang/4.0.0/include/stdarg.h",
	"wchar_t":               "/usr/lib/llvm-4.0/bin/../lib/clang/4.0.0/include/stddef.h",
	"idtype_t":              "/usr/include/x86_64-linux-gnu/bits/waitflags.h",
	"__WAIT_STATUS":         "/usr/include/stdlib.h",
	"int32_t":               "/usr/include/x86_64-linux-gnu/sys/types.h",
	"__sig_atomic_t":        "/usr/include/x86_64-linux-gnu/bits/sigset.h",
	"__sigset_t":            "/usr/include/x86_64-linux-gnu/bits/sigset.h",
	"sigset_t":              "/usr/include/x86_64-linux-gnu/sys/select.h",
	"__fd_mask":             "/usr/include/x86_64-linux-gnu/sys/select.h",
	"fd_set":                "/usr/include/x86_64-linux-gnu/sys/select.h",
	"fd_mask":               "/usr/include/x86_64-linux-gnu/sys/select.h",
	"pthread_t":             "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_attr_t":        "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"__pthread_list_t":      "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_mutex_t":       "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_mutexattr_t":   "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_cond_t":        "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_condattr_t":    "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_key_t":         "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_once_t":        "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_rwlock_t":      "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_rwlockattr_t":  "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_spinlock_t":    "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_barrier_t":     "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"pthread_barrierattr_t": "/usr/include/x86_64-linux-gnu/bits/pthreadtypes.h",
	"__compar_fn_t":         "/usr/include/stdlib.h",
	"comparison_fn_t":       "/usr/include/stdlib.h",
	"__compar_d_fn_t":       "/usr/include/stdlib.h",
	"float_t":               "/usr/include/x86_64-linux-gnu/bits/mathdef.h",
	"double_t":              "/usr/include/x86_64-linux-gnu/bits/mathdef.h",
	"_LIB_VERSION_TYPE":     "/usr/include/math.h",

	"intptr_t":  "/usr/include/unistd.h",
	"socklen_t": "/usr/include/unistd.h",

	"siginfo_t": "/usr/include/sys/signal.h",
}

func transpileTypedefDecl(p *program.Program, n *ast.TypedefDecl) (
	decls []goast.Decl, err error) {

	// implicit code from clang at the head of each clang AST tree
	if n.IsImplicit && n.Pos.File == ast.PositionBuiltIn {
		return
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile Typedef Decl : err = %v", err)
		} else {
			if !p.IncludeHeaderIsExists(n.Pos.File) {
				// no need add struct from C STD
				decls = nil
				return
			}
			if h, ok := ignoreTypedef[n.Name]; ok && p.IncludeHeaderIsExists(h) {
				decls = nil
				return
			}
			// Only for adding to ignore list
			// fmt.Printf("%40s:\t\"%s\",\n", "\""+n.Name+"\"", n.Pos.File)
		}
	}()
	n.Name = types.CleanCType(types.GenerateCorrectType(n.Name))
	n.Type = types.CleanCType(types.GenerateCorrectType(n.Type))
	n.Type2 = types.CleanCType(types.GenerateCorrectType(n.Type2))
	name := n.Name

	if "struct "+n.Name == n.Type || "union "+n.Name == n.Type {
		p.TypedefType[n.Name] = n.Type
		return
	}

	if types.IsFunction(n.Type) {
		var field *goast.Field
		field, err = newFunctionField(p, n.Name, n.Type)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		} else {
			// registration type
			p.TypedefType[n.Name] = n.Type

			decls = append(decls, &goast.GenDecl{
				Tok: token.TYPE,
				Specs: []goast.Spec{
					&goast.TypeSpec{
						Name: util.NewIdent(name),
						Type: field.Type,
					},
				},
			})
			err = nil
			return
		}
	}

	// added for support "typedef enum {...} dd" with empty name of struct
	// Result in Go: "type dd int"
	if strings.Contains(n.Type, "enum") {
		// Registration new type in program.Program
		if !p.IsTypeAlreadyDefined(n.Name) {
			p.DefineType(n.Name)
			p.EnumTypedefName[n.Name] = true
		}
		decls = append(decls, &goast.GenDecl{
			Tok: token.TYPE,
			Specs: []goast.Spec{
				&goast.TypeSpec{
					Name: util.NewIdent(name),
					Type: util.NewTypeIdent("int"),
				},
			},
		})
		err = nil
		return
	}

	if p.IsTypeAlreadyDefined(name) {
		err = nil
		return
	}

	p.DefineType(name)

	resolvedType, err := types.ResolveType(p, n.Type)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
	}

	// There is a case where the name of the type is also the definition,
	// like:
	//
	//     type _RuneEntry _RuneEntry
	//
	// This of course is impossible and will cause the Go not to compile.
	// It itself is caused by lack of understanding (at this time) about
	// certain scenarios that types are defined as. The above example comes
	// from:
	//
	//     typedef struct {
	//        // ... some fields
	//     } _RuneEntry;
	//
	// Until which time that we actually need this to work I am going to
	// suppress these.
	if name == resolvedType {
		err = nil
		return
	}

	if name == "__darwin_ct_rune_t" {
		resolvedType = p.ImportType("github.com/Konstantin8105/c4go/darwin.CtRuneT")
	}

	if name == "div_t" || name == "ldiv_t" || name == "lldiv_t" {
		intType := "int"
		if name == "ldiv_t" {
			intType = "long int"
		} else if name == "lldiv_t" {
			intType = "long long int"
		}

		// I don't know to extract the correct fields from the typedef to create
		// the internal definition. This is used in the noarch package
		// (stdio.go).
		//
		// The name of the struct is not prefixed with "struct " because it is a
		// typedef.
		p.Structs[name] = &program.Struct{
			Name: name,
			Type: program.StructType,
			Fields: map[string]interface{}{
				"quot": intType,
				"rem":  intType,
			},
		}
	}

	err = nil
	if resolvedType == "" {
		resolvedType = "interface{}"
	}

	if v, ok := p.Structs["struct "+resolvedType]; ok {
		// Registration "typedef struct" with non-empty name of struct
		p.Structs["struct "+name] = v
	} else if v, ok := p.EnumConstantToEnum["enum "+resolvedType]; ok {
		// Registration "enum constants"
		p.EnumConstantToEnum["enum "+resolvedType] = v
	} else {
		// Registration "typedef type type2"
		p.TypedefType[n.Name] = n.Type
	}

	decls = append(decls, &goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: util.NewIdent(name),
				Type: util.NewTypeIdent(resolvedType),
			},
		},
	})

	return
}

func transpileVarDecl(p *program.Program, n *ast.VarDecl) (
	decls []goast.Decl, theType string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileVarDecl : err = %v", err)
		}
	}()

	n.Name = types.GenerateCorrectType(n.Name)
	n.Type = types.GenerateCorrectType(n.Type)
	n.Type2 = types.GenerateCorrectType(n.Type2)

	// There may be some startup code for this global variable.
	if p.Function == nil {
		name := n.Name
		switch name {
		// Below are for macOS.
		case "__stdinp", "__stdoutp", "__stderrp":
			theType = "*noarch.File"
			p.AddImport("github.com/Konstantin8105/c4go/noarch")
			p.AppendStartupExpr(
				util.NewBinaryExpr(
					goast.NewIdent(name),
					token.ASSIGN,
					util.NewTypeIdent(
						"noarch."+util.Ucfirst(name[2:len(name)-1])),
					"*noarch.File",
					true,
				),
			)
			return []goast.Decl{&goast.GenDecl{
				Tok: token.VAR,
				Specs: []goast.Spec{&goast.ValueSpec{
					Names: []*goast.Ident{{Name: name}},
					Type:  util.NewTypeIdent(theType),
					Doc:   p.GetMessageComments(),
				}},
			}}, "", nil

		// Below are for linux.
		case "stdout", "stdin", "stderr":
			theType = "*noarch.File"
			p.AddImport("github.com/Konstantin8105/c4go/noarch")
			p.AppendStartupExpr(
				util.NewBinaryExpr(
					goast.NewIdent(name),
					token.ASSIGN,
					util.NewTypeIdent("noarch."+util.Ucfirst(name)),
					theType,
					true,
				),
			)
			return []goast.Decl{&goast.GenDecl{
				Tok: token.VAR,
				Specs: []goast.Spec{&goast.ValueSpec{
					Names: []*goast.Ident{{Name: name}},
					Type:  util.NewTypeIdent(theType),
				}},
				Doc: p.GetMessageComments(),
			}}, "", nil

		default:
			// No init needed.
		}
	}

	// Ignore extern as there is no analogy for Go right now.
	if n.IsExtern && len(n.ChildNodes) == 0 {
		return
	}

	if strings.Contains(n.Type, "va_list") &&
		strings.Contains(n.Type2, "va_list_tag") {
		// variable for va_list. see "variadic function"
		// header : <stdarg.h>
		// Example :
		// DeclStmt 0x2fd87e0 <line:442:2, col:14>
		// `-VarDecl 0x2fd8780 <col:2, col:10> col:10 used args 'va_list':'struct __va_list_tag [1]'
		// Result:
		// ... - convert to - c4goArgs ...interface{}
		// var args = c4goArgs
		return []goast.Decl{&goast.GenDecl{
			Tok: token.VAR,
			Specs: []goast.Spec{
				&goast.ValueSpec{
					Names:  []*goast.Ident{util.NewIdent(n.Name)},
					Values: []goast.Expr{util.NewIdent("c4goArgs")},
				},
			},
		}}, "", nil
	}

	/*
		Example of DeclStmt for C code:
		void * a = NULL;
		void(*t)(void) = a;
		Example of AST:
		`-VarDecl 0x365fea8 <col:3, col:20> col:9 used t 'void (*)(void)' cinit
		  `-ImplicitCastExpr 0x365ff48 <col:20> 'void (*)(void)' <BitCast>
		    `-ImplicitCastExpr 0x365ff30 <col:20> 'void *' <LValueToRValue>
		      `-DeclRefExpr 0x365ff08 <col:20> 'void *' lvalue Var 0x365f8c8 'r' 'void *'
	*/

	if len(n.Children()) > 0 {
		if v, ok := (n.Children()[0]).(*ast.ImplicitCastExpr); ok {
			if len(v.Type) > 0 {
				// Is it function ?
				if types.IsFunction(v.Type) {
					var prefix string
					var fields, returns []string
					prefix, fields, returns, err = types.SeparateFunction(p, v.Type)
					if err != nil {
						err = fmt.Errorf("Cannot resolve function : %v", err)
						return
					}
					if len(prefix) != 0 {
						p.AddMessage(p.GenerateWarningMessage(
							fmt.Errorf("Prefix is not used : `%v`", prefix), n))
					}
					functionType := GenerateFuncType(fields, returns)
					nameVar1 := n.Name

					if vv, ok := v.Children()[0].(*ast.ImplicitCastExpr); ok {
						if decl, ok := vv.Children()[0].(*ast.DeclRefExpr); ok {
							nameVar2 := decl.Name

							return []goast.Decl{&goast.GenDecl{
								Tok: token.VAR,
								Specs: []goast.Spec{&goast.ValueSpec{
									Names: []*goast.Ident{{Name: nameVar1}},
									Type:  functionType,
									Values: []goast.Expr{&goast.TypeAssertExpr{
										X:    &goast.Ident{Name: nameVar2},
										Type: functionType,
									}},
									Doc: p.GetMessageComments(),
								},
								}}}, "", nil
						}
					}
				}
			}
		}
	}

	theType = n.Type

	p.GlobalVariables[n.Name] = theType

	name := n.Name
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	// TODO: Some platform structs are ignored.
	// https://github.com/Konstantin8105/c4go/issues/85
	if name == "_LIB_VERSION" ||
		name == "_IO_2_1_stdin_" ||
		name == "_IO_2_1_stdout_" ||
		name == "_IO_2_1_stderr_" ||
		name == "_DefaultRuneLocale" ||
		name == "_CurrentRuneLocale" {
		theType = "unknown10"
		return
	}

	defaultValue, _, newPre, newPost, err := getDefaultValueForVar(p, n)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
		err = nil // Error is ignored
	}
	// for ignore zero value. example:
	// int i = 0;
	// tranpile to:
	// var i int // but not "var i int = 0"
	if len(defaultValue) == 1 && defaultValue[0] != nil {
		if bl, ok := defaultValue[0].(*goast.BasicLit); ok {
			if bl.Kind == token.INT && bl.Value == "0" {
				defaultValue = nil
			}
			if bl.Kind == token.FLOAT && bl.Value == "0" {
				defaultValue = nil
			}
		} else if call, ok := defaultValue[0].(*goast.CallExpr); ok {
			if len(call.Args) == 1 {
				if bl, ok := call.Args[0].(*goast.BasicLit); ok {
					if bl.Kind == token.INT && bl.Value == "0" {
						defaultValue = nil
					}
					if bl.Kind == token.FLOAT && bl.Value == "0" {
						defaultValue = nil
					}
				}
			}
		} else if ind, ok := defaultValue[0].(*goast.Ident); ok {
			if ind.Name == "nil" {
				defaultValue = nil
			}
		}
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// Allocate slice so that it operates like a fixed size array.
	arrayType, arraySize := types.GetArrayTypeAndSize(n.Type)

	if arraySize != -1 && defaultValue == nil {
		var goArrayType string
		goArrayType, err = types.ResolveType(p, arrayType)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
			err = nil // Error is ignored
		}

		defaultValue = []goast.Expr{
			util.NewCallExpr(
				"make",
				&goast.ArrayType{
					Elt: util.NewTypeIdent(goArrayType),
				},
				util.NewIntLit(arraySize),
				// If len and capacity is same, then
				// capacity is not need
				// util.NewIntLit(arraySize),
			),
		}
	}

	if len(preStmts) != 0 || len(postStmts) != 0 {
		p.AddMessage(p.GenerateWarningMessage(
			fmt.Errorf("Not acceptable length of Stmt : pre(%d), post(%d)",
				len(preStmts), len(postStmts)), n))
	}

	theType, err = types.ResolveType(p, n.Type)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
		err = nil // Error is ignored
		theType = "UnknownType"
	}
	typeResult := util.NewTypeIdent(theType)

	return []goast.Decl{&goast.GenDecl{
		Tok: token.VAR,
		Specs: []goast.Spec{
			&goast.ValueSpec{
				Names:  []*goast.Ident{util.NewIdent(n.Name)},
				Type:   typeResult,
				Values: defaultValue,
				Doc:    p.GetMessageComments(),
			},
		},
	}}, "", nil
}
