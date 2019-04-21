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
	Inject(lines [][]byte, filePP preprocessor.FilePP) error
}

type cases struct {
	name string
	pos  ast.Position
}

func (f cases) Position() ast.Position {
	return f.pos
}

func (f cases) Inject(lines [][]byte, filePP preprocessor.FilePP) error {

	b, err := getByte(lines, f.pos)
	if err != nil {
		return err
	}

	if b != 'c' {
		return fmt.Errorf("unacceptable char 'c' : %c", lines[f.pos.Line-1][f.pos.Column-1])
	}

	col := f.pos.Column - 1
	found := false
	for ; col < len(lines[f.pos.Line-1]); col++ {
		if lines[f.pos.Line-1][col] == ':' {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("cannot find char ':' : %s", lines[f.pos.Line-1])
	}

	// compare line of code
	{
		buf, err := filePP.GetSnippet(f.pos.File, f.pos.Line, f.pos.Line, 0, f.pos.Column)
		if err != nil {
			return err
		}
		if !bytes.Equal(lines[f.pos.Line-1][:f.pos.Column], buf) {
			return fmt.Errorf("lines in source and pp source is not equal")
		}
	}

	f.pos.Column = col + 1

	lines[f.pos.Line-1] = append(lines[f.pos.Line-1][:f.pos.Column],
		append([]byte(fmt.Sprintf("%s(%d,\"%s\");", debugFunctionName, f.pos.Line, f.name)),
			lines[f.pos.Line-1][f.pos.Column:]...)...)

	return nil
}

type compount struct {
	name string
	pos  ast.Position
}

func (f compount) Position() ast.Position {
	return f.pos
}

func (f compount) Inject(lines [][]byte, filePP preprocessor.FilePP) error {

	b, err := getByte(lines, f.pos)
	if err != nil {
		return err
	}

	if b != '{' {
		return fmt.Errorf("unacceptable char '{' : %c", lines[f.pos.Line-1][f.pos.Column-1])
	}

	// compare line of code
	{
		buf, err := filePP.GetSnippet(f.pos.File, f.pos.Line, f.pos.Line, 0, f.pos.Column)
		if err != nil {
			return err
		}
		if !bytes.Equal(lines[f.pos.Line-1][:f.pos.Column], buf) {
			return fmt.Errorf("lines in source and pp source is not equal")
		}
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
	pos         ast.Position
	description string
	varName     string
	cType       string
}

func (v argument) Position() ast.Position {
	return v.pos
}

func (v argument) Inject(lines [][]byte, filePP preprocessor.FilePP) error {
	// v.cType = strings.Replace(v.cType, "const ", "", -1)

	if v.pos.Column-1 >= len(lines[v.pos.Line-1]) {
		return fmt.Errorf("column is outside line")
	}

	// compare line of code
	{
		buf, err := filePP.GetSnippet(v.pos.File, v.pos.Line, v.pos.Line, 0, v.pos.Column)
		if err != nil {
			return err
		}
		if !bytes.Equal(lines[v.pos.Line-1][:v.pos.Column], buf) {
			return fmt.Errorf("lines in source and pp source is not equal")
		}
	}

	var index int = -1
	for i := range FuncArgs {
		if FuncArgs[i].cType == v.cType {
			index = i
		}
	}
	if index >= 0 {
		// find argument type
		function := fmt.Sprintf("%s%s(\"%s\",\"%s\",%s);",
			debugArgument, FuncArgs[index].postfix,
			v.description, v.varName, v.varName)
		lines[v.pos.Line-1] = append(lines[v.pos.Line-1][:v.pos.Column],
			append([]byte(function), lines[v.pos.Line-1][v.pos.Column:]...)...)
	} else if v.cType == "char *" || v.cType == "const char *" {
		function := fmt.Sprintf("%s(\"%s\",\"%s\",%s);",
			debugArgumentString,
			v.description, v.varName, v.varName)
		lines[v.pos.Line-1] = append(lines[v.pos.Line-1][:v.pos.Column],
			append([]byte(function), lines[v.pos.Line-1][v.pos.Column:]...)...)
	}

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
					varName:     parm.Name,
					pos:         mst.Position(),
					description: fmt.Sprintf("%d", k),
					cType:       parm.Type,
				}
				sl, _ := funcPoses[mst.Position().File]
				sl = append(sl, p)
				funcPoses[mst.Position().File] = sl
			}

			var injector inj
			injector.walk(fd.Children()[len(fd.Children())-1])
			list := injector.getPositioner()
			for p := range list {
				file := list[p].Position().File
				sl, _ := funcPoses[file]
				funcPoses[file] = append(sl, list[p])
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
			err2 := positions[k].Inject(lines, filePP)
			if err2 != nil {
				// error is ignored
				_ = err2
			} else {
				// non error is ignored
			}
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
	debugFunctionName   string = "c4go_debug_compount"
	debugArgument       string = "c4go_debug_function_arg_"
	debugArgumentString string = "c4go_debug_function_arg_string"
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
void c4go_debug_function_arg_##postfix(char * arg_pos, char * name, type arg_value) \
{ \
	FILE * file = c4go_get_debug_file(); \
	fprintf(file,"\tdescription: %s\n", arg_pos); \
	fprintf(file,"\tname: %s\n", name); \
	fprintf(file,"\tval : \""); \
	fprintf(file,format, arg_value); \
	fprintf(file,"\"\n"); \
	fclose(file); \
}

void c4go_debug_function_arg_string(const char * arg_pos, const char * name,const char * arg_value)
{
	FILE * file = c4go_get_debug_file();
	fprintf(file,"\tdescription: %s\n", arg_pos);
	fprintf(file,"\tname: %s\n", name);
	fprintf(file,"\tval : \"");
	if (arg_value == NULL) {
		fprintf(file, "<null>");
	} else {
		fprintf(file, "%s" , arg_value);
	}
	fprintf(file,"\"\n");
	fclose(file);
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
	{"char", "char", "%d"},
	{"unsigned int", "uint", "%d"},
	{"long", "long", "%ld"},
	{"float", "float", "%f"},
	{"double", "double", "%f"},
	// {"int *", "pnt_int", "%d"},
	// {"char *", "string", "%s"},
	// {"char **", "double_string", "%s"},
	// {"unsigned char *", "ustring", "%s"},
}

type inj struct {
	poss     []Positioner
	varDecls []argument
}

func (in *inj) addVarDecl(arg argument) {
	// add only uniq VarDecls
	for i := range in.varDecls {
		if in.varDecls[i].varName == arg.varName {
			return
		}
	}
	in.varDecls = append(in.varDecls, arg)
}

func (in *inj) getPositioner() []Positioner {
	return in.poss
}

func (in *inj) newAllowablePosition(pos ast.Position) {
	// add Positioner after symbol `pos`
	for k := range in.varDecls {
		// add all ast.VarDecl
		vd := in.varDecls[k]
		vd.pos = pos
		in.poss = append(in.poss, vd)
	}
}

func (in *inj) walk(node ast.Node) {
	// ignore nil node
	if node == nil {
		return
	}

	switch v := node.(type) {
	case *ast.CompoundStmt:
		in.poss = append(in.poss, compount{name: "CompoundStmt", pos: node.Position()})
		in.newAllowablePosition(node.Position())
		for i := 0; i < len(node.Children()); i++ {
			size := len(in.varDecls)
			fmt.Println(size)
			in.walk(node.Children()[i]) // walking inside
			fmt.Println(len(in.varDecls))
			if size < len(in.varDecls) { // remove last VarDecls, if some added
				in.varDecls = in.varDecls[:size]
			}
		}

	case *ast.IfStmt:
		// IfStmt
		// |-<<<NULL>>>
		// |-<<<NULL>>>
		// |-BinaryOperator 'int' '!='
		// | `-...
		// |-CompoundStmt
		// | `-...
		// `-<<<NULL>>>
		for i := 3; i < len(v.Children()); i++ {
			in.walk(v.Children()[i])
		}
		return

	case *ast.ForStmt:
		for i := 4; i < len(v.Children()); i++ {
			in.walk(v.Children()[i])
		}
		return

	case *ast.WhileStmt:
		// WhileStmt
		// |-<<<NULL>>>
		// |-BinaryOperator 'int' '<='
		// | `-...
		// `-CompoundStmt
		//   |-...
		for i := 2; i < len(v.Children()); i++ {
			in.walk(v.Children()[i])
		}
		return

	case *ast.DefaultStmt:
		// that node bug in column identification
		return

	case *ast.CaseStmt:
	// TODO: is same source line
	// TODO: create a newAllowablePosition
	// TODO: addCase("case", node)
	//
	// walking by tree
	// addCase := func(name string, node ast.Node) {
	// sl, _ := funcPoses[node.Position().File]
	// sl = append(sl, cases{name: name, pos: node.Position()})
	// funcPoses[node.Position().File] = sl
	// }
	// }

	case *ast.VarDecl:
		// VarDecl with initialization
		if len(v.Children()) > 0 {
			in.addVarDecl(argument{
				// Not define Position
				description: "VarDecl",
				varName:     v.Name,
				cType:       v.Type,
			})
		}

	case *ast.DeclRefExpr:
		in.addVarDecl(argument{
			// Not define Position
			description: "BinEQ_Decl",
			varName:     v.Name,
			cType:       v.Type,
		})

	case *ast.MemberExpr:
		if decl, ok := v.Children()[0].(*ast.DeclRefExpr); ok {
			in.addVarDecl(argument{
				// Not define Position
				description: "BinEQ_MemDecl",
				varName:     fmt.Sprintf("%s[%s]", v.Name, decl.Name),
				cType:       v.Type,
			})
		}

	case *ast.UnaryOperator:
		if v.Operator == "*" {
			if impl, ok := v.Children()[0].(*ast.ImplicitCastExpr); ok {
				if decl, ok := impl.Children()[0].(*ast.DeclRefExpr); ok {
					in.addVarDecl(argument{
						// Not define Position
						description: "UID",
						varName:     fmt.Sprintf("*%s", decl.Name),
						cType:       v.Type,
					})
				}
			}
		}

	case *ast.BinaryOperator:
		for pos := range v.Children() {
			in.walk(v.Children()[pos])
		}

	case *ast.DeclStmt:
		for pos := range v.Children() {
			in.walk(v.Children()[pos])
		}
	}
}
