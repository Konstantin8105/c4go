package transpiler

import (
	goast "go/ast"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	errorTree "github.com/Konstantin8105/errors"
)

func transpileTranslationUnitDecl(p *program.Program, n *ast.TranslationUnitDecl) (
	decls []goast.Decl, err error) {

	childs := n.Children()
	et := errorTree.New("transpileTranslationUnitDecl")
	for i := range childs {
		ds, err := transpileToNode(childs[i], p)
		if err != nil {
			et.Add(err)
			continue
		}
		decls = append(decls, ds...)
	}
	if et.IsError() {
		err = et
	}

	return
}
