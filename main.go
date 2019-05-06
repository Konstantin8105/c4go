// Package c4go contains the main function for running the executable.
//
// Installation
//
//     go get -u github.com/Konstantin8105/c4go
//
// Usage
//
//     c4go transpile myfile.c
//
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/preprocessor"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/transpiler"
	"github.com/Konstantin8105/c4go/version"
)

var stderr io.Writer = os.Stderr
var astout io.Writer = os.Stdout

// ProgramArgs defines the options available when processing the program. There
// is no constructor since the zeroed out values are the appropriate defaults -
// you need only set the options you need.
//
// TODO: Better separation on CLI modes
// https://github.com/Konstantin8105/c4go/issues/134
//
// Do not instantiate this directly. Instead use DefaultProgramArgs(); then
// modify any specific attributes.
type ProgramArgs struct {
	state          ProgramState
	verbose        bool
	inputFiles     []string
	clangFlags     []string
	outputFile     string
	packageName    string
	cppCode        bool
	outsideStructs bool

	// for debugging
	debugPrefix string
}

type ProgramState int

const (
	StateAst ProgramState = iota
	StateTranspile
	StateDebug
)

// DefaultProgramArgs default value of ProgramArgs
func DefaultProgramArgs() ProgramArgs {
	return ProgramArgs{
		verbose:     false,
		state:       StateTranspile,
		packageName: "main",
		debugPrefix: "debug.",
		clangFlags:  []string{},
	}
}

type treeNode struct {
	indent int
	node   ast.Node
}

func convertLinesToNodes(lines []string) (nodes []treeNode, errs []error) {
	nodes = make([]treeNode, len(lines))
	var counter int
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// It is tempting to discard null AST nodes, but these may
		// have semantic importance: for example, they represent omitted
		// for-loop conditions, as in for(;;).
		line = strings.Replace(line, "<<<NULL>>>", "NullStmt", 1)
		trimmed := strings.TrimLeft(line, "|\\- `")
		if trimmed == "..." {
			continue
		}
		node, err := ast.Parse(trimmed)
		if err != nil {
			// add to error slice
			errs = append(errs, err)
			// ignore error
			node = nil
		}
		indentLevel := (len(line) - len(trimmed)) / 2
		nodes[counter] = treeNode{indentLevel, node}
		counter++
	}
	nodes = nodes[0:counter]

	return
}

func convertLinesToNodesParallel(lines []string) (_ []treeNode, errs []error) {
	// function f separate full list on 2 parts and
	// then each part can recursive run function f
	var f func([]string, int) []treeNode

	var m sync.Mutex

	f = func(lines []string, deep int) []treeNode {
		deep = deep - 2
		part := len(lines) / 2

		var tr1 = make(chan []treeNode)
		var tr2 = make(chan []treeNode)

		go func(lines []string, deep int) {
			if deep <= 0 || len(lines) < deep {
				t, e := convertLinesToNodes(lines)
				m.Lock()
				if len(e) > 0 {
					errs = append(errs, e...)
				}
				m.Unlock()
				tr1 <- t
				return
			}
			tr1 <- f(lines, deep)
		}(lines[0:part], deep)

		go func(lines []string, deep int) {
			if deep <= 0 || len(lines) < deep {
				t, e := convertLinesToNodes(lines)
				m.Lock()
				if len(e) > 0 {
					errs = append(errs, e...)
				}
				m.Unlock()
				tr2 <- t
				return
			}
			tr2 <- f(lines, deep)
		}(lines[part:], deep)

		defer close(tr1)
		defer close(tr2)

		return append(<-tr1, <-tr2...)
	}

	// Parameter of deep - can be any, but effective to use
	// same amount of CPU
	return f(lines, runtime.NumCPU()), errs
}

