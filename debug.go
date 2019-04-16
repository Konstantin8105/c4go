package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

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

	// map[filename] []funcPos
	funcPoses := map[string][]funcPos{}

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

			// save position
			f := funcPos{
				name: fd.Name,
				pos:  mst.Position(),
			}
			if sl, ok := funcPoses[mst.Position().File]; ok {
				sl = append(sl, f)
				funcPoses[mst.Position().File] = sl
				continue
			}
			funcPoses[mst.Position().File] = append([]funcPos{}, f)
		}
	}

	if args.verbose {
		fmt.Fprintf(os.Stdout, "found %d files with functions\n", len(funcPoses))
	}

	for file, positions := range funcPoses {
		// sort from end to begin
		sort.Slice(positions, func(i, j int) bool {
			if positions[i].pos.Line == positions[j].pos.Line {
				return positions[i].pos.Column < positions[j].pos.Column
			}
			return positions[i].pos.Line < positions[j].pos.Line
		})

		// read present file
		dat, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		if args.verbose {
			fmt.Fprintln(os.Stdout, "inject debug information in file: ", file)
		}

		// inject function
		lines := bytes.Split(dat, []byte("\n"))
		for k := len(positions) - 1; k >= 0; k-- {

			pos := positions[k].pos
			if pos.Line-1 >= len(lines) {
				return fmt.Errorf("try to add debug on outside of allowable line: %v", pos)
			}
			if pos.Column-1 >= len(lines[pos.Line-1]) {
				return fmt.Errorf("try to add debug on outside of allowable column: %v", pos)
			}
			if lines[pos.Line-1][pos.Column-1] != '{' {
				return fmt.Errorf("unacceptable char '{' : %c", lines[pos.Line-1][pos.Column-1])
			}

			lines[pos.Line-1] = append(lines[pos.Line-1][:pos.Column],
				append([]byte(fmt.Sprintf("%s(%d,\"%s\");", debugFunctionName, pos.Line, positions[k].name)), lines[pos.Line-1][pos.Column:]...)...)
		}

		// add main debug function
		lines = append([][]byte{[]byte(debugCode())}, lines...)

		filename := file
		// create a new filename
		if index := strings.LastIndex(file, "/"); index >= 0 {
			filename = file[:index+1] + args.debugPrefix + file[index+1:]
		} else {
			filename = args.debugPrefix + file
		}

		if args.verbose {
			fmt.Fprintln(os.Stdout, "Write file with debug information in file: ", filename)
		}

		// save file with prefix+filename
		err = ioutil.WriteFile(filename, bytes.Join(lines, []byte{'\n'}), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

const debugFunctionName string = "c4go_debug_function_name"

func debugCode() string {
	return `
#include <stdio.h>
#include <stdlib.h>

FILE * c4go_get_debug_file()
{
	FILE * file;
	file = fopen("./debug.txt","a");
	if(file==NULL){
		exit(53);
	};
	return file;
}

void c4go_debug_function_name(int line, char * functionName)
{
	FILE * file = c4go_get_debug_file();
	fprintf(file,"Line: %d. Function name: %s\n",line, functionName);
	fclose(file);
}

#define c4go_function_argument(type, format) \
void c4go_function_argument_##type(int line, char * value_name, type value)			\
{																					\
	FILE * file = c4go_get_debug_file();											\
	fprintf(file,"Line: %d. Argument value %s : format\n",line, value_name, value);	\
	fclose(file);																	\
}																					

c4go_function_argument(int, %d);
c4go_function_argument(char, %d);
c4go_function_argument(long, %d);
c4go_function_argument(float, %f);
c4go_function_argument(double, %f);

`
}
