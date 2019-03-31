package transpiler

import (
	"bytes"
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

func transpileVAArgExpr(n *ast.VAArgExpr, p *program.Program) (
	expr goast.Expr,
	exprType string,
	preStmts []goast.Stmt,
	postStmts []goast.Stmt,
	err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileVAArgExpr. %v", err)
		}
	}()
	// -VAArgExpr 'int'
	//  `-ImplicitCastExpr 'struct __va_list_tag *' <ArrayToPointerDecay>
	//    `-DeclRefExpr 'va_list':'struct __va_list_tag [1]' lvalue Var 'ap' 'va_list':'struct __va_list_tag [1]'

	// Example of C code:
	//
	// va_list ap;
	// int i;
	// va_start(ap, num_args);
	// for(i = 0; i < num_args; i++) {
	//    val += va_arg(ap, int); // <<- This line
	// }
	// va_end(ap);

	var varName string
	if impl, ok := n.Children()[0].(*ast.ImplicitCastExpr); ok {
		if decl, ok := impl.Children()[0].(*ast.DeclRefExpr); ok {
			varName = decl.Name
		} else {
			err = fmt.Errorf("second node is not DeclRefExpr : %T",
				impl.Children()[0])
			return
		}
	} else {
		err = fmt.Errorf("first node is not ImplicitCastExpr : %T",
			n.Children()[0])
		return
	}

	varType, err := types.ResolveType(p, n.Type)
	if err != nil {
		return
	}

	type code struct {
		GoType string
		Name   string
	}

	src := `
package main

func main() {
	var {{.Name}} []interface{}
	var rr int = 10
	{{.Name}} = append({{.Name}}, rr)
	var c4goVaListPosition int
	/////////////////////////////////
	// Begin of needed code
	/////////////////////////////////
	rrr := func () (c4go_def {{ .GoType }}) {
		switch v := {{.Name}}[c4goVaListPosition].(type) {
			case int:
				return {{ .GoType }} (v)
			case int32:
				return {{ .GoType }} (v)
			case int64:
				return {{ .GoType }} (v)
			case float32: 
				return {{ .GoType }} (v)
			case float64: 
				return {{ .GoType }} (v)
		}
		return
	}()
	c4goVaListPosition++
	/////////////////////////////////
	// End of code
	/////////////////////////////////
	_ = rrr
}
`
	un := code{
		GoType: varType,
		Name:   varName,
	}

	if strings.Contains(varType, "[]") {
		src = `
package main

func main() {
	var {{.Name}} []interface{}
	var rr int = 10
	{{.Name}} = append({{.Name}}, rr)
	var c4goVaListPosition int
	/////////////////////////////////
	// Begin of needed code
	/////////////////////////////////
	rrr := {{.Name}}[c4goVaListPosition].({{ .GoType }})
	c4goVaListPosition++
	/////////////////////////////////
	// End of code
	/////////////////////////////////
	_ = rrr
}
`
	}

	tmpl := template.Must(template.New("").Parse(src))
	var source bytes.Buffer
	err = tmpl.Execute(&source, un)
	if err != nil {
		err = fmt.Errorf("cannot execute template \"%s\" for data %v : %v",
			source.String(), un, err)
		return
	}

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", source.String(), 0)
	if err != nil {
		err = fmt.Errorf("cannot parse source \"%s\" : %v",
			source.String(), err)
		return
	}

	expr = f.Decls[0].(*goast.FuncDecl).Body.List[4].(*goast.AssignStmt).Rhs[0]
	postStmts = []goast.Stmt{
		f.Decls[0].(*goast.FuncDecl).Body.List[5],
	}

	exprType = n.Type
	return
}

func getVaListStruct() string {
	return `

// va_list is C4GO implementation of va_list from "stdarg.h"
type va_list struct{
	position int
	slice    []interface{}
}

func create_va_list(list []interface{}) *va_list{
	return &va_list{
		position: 0,
		slice   : list,
	}
}

func va_start(v * va_list, count interface{}) {
	v.position = 0
}

func va_end(v * va_list) {
	// do nothing
}

func va_arg(v * va_list) interface{} {
	defer func(){
		 v.position++	
	}()
	return v.slice[v.position]
}

`
}

const create_va_list string = "create_va_list"

func VaListInit(name string) []goast.Decl {
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
				Names: []*goast.Ident{util.NewIdent(name)},
				Values: []goast.Expr{
					util.NewCallExpr(create_va_list, util.NewIdent("c4goArgs")),
				},
			},
		},
	}}
}