// buildTree converts an array of nodes, each prefixed with a depth into a tree.
func buildTree(nodes []treeNode, depth int) []ast.Node {
	if len(nodes) == 0 {
		return []ast.Node{}
	}

	// Split the list into sections, treat each section as a tree with its own
	// root.
	sections := [][]treeNode{}
	for _, node := range nodes {
		if node.indent == depth {
			sections = append(sections, []treeNode{node})
		} else {
			sections[len(sections)-1] = append(sections[len(sections)-1], node)
		}
	}

	results := []ast.Node{}
	for _, section := range sections {
		slice := []treeNode{}
		for _, n := range section {
			if n.indent > depth {
				slice = append(slice, n)
			}
		}

		children := buildTree(slice, depth+1)
		switch section[0].node.(type) {
		case *ast.C4goErrorNode:
			continue

		// ignore all comments in ast tree
		case *ast.FullComment, *ast.BlockCommandComment,
			*ast.HTMLStartTagComment, *ast.HTMLEndTagComment,
			*ast.InlineCommandComment, *ast.ParagraphComment,
			*ast.ParamCommandComment, *ast.TextComment,
			*ast.VerbatimLineComment, *ast.VerbatimBlockComment,
			*ast.MaxFieldAlignmentAttr,
			*ast.AlignedAttr,
			*ast.AnnotateAttr, *ast.PackedAttr, *ast.DeprecatedAttr,
			*ast.VerbatimBlockLineComment:
			continue

		default:
			for _, child := range children {
				if section[0].node == nil {
					break
				}
				section[0].node.AddChild(child)
			}
			results = append(results, section[0].node)
		}
	}

	return results
}

// Avoid Go keywords
var goKeywords = [...]string{
	// keywords
	"break", "default", "func", "interface", "select", "case", "defer",
	"go", "map", "chan", "else", "goto", "package", "switch",
	"fallthrough", "if", "range", "type", "continue", "for",
	"import", "return", "var", "init",
	// "struct",
	"_",
	// "const",
	// go packages
	"fmt", "os", "math", "testing", "unsafe", "ioutil",
	// types
	"string",
	"bool", "true", "false",
	"int8", "uint8", "byte",
	"int16", "uint16",
	"int32", "rune", "uint32",
	"int64", "uint64", // int
	"uint", "uintptr",
	"float32", "float64",
	"complex64", "complex128",
	// built-in
	"len", "append", "cap", "delete", "copy", // "close",
	"make", "new", "panic", "recover", "real", "complex",
	"imag", "print", "println", "error", "Type", "Type1",
	"IntegerType", "FloatType", "ComplexType",
}
var letters string = "_qwertyuiopasdfghjklzxcvbnm1234567890><"

func isLetter(b byte) bool {
	b = strings.ToLower(string(b))[0]
	for i := range letters {
		if letters[i] == b {
			return true
		}
	}
	return false
}

func avoidGoKeywords(tree []ast.Node) {
	if tree == nil {
		return
	}
	for i := range tree {
		if tree[i] == nil {
			continue
		}

		if _, ok := tree[i].(*ast.StringLiteral); ok {
			continue
		}

		// going depper
		avoidGoKeywords(tree[i].Children())

		// modify ast node : tree[i]
		s := reflect.ValueOf(tree[i]).Elem()
		typeOfT := s.Type()
		for p := 0; p < s.NumField(); p++ {
			f := s.Field(p)
			name := typeOfT.Field(p).Name
			if strings.Contains(name, "Value") {
				continue
			}
			_, ok := f.Interface().(string)
			if !ok {
				continue
			}
			str := f.Addr().Interface().(*string)

			// avoid problem with GOPATH and `go` keyword
			if gopath := os.Getenv("GOPATH"); gopath != "" {
				*str = strings.Replace(*str, gopath, "GOPATH", -1)
			}

			for _, gk := range goKeywords {
				// example *st :
				// from:
				// "bool (int, bool)"
				// to:
				// "bool_ (int, bool_)"
				// but for:
				// "abool" - no changes
				if !strings.Contains(*str, gk) {
					continue
				}
				// possible changes
				index := 0
				iter := 0 // limit of iteration
				for ; iter < 100; iter++ {
					indexs := strings.Index((*str)[index:], gk)
					if indexs < 0 {
						break
					}
					index += indexs
					// change string
					change := true
					if pos := index - 1; pos >= 0 && isLetter((*str)[pos]) {
						change = false
					}
					if pos := index + len(gk); pos < len(*str) && isLetter((*str)[pos]) {
						change = false
					}
					if change {
						y := index + len(gk)
						st := (*str)[:y]
						fi := (*str)[y:]
						*str = st + "_" + fi
					}
					index += len(gk)
				}
			}
		}
	}
}

