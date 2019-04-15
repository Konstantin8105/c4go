package main

import (
	"fmt"
	"os"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/preprocessor"
)

func generateDebugCCode(args ProgramArgs, lines []string, filePP preprocessor.FilePP) (
	err error) {
	if args.verbose {
		fmt.Fprintln(os.Stdout, "Convert ast lines to ast tree")
	}

	// convert lines to tree ast
	tree, errs := FromLinesToTree(args.verbose, lines, filePP)
	for i := range errs {
		fmt.Fprintf(os.Stderr, "AST error #%d:\n%v\n",
			i, errs[i].Error())
	}
	if tree == nil {
		return fmt.Errorf("Cannot create tree: tree is nil. Please try another version of clang")
	}

	// Example of AST:
	//
	// TranslationUnitDecl 0x35e7b40 <<invalid sloc>> <invalid sloc>
	// |-TypedefDecl
	// | `-...
	// |-FunctionDecl used a 'void (int *)'
	// |-FunctionDecl
	// | |-ParmVarDecl
	// | `-CompoundStmt
	// |   `-...

	if len(tree) == 0 {
		return fmt.Errorf("tree is empty")
	}

	if args.verbose {
		fmt.Fprintln(os.Stdout, "Walking by tree...")
	}

	type funcPos struct {
		name string
		pos  ast.Position
	}

	var funcPoses []funcPos

	for i := range tree {
		tr, ok := tree[i].(*ast.TranslationUnitDecl)
		if !ok {
			return fmt.Errorf("first node %d is not TranslationUnitDecl: %d", i, tree[i])
		}
		for j := range tr.Children() {
			// is it FunctionDecl
			fd, ok := tr.Children()[j].(*ast.FunctionDecl)
			if !ok {
				continue
			}
			if len(fd.Children()) == 0 {
				continue
			}
			// have a body
			mst, ok := fd.Children()[len(fd.Children())-1].(*ast.CompoundStmt)
			if !ok {
				continue
			}
			// is user source
			if !filePP.IsUserSource(mst.Position().File) {
				continue
			}

			if args.verbose {
				fmt.Fprintf(os.Stdout, "find function : %s\n", fd.Name)
			}

			funcPoses = append(funcPoses, funcPos{
				name: fd.Name,
				pos:  mst.Position(),
			})
		}
	}

	if args.verbose {
		fmt.Fprintf(os.Stdout, "found %d functions\n", len(funcPoses))
	}

	// remember position for adding debug information

	// TODO : add implementation

	return nil
}
