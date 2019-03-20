package transpiler

import (
	"fmt"
	"go/token"
	"sort"

	goast "go/ast"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

const unsafeConvertFunctionName string = "c4goUnsafeConvert_"

func ConvertValueToPointer(nodes []ast.Node, p *program.Program) (expr goast.Expr, ok bool) {
	if len(nodes) != 1 {
		return nil, false
	}

	decl, ok := nodes[0].(*ast.DeclRefExpr)
	if !ok {
		return nil, false
	}

	if types.IsPointer(decl.Type, p) {
		return nil, false
	}

	// get base type if it typedef
	var td string = decl.Type
	for {
		if t, ok := p.TypedefType[td]; ok {
			td = t
			continue
		}
		break
	}

	resolvedType, err := types.ResolveType(p, td)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, decl))
		return
	}

	if !types.IsGoBaseType(resolvedType) {
		return nil, false
	}

	// can simplify
	p.UnsafeConvertValueToPointer[resolvedType] = true

	return util.NewCallExpr(fmt.Sprintf("%s%s", unsafeConvertFunctionName, resolvedType),
		&goast.UnaryExpr{
			Op: token.AND,
			X:  goast.NewIdent(decl.Name),
		}), true
}

func GetUnsafeConvertDecls(p *program.Program) {
	if len(p.UnsafeConvertValueToPointer) == 0 {
		return
	}

	p.AddImport("unsafe")

	var names []string
	for t := range p.UnsafeConvertValueToPointer {
		names = append(names, t)
	}
	sort.Sort(sort.StringSlice(names))

	for _, t := range names {
		functionName := fmt.Sprintf("%s%s", unsafeConvertFunctionName, t)
		varName := "c4go_name"
		p.File.Decls = append(p.File.Decls, &goast.FuncDecl{
			Name: goast.NewIdent(functionName),
			Type: &goast.FuncType{
				Params: &goast.FieldList{
					List: []*goast.Field{
						{
							Names: []*goast.Ident{goast.NewIdent(varName)},
							Type:  goast.NewIdent("*" + t),
						},
					},
				},
				Results: &goast.FieldList{
					List: []*goast.Field{
						{
							Type: &goast.ArrayType{
								Lbrack: 1,
								Elt:    goast.NewIdent(t),
							},
						},
					},
				},
			},
			Body: &goast.BlockStmt{
				List: []goast.Stmt{
					&goast.ReturnStmt{
						Results: []goast.Expr{
							&goast.SliceExpr{
								X: util.NewCallExpr(fmt.Sprintf("(*[1000000]%s)", t),
									util.NewCallExpr("unsafe.Pointer",
										goast.NewIdent(varName)),
								),
							},
						},
					},
				},
			},
		})
	}

	return
}