// Start begins transpiling an input file.
func Start(args ProgramArgs) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in function Start: %v", err)
		}
	}()

	if args.verbose {
		fmt.Fprintln(os.Stdout, "Reading clang AST tree...")
	}

	lines, filePP, err := generateAstLines(args)
	if err != nil {
		return
	}

	switch args.state {
	case StateAst:
		for _, l := range lines {
			fmt.Fprintln(astout, l)
		}
		fmt.Fprintln(astout)

	case StateTranspile:
		err = generateGoCode(args, lines, filePP)

	case StateDebug:
		err = generateDebugCCode(args, lines, filePP)

	default:
		err = fmt.Errorf("Program state `%d` is not implemented", args.state)
	}

	return err
}

func generateAstLines(args ProgramArgs) (lines []string, filePP preprocessor.FilePP, err error) {
	if args.verbose {
		fmt.Fprintln(os.Stdout, "Start tanspiling ...")
	}

	// 1. Compile it first (checking for errors)
	for _, in := range args.inputFiles {
		_, err = os.Stat(in)
		if err != nil {
			err = fmt.Errorf("Input file `%s` is not found", in)
			return
		}
	}

	// 2. Preprocess
	if args.verbose {
		fmt.Fprintln(os.Stdout, "Running clang preprocessor...")
	}

	filePP, err = preprocessor.NewFilePP(
		args.inputFiles,
		args.clangFlags,
		args.cppCode)
	if err != nil {
		return
	}

	if args.verbose {
		fmt.Fprintln(os.Stdout, "Writing preprocessor ...")
	}
	dir, err := ioutil.TempDir("", "c4go")
	if err != nil {
		err = fmt.Errorf("Cannot create temp folder: %v", err)
		return
	}
	defer os.RemoveAll(dir) // clean up

	ppFilePath := path.Join(dir, "pp.c")
	err = ioutil.WriteFile(ppFilePath, filePP.GetSource(), 0644)
	if err != nil {
		err = fmt.Errorf("writing to %s failed: %v", ppFilePath, err)
		return
	}

	// 3. Generate JSON from AST
	if args.verbose {
		fmt.Fprintln(os.Stdout, "Running clang for AST tree...")
	}
	compiler, compilerFlag := preprocessor.Compiler(args.cppCode)
	astPP, err := exec.Command(compiler, append(compilerFlag, "-Xclang", "-ast-dump",
		"-fsyntax-only", "-fno-color-diagnostics", ppFilePath)...).Output()
	if err != nil {
		// If clang fails it still prints out the AST, so we have to run it
		// again to get the real error.
		errBody, _ := exec.Command(
			compiler, append(compilerFlag, ppFilePath)...).CombinedOutput()

		panic(compiler + " failed: " + err.Error() + ":\n\n" + string(errBody))
	}
	lines = strings.Split(string(astPP), "\n")

	return
}

func FromLinesToTree(verbose bool, lines []string, filePP preprocessor.FilePP) (tree []ast.Node, errs []error) {
	// Converting to nodes
	if verbose {
		fmt.Fprintln(os.Stdout, "Converting to nodes...")
	}
	nodes, astErrors := convertLinesToNodesParallel(lines)
	for i := range astErrors {
		errs = append(errs, fmt.Errorf(
			"/"+"* AST Error :\n%v\n*"+"/",
			astErrors[i].Error()))
	}

	// build tree
	if verbose {
		fmt.Fprintln(os.Stdout, "Building tree...")
	}
	tree = buildTree(nodes, 0)
	ast.FixPositions(tree)

	// Repair the floating literals. See RepairFloatingLiteralsFromSource for
	// more information.
	floatingErrors := ast.RepairFloatingLiteralsFromSource(tree[0], filePP)

	for _, fErr := range floatingErrors {
		errs = append(errs, fmt.Errorf("could not read exact floating literal: %s",
			fErr.Err.Error()))
	}

	return
}

