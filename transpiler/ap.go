package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"strings"

	"github.com/Konstantin8105/c4go/util"
)

// GetUintptrForSlice - return uintptr for slice
// Example : int64(uintptr(unsafe.Pointer((*(**int)(unsafe.Pointer(& ...slice... )))))))
func GetUintptrForSlice(expr goast.Expr, sizeof int) (goast.Expr, string) {
	returnType := "long long"

	if _, ok := expr.(*goast.SelectorExpr); ok {
		expr = &goast.IndexExpr{
			X:     expr,
			Index: goast.NewIdent("0"),
		}
	}

	if sl, ok := expr.(*goast.SliceExpr); ok {
		// from :
		//
		// 88  0: *ast.SliceExpr {
		// 89  .  X: *ast.Ident {
		// 91  .  .  Name: "b"
		// 93  .  }
		// 95  .  Low: *ast.BasicLit { ... }
		// 99  .  }
		// 102  }
		//
		// to:
		//
		// 0  *ast.IndexExpr {
		// 1  .  X: *ast.Ident {
		// 3  .  .  Name: "b"
		// 4  .  }
		// 6  .  Index: *ast.BasicLit { ... }
		// 12  }
		if sl.Low == nil {
			sl.Low = goast.NewIdent("0")
		}
		util.PanicIfNil(sl.X, "slice is nil")
		util.PanicIfNil(sl.Low, "slice low is nil")
		expr = &goast.IndexExpr{
			X:     sl.X,
			Index: sl.Low,
		}
	}

	if sl, ok := expr.(*goast.SliceExpr); ok {
		if c, ok := sl.X.(*goast.CallExpr); ok {
			if fin, ok := c.Fun.(*goast.Ident); ok && strings.Contains(fin.Name, "100000") {
				if len(c.Args) == 1 {
					if cc, ok := c.Args[0].(*goast.CallExpr); ok {
						if fin, ok := cc.Fun.(*goast.Ident); ok && strings.Contains(fin.Name, "unsafe.Pointer") {
							if len(cc.Args) == 1 {
								if un, ok := cc.Args[0].(*goast.UnaryExpr); ok && un.Op == token.AND {
									expr = un.X
								}
							}
						}
					}
				}
			}
		}
	}

	if _, ok := expr.(*goast.CallExpr); ok {
		name := "c4go_temp_name"
		expr = util.NewAnonymousFunction(
			// body
			[]goast.Stmt{
				&goast.ExprStmt{
					X: &goast.BinaryExpr{
						X:  goast.NewIdent(name),
						Op: token.DEFINE,
						Y:  expr,
					},
				},
			},
			// defer
			nil,
			// returnValue
			util.NewCallExpr("int64", util.NewCallExpr("uintptr", util.NewCallExpr("unsafe.Pointer",
				&goast.StarExpr{
					Star: 1,
					X: &goast.CallExpr{
						Fun:    goast.NewIdent("(**byte)"),
						Lparen: 1,
						Args: []goast.Expr{&goast.CallExpr{
							Fun:    goast.NewIdent("unsafe.Pointer"),
							Lparen: 1,
							Args: []goast.Expr{
								&goast.UnaryExpr{
									Op: token.AND,
									X:  goast.NewIdent(name),
								},
							},
						}},
					},
				},
			))),
			// returnType
			"int64",
		)
		return expr, returnType
	}

	return &goast.BinaryExpr{
		X: util.NewCallExpr("int64", util.NewCallExpr("uintptr", util.NewCallExpr("unsafe.Pointer",
			&goast.StarExpr{
				Star: 1,
				X: &goast.CallExpr{
					Fun:    goast.NewIdent("(**byte)"),
					Lparen: 1,
					Args: []goast.Expr{&goast.CallExpr{
						Fun:    goast.NewIdent("unsafe.Pointer"),
						Lparen: 1,
						Args: []goast.Expr{&goast.UnaryExpr{
							Op: token.AND,
							X:  expr,
						}},
					}},
				},
			},
		))),
		Op: token.QUO,
		Y:  util.NewCallExpr("int64", goast.NewIdent(fmt.Sprintf("%d", sizeof))),
	}, returnType
}

// CreateSliceFromReference - create a slice, like :
// (*[1]int)(unsafe.Pointer(&a))[:]
func CreateSliceFromReference(goType string, expr goast.Expr) goast.Expr {
	// If the Go type is blank it means that the C type is 'void'.
	if goType == "" {
		goType = "interface{}"
	}

	// This is a hack to convert a reference to a variable into a slice that
	// points to the same location. It will look similar to:
	//
	//     (*[1]int)(unsafe.Pointer(&a))[:]
	//
	// You must always call this Go before using CreateSliceFromReference:
	//
	//     p.AddImport("unsafe")
	//
	return &goast.SliceExpr{
		X: util.NewCallExpr(fmt.Sprintf("(*[100000000]%s)", goType),
			util.NewCallExpr("unsafe.Pointer",
				&goast.UnaryExpr{
					X:  expr,
					Op: token.AND,
				}),
		),
	}
}
