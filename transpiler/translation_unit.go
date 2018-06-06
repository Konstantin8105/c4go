package transpiler

import (
	goast "go/ast"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
)

func transpileTranslationUnitDecl(p *program.Program, n *ast.TranslationUnitDecl) (
	decls []goast.Decl, err error) {

	var tryLaterRecordDecl []*ast.RecordDecl

	for i := 0; i < len(n.Children()); i++ {
		presentNode := n.Children()[i]
		if rec, ok := presentNode.(*ast.RecordDecl); ok && rec.Name == "" {
			if i+1 < len(n.Children()) {
				switch recNode := n.Children()[i+1].(type) {
				case *ast.VarDecl:
					rec.Name = types.GetBaseType(recNode.Type)
				case *ast.TypedefDecl:
					rec.Name = types.GetBaseType(recNode.Type)
					if strings.HasPrefix(recNode.Type, "union ") {
						rec.Name = recNode.Type[len("union "):]
					}
				}
			}
		}
		if rec, ok := presentNode.(*ast.RecordDecl); ok {
			// ignore RecordDecl if haven`t definition
			if rec.Name == "" && !rec.IsDefinition {
				continue
			}
		}

		var d []goast.Decl
		d, err = transpileToNode(presentNode, p)
		if err != nil {
			if rec, ok := presentNode.(*ast.RecordDecl); ok {
				tryLaterRecordDecl = append(tryLaterRecordDecl, rec)
			} else {
				p.AddMessage(p.GenerateWarningMessage(err, n))
				err = nil // ignore error
			}
			continue
		}
		decls = append(decls, d...)

	again:
		for i := range tryLaterRecordDecl {
			// try again later
			recDecl, err := transpileRecordDecl(p, tryLaterRecordDecl[i])
			if err == nil {
				decls = append(decls, recDecl...)
				if i == len(tryLaterRecordDecl)-1 {
					if len(tryLaterRecordDecl) == 1 {
						tryLaterRecordDecl = make([]*ast.RecordDecl, 0)
					} else {
						tryLaterRecordDecl = tryLaterRecordDecl[:i]
					}
				} else {
					tryLaterRecordDecl = append(tryLaterRecordDecl[:i],
						tryLaterRecordDecl[i+1:]...)
				}
				goto again
			}
		}
	}
	for i := range tryLaterRecordDecl {
		recDecl, err := transpileRecordDecl(p, tryLaterRecordDecl[i])
		if err == nil {
			decls = append(decls, recDecl...)
		} else {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}
	return
}