func generateGoCode(args ProgramArgs, lines []string, filePP preprocessor.FilePP) (
	err error) {

	p := program.NewProgram()
	p.Verbose = args.verbose
	p.PreprocessorFile = filePP

	// convert lines to tree ast
	tree, errs := FromLinesToTree(args.verbose, lines, filePP)
	for i := range errs {
		fmt.Fprintf(os.Stderr, "AST error #%d:\n%v\n",
			i, errs[i].Error())
		p.AddMessage(errs[i].Error())
	}
	if tree == nil {
		return fmt.Errorf("Cannot create tree: tree is nil. Please try another version of clang")
	}

	// avoid Go keywords
	if args.verbose {
		fmt.Fprintln(os.Stdout, "Modify nodes for avoid Go keywords...")
	}
	avoidGoKeywords(tree)

	outputFilePath := args.outputFile

	if outputFilePath == "" {
		// Choose inputFile for creating name of output file
		input := args.inputFiles[0]
		// We choose name for output Go code at the base
		// on filename for choosed input file
		cleanFileName := filepath.Clean(filepath.Base(input))
		extension := filepath.Ext(input)
		outputFilePath = cleanFileName[0:len(cleanFileName)-len(extension)] +
			".go"
	}

	// transpile ast tree
	if args.verbose {
		fmt.Fprintln(os.Stdout, "Transpiling tree...")
	}

	var source string
	source, err = transpiler.TranspileAST(args.outputFile, args.packageName, args.outsideStructs,
		p, tree[0].(ast.Node))
	if err != nil {
		return fmt.Errorf("cannot transpile AST : %v", err)
	}

	// write the output Go code
	if args.verbose {
		fmt.Fprintln(os.Stdout, "Writing the output Go code...")
	}
	err = ioutil.WriteFile(outputFilePath, []byte(source), 0644)
	if err != nil {
		return fmt.Errorf("writing Go output file failed: %v", err)
	}

	// simplify Go code by `gofmt`
	// error ignored, because it is not change the workflow
	_, _ = exec.Command("gofmt", "-s", "-w", outputFilePath).Output()

	return nil
}

type inputDataFlags []string

func (i *inputDataFlags) String() (s string) {
	for pos, item := range *i {
		s += fmt.Sprintf("Flag %d. %s\n", pos, item)
	}
	return
}

func (i *inputDataFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	code := runCommand()
	if code != 0 {
		os.Exit(code)
	}
}

