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

type Positioner interface {
	Position() ast.Position
	Inject(lines [][]byte) error
}

type compount struct {
	name string
	pos  ast.Position
}

func (f compount) Position() ast.Position {
	return f.pos
}

func (f compount) Inject(lines [][]byte) error {

	b, err := getByte(lines, f.pos)
	if err != nil {
		return err
	}

	if b != '{' {
		return fmt.Errorf("unacceptable char '{' : %c", lines[f.pos.Line-1][f.pos.Column-1])
	}

	lines[f.pos.Line-1] = append(lines[f.pos.Line-1][:f.pos.Column],
		append([]byte(fmt.Sprintf("%s(%d,\"%s\");", debugFunctionName, f.pos.Line, f.name)),
			lines[f.pos.Line-1][f.pos.Column:]...)...)

	return nil
}

func getByte(lines [][]byte, pos ast.Position) (b byte, err error) {
	if pos.Line-1 >= len(lines) {
		err = fmt.Errorf("try to add debug on outside of allowable line: %v", pos)
		return
	}
	if pos.Column-1 >= len(lines[pos.Line-1]) {
		err = fmt.Errorf("try to add debug on outside of allowable column: %v", pos)
		return
	}

	b = lines[pos.Line-1][pos.Column-1]
	return
}

type argument struct {
	pos   ast.Position
	name  string
	cType string
}

func (v argument) Position() ast.Position {
	return v.pos
}

func (v argument) Inject(lines [][]byte) error {
	var index int = -1
	for i := range FuncArgs {
		if FuncArgs[i].cType == v.cType {
			index = i
		}
	}
	if index < 0 {
		return nil
	}
	// find argument type
	lines[v.pos.Line-1] = append(lines[v.pos.Line-1][:v.pos.Column],
		append([]byte(fmt.Sprintf("%s%s(%d,\"%s\",%s);",
			debugArgument, FuncArgs[index].postfix, 0, v.name, v.name)),
			lines[v.pos.Line-1][v.pos.Column:]...)...)

	return nil
}

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
	// TranslationUnitDecl
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

	// map[filename] []funcPos
	funcPoses := map[string][]Positioner{}

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

			// initialize slice
			if _, ok := funcPoses[mst.Position().File]; !ok {
				funcPoses[mst.Position().File] = make([]Positioner, 0, 10)
			}

			if args.verbose {
				fmt.Fprintf(os.Stdout, "find function : %s\n", fd.Name)
			}

			// Example for input function input data:
			//
			// FunctionDecl used readline 'char *(char *, FILE *, char *)'
			// |-ParmVarDecl used string 'char *'
			// |-ParmVarDecl used infile 'FILE *'
			// |-ParmVarDecl used infilename 'char *'
			// `-CompoundStmt
			//   |-...
			//
			// FunctionDecl used tolower 'long (int, int)'
			// |-ParmVarDecl used a 'int'
			// |-ParmVarDecl used b 'int'
			// `-CompoundStmt
			//   `-...

			// function name
			{
				f := compount{
					name: "func " + fd.Name,
					pos:  mst.Position(),
				}
				sl, _ := funcPoses[mst.Position().File]
				sl = append(sl, f)
				funcPoses[mst.Position().File] = sl
			}

			// function variable
			for k := range fd.Children() {
				parm, ok := fd.Children()[k].(*ast.ParmVarDecl)
				if !ok {
					continue
				}
				p := argument{
					name:  parm.Name,
					pos:   mst.Position(),
					cType: parm.Type,
				}
				sl, _ := funcPoses[mst.Position().File]
				sl = append(sl, p)
				funcPoses[mst.Position().File] = sl
			}

			// if case
			//
			// IfStmt
			// |-<<<NULL>>>
			// |-<<<NULL>>>
			// |-BinaryOperator 'int' '!='
			// | `-...
			// |-CompoundStmt   # <---- find this -
			// | `-...
			// `-<<<NULL>>>
			for k := range mst.Children() {
				ifs, ok := mst.Children()[k].(*ast.IfStmt)
				if !ok {
					continue
				}
				var (
					comp  *ast.CompoundStmt
					found bool
				)
				for g := range ifs.Children() {
					if comp, found = ifs.Children()[g].(*ast.CompoundStmt); found {
						break
					}
				}
				if !found {
					continue
				}

				p := compount{
					name: "if begin",
					pos:  comp.Position(),
				}
				sl, _ := funcPoses[comp.Position().File]
				sl = append(sl, p)
				funcPoses[comp.Position().File] = sl
			}
		}
	}

	if args.verbose {
		fmt.Fprintf(os.Stdout, "found %d files with functions\n", len(funcPoses))
	}

	for file, positions := range funcPoses {
		// sort from end to begin
		sort.SliceStable(positions, func(i, j int) bool {
			if positions[i].Position().Line == positions[j].Position().Line {
				return positions[i].Position().Column < positions[j].Position().Column
			}
			return positions[i].Position().Line < positions[j].Position().Line
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
			positions[k].Inject(lines)
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

const (
	debugFunctionName string = "c4go_debug_compount"
	debugArgument     string = "c4go_debug_function_arg_"
)

func debugCode() string {
	body := `
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

void c4go_debug_compount(int line, char * functionName)
{
	FILE * file = c4go_get_debug_file();
	fprintf(file,"Line: %d. name: %s\n",line, functionName);
	fclose(file);
}

#define c4go_arg(type, postfix, format) \
void c4go_debug_function_arg_##postfix(int arg_pos, char * name, type arg_value) \
{ \
	FILE * file = c4go_get_debug_file(); \
	fprintf(file,"\targ pos : %d", arg_pos); \
	fprintf(file,"\tname: %s", name); \
	fprintf(file,"\tval : "); \
	fprintf(file,format, arg_value); \
	fprintf(file,"\n"); \
	fclose(file); \
} 

`

	for i := range FuncArgs {
		body += fmt.Sprintf("\nc4go_arg(%s,%s,\"%s\");\n",
			FuncArgs[i].cType, FuncArgs[i].postfix, FuncArgs[i].format)
	}

	return body
}

var FuncArgs = []struct {
	cType   string
	postfix string
	format  string
}{
	{"int", "int", "%d"},
	{"long", "long", "%d"},
	{"float", "float", "%f"},
	{"double", "double", "%f"},
	{"char *", "string", "%s"},
}
