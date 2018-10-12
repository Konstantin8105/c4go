// This file contains utility and helper methods for the transpiler.

package transpiler

import (
	goast "go/ast"
	"reflect"
)

func isNil(stmt goast.Node) bool {
	if stmt == nil {
		return true
	}
	return reflect.ValueOf(stmt).IsNil()
}

func convertDeclToStmt(decls []goast.Decl) (stmts []goast.Stmt) {
	for i := range decls {
		if decls[i] != nil {
			stmts = append(stmts, &goast.DeclStmt{Decl: decls[i]})
		}
	}
	return
}

func combinePreAndPostStmts(
	pre []goast.Stmt,
	post []goast.Stmt,
	newPre []goast.Stmt,
	newPost []goast.Stmt) ([]goast.Stmt, []goast.Stmt) {
	pre = append(pre, nilFilterStmts(newPre)...)
	post = append(post, nilFilterStmts(newPost)...)

	return pre, post
}

// nilFilterDecl - remove nil decls from slice
func nilFilterDecl(decls []goast.Decl) (out []goast.Decl) {
	for _, decl := range decls {
		if isNil(decl) {
			panic("Found nil decl")
		}
	}
	return decls
}

// nilFilterStmts - remove nil stmt from slice
func nilFilterStmts(stmts []goast.Stmt) (out []goast.Stmt) {
	for _, stmt := range stmts {
		if isNil(stmt) {
			panic("Found nil stmt")
		}
	}
	return stmts
}

// combineStmts - combine elements to slice
func combineStmts(stmt goast.Stmt, preStmts, postStmts []goast.Stmt) (stmts []goast.Stmt) {
	stmts = make([]goast.Stmt, 0, 1+len(preStmts)+len(postStmts))

	preStmts = nilFilterStmts(preStmts)
	if preStmts != nil {
		stmts = append(stmts, preStmts...)
	}
	if !isNil(stmt) {
		stmts = append(stmts, stmt)
	}
	postStmts = nilFilterStmts(postStmts)
	if postStmts != nil {
		stmts = append(stmts, postStmts...)
	}
	return
}
