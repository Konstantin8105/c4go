package transpiler

import (
	"bytes"
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"html/template"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
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
	// -VAArgExpr 0x3239438 <col:14, col:38> 'int'
	//  `-ImplicitCastExpr 0x3239420 <col:31> 'struct __va_list_tag *' <ArrayToPointerDecay>
	//    `-DeclRefExpr 0x32393e0 <col:31> 'va_list':'struct __va_list_tag [1]' lvalue Var 0x3239058 'ap' 'va_list':'struct __va_list_tag [1]'

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
	rrr := {{.Name}}[c4goVaListPosition].({{ .GoType }})
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