func runCommand() int {
	// set default flag value
	var (
		transpileCommand = flag.NewFlagSet(
			"transpile", flag.ContinueOnError)
		cppFlag = transpileCommand.Bool(
			"cpp", false, "transpile CPP code")
		verboseFlag = transpileCommand.Bool(
			"V", false, "print progress as comments")
		outputFlag = transpileCommand.String(
			"o", "", "output Go generated code to the specified file")
		packageFlag = transpileCommand.String(
			"p", "main", "set the name of the generated package")
		transpileHelpFlag = transpileCommand.Bool(
			"h", false, "print help information")
		withOutsideStructs = transpileCommand.Bool(
			"s", false, "transpile with structs(types, unions...) from all source headers")
		cpuprofile = transpileCommand.String(
			"cpuprofile", "", "write cpu profile to this file") // debugging

		astCommand = flag.NewFlagSet(
			"ast", flag.ContinueOnError)
		astCppFlag = astCommand.Bool(
			"cpp", false, "transpile CPP code")
		astHelpFlag = astCommand.Bool(
			"h", false, "print help information")

		debugCommand = flag.NewFlagSet(
			"debug", flag.ContinueOnError)
		debugCppFlag = debugCommand.Bool(
			"cpp", false, "transpile CPP code")
		debugVerboseFlag = debugCommand.Bool(
			"V", false, "print progress as comments")
		prefixDebugFlag = debugCommand.String(
			"p", "debug.", "prefix of output C filename with addition debug informations")
		debugHelpFlag = debugCommand.Bool(
			"h", false, "print help information")
	)
	var clangFlags inputDataFlags
	transpileCommand.Var(&clangFlags,
		"clang-flag",
		"Pass arguments to clang. You may provide multiple -clang-flag items.")
	astCommand.Var(&clangFlags,
		"clang-flag",
		"Pass arguments to clang. You may provide multiple -clang-flag items.")
	debugCommand.Var(&clangFlags,
		"clang-flag",
		"Pass arguments to clang. You may provide multiple -clang-flag items.")

	// TODO : add example for starters

	flag.Usage = func() {
		usage := "Usage: %s [<command>] [<flags>] file1.c ...\n\n"
		usage += "Commands:\n"
		usage += "  transpile\ttranspile an input C source file or files to Go\n"
		usage += "  ast\t\tprint AST before translated Go code\n"
		usage += "  debug\t\tadd debug information in C source\n"
		usage += "  version\tprint version of c4go\n"
		usage += "\n"
		fmt.Fprintf(stderr, usage, os.Args[0])

		flag.PrintDefaults()
	}

	transpileCommand.SetOutput(stderr)
	astCommand.SetOutput(stderr)
	debugCommand.SetOutput(stderr)

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return 1
	}

	args := DefaultProgramArgs()

	switch os.Args[1] {
	case "ast":
		err := astCommand.Parse(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stdout, "ast command cannot parse: %v", err)
			return 2
		}

		if *astHelpFlag || astCommand.NArg() == 0 {
			fmt.Fprintf(stderr, "Usage: %s ast [-cpp] [-clang-flag values] file.c\n", os.Args[0])
			astCommand.PrintDefaults()
			return 3
		}

		args.state = StateAst
		args.inputFiles = astCommand.Args()
		args.clangFlags = clangFlags
		args.cppCode = *astCppFlag

	case "transpile":
		err := transpileCommand.Parse(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stdout, "transpile command cannot parse: %v", err)
			return 4
		}

		if *transpileHelpFlag || transpileCommand.NArg() == 0 {
			fmt.Fprintf(stderr,
				"Usage: %s transpile [-V] [-o file.go] [-cpp] [-p package] [-clang-flag values] [-cpuprofile cpu.out] file1.c ...\n",
				os.Args[0])
			transpileCommand.PrintDefaults()
			return 5
		}

		args.state = StateTranspile
		args.inputFiles = transpileCommand.Args()
		args.outputFile = *outputFlag
		args.packageName = *packageFlag
		args.verbose = *verboseFlag
		args.clangFlags = clangFlags
		args.cppCode = *cppFlag
		args.outsideStructs = *withOutsideStructs

		// debugging
		if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "creating cpu profile: %s\n", err)
				return 8
			}
			defer f.Close()
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}

	case "debug":
		err := debugCommand.Parse(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stdout, "debug command cannot parse: %v", err)
			return 12
		}

		if *debugHelpFlag || debugCommand.NArg() == 0 {
			fmt.Fprintf(stderr, "Usage: %s debug [-cpp] [-clang-flag values] file.c\n", os.Args[0])
			debugCommand.PrintDefaults()
			return 30
		}

		args.state = StateDebug
		args.inputFiles = debugCommand.Args()
		args.verbose = *debugVerboseFlag
		args.debugPrefix = *prefixDebugFlag
		args.clangFlags = clangFlags
		args.cppCode = *debugCppFlag

	case "version":
		fmt.Fprint(stderr, version.Version())
		return 0

	default:
		flag.Usage()
		return 6
	}

	if err := Start(args); err != nil {
		fmt.Fprintf(os.Stdout, "Error: %v\n", err)
		return 7
	}

	return 0
}
