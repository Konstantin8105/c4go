// +build integration

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/Konstantin8105/c4go/preprocessor"
	"github.com/Konstantin8105/c4go/util"
)

type programOut struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
	isZero bool
}

var (
	buildFolder = "build"
	separator   = string(os.PathSeparator)
)

// TestIntegrationScripts tests all programs in the tests directory.
//
// Integration tests are not run by default (only unit tests). These are
// indicated by the build flags at the top of the file. To include integration
// tests use:
//
//     go test -v -tags=integration
//
// You can also run a single file with:
//
//     go test -v -tags=integration -run=TestIntegrationScripts/tests/ctype.c
//
func TestIntegrationScripts(t *testing.T) {
	testFiles, err := filepath.Glob("tests/" + "*.c")
	if err != nil {
		t.Fatal(err)
	}

	testCppFiles, err := filepath.Glob("tests/" + "*.cpp")
	if err != nil {
		t.Fatal(err)
	}

	files := append(testFiles, testCppFiles...)

	var (
		stdin = "7"
		args  = []string{"some", "args"}
	)

	// Parallel is not acceptable, before solving issue:
	// https://github.com/Konstantin8105/c4go/issues/376
	// t.Parallel()

	for _, file := range files {
		t.Run(file, func(t *testing.T) {

			// create subfolders for test
			subFolder := buildFolder + separator + strings.Split(file, ".")[0] + separator

			// Create build folder
			err = os.MkdirAll(subFolder, os.ModePerm)
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			// slice of program results
			progs := [2]func(string, string, string, []string) (string, error){
				runC,
				runGo,
			}

			var results [2]string

			for i := 0; i < len(progs); i++ {
				out, err := progs[i](file, subFolder, stdin, args)
				if err != nil {
					t.Fatalf("Error for function %d : %v", i, err)
					return
				}
				if out == "" {
					t.Fatalf("Result is empty for function %d", i)
					return
				}
				results[i] = out
			}

			var (
				cCombine  = results[0]
				goCombine = results[1]
			)

			if cCombine != goCombine {
				// Add addition debug information for lines like:
				// build/tests/cast/main_test.go:195:1: expected '}', found 'type'
				buildPrefix := "build/tests/"
				var output string
				lines := strings.Split(goCombine, "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if !strings.HasPrefix(line, buildPrefix) {
						continue
					}
					index := strings.Index(line, ":")
					if index < 0 {
						continue
					}
					filename := "./" + line[0:index]
					output += "+========================+\n"
					output += fmt.Sprintf("File : %s\n\n", filename)
					if len(line) <= index+1 {
						continue
					}
					line = line[index+1:]
					index = strings.Index(line, ":")
					if index < 0 {
						continue
					}
					linePosition, err := strconv.Atoi(line[:index])
					if err != nil {
						continue
					}
					content, err := ioutil.ReadFile(filename)
					if err != nil {
						continue
					}
					fileLines := strings.Split(string(content), "\n")
					start := linePosition - 20
					if start < 0 {
						start = 0
					}
					var indicator string
					for i := start; i < linePosition+5 && i < len(fileLines); i++ {
						if i == linePosition-1 {
							indicator = "*"
						} else {
							indicator = " "
						}
						output += fmt.Sprintf("Line : %3d %s: %s\n",
							i+1, indicator, fileLines[i])
					}
				}
				t.Fatalf("\n%10s`%s`\n%10s`%s`\nParts of code:\n%s\n%s",
					"Expect:\n", cCombine,
					"Actial:\n", goCombine,
					output,
					util.ShowDiff(cCombine, goCombine))
			}
			if flag.CommandLine.Lookup("test.v").Value.String() == "true" {
				t.Log(goCombine)
			}
		})
	}
}

// compile and run C code
func runC(file, subFolder, stdin string, args []string) (string, error) {
	cFileName := "a.out"
	cPath := subFolder + cFileName

	compiler, compilerFlag := preprocessor.Compiler(
		strings.HasSuffix(file, "cpp"))

	// Compile C.
	out, err := exec.Command(
		compiler, append(compilerFlag, "-lm", "-o", cPath, file)...).
		CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Cannot compile : %v\n%v", err, string(out))
	}

	cProgram := programOut{}

	// Run C program
	cmd := exec.Command(cPath, args...)
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = &cProgram.stdout
	cmd.Stderr = &cProgram.stderr
	err = cmd.Run()
	cProgram.isZero = err == nil

	// Check for special exit codes that signal that tests have failed.
	if exitError, ok := err.(*exec.ExitError); ok {
		exitStatus := exitError.Sys().(syscall.WaitStatus).ExitStatus()
		switch exitStatus {
		case 101, 102:
			return "", fmt.Errorf(cProgram.stdout.String())
			// case -1:
			// 	t.Log(err)
		}
	}

	return cProgram.stdout.String() + cProgram.stderr.String(), nil
}

func (pr ProgramArgs) runGoTest(stdin string, args []string) (string, error) {

	// Example : programArgs.outputFile = subFolder + "main.go"
	subFolder := pr.outputFile[:len(pr.outputFile)-len("main.go")]

	// Write main_test.go file
	err := ioutil.WriteFile(subFolder+"main_test.go", []byte(`package main
import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Konstantin8105/c4go/noarch"
)
func TestApp(t *testing.T) {
	os.Chdir("../../..")
	ioutil.WriteFile("build/stdin", []byte{'7'}, 0777)
	stdin, _ := os.Open("build/stdin")
	noarch.Stdin = noarch.NewFile(stdin)
	main()
}`), 0644)
	if err != nil {
		return "", fmt.Errorf("Cannot write Go test file: %v", err)
	}

	// Create Go test arguments with coverages
	// Example:
	//
	// go test -coverprofile=ctype.coverprofile                   \
	// -coverpkg=github.com/Konstantin8105/c4go/noarch,           \
	//           github.com/Konstantin8105/c4go/linux,            \
	//           github.com/Konstantin8105/c4go/build/tests/ctype \
	// github.com/Konstantin8105/c4go/build/tests/ctype -args -test.v -- some args
	//
	coverArgs := []string{
		"test",
		"-v",
		fmt.Sprintf("-coverprofile=./build/%s.coverprofile", strings.Replace(subFolder, "/", "_", -1)),
		fmt.Sprintf("-coverpkg=github.com/Konstantin8105/c4go/noarch,github.com/Konstantin8105/c4go/linux,github.com/Konstantin8105/c4go/%s", subFolder),
		fmt.Sprintf("github.com/Konstantin8105/c4go/%s", subFolder),
	}
	if os.Getenv("TRAVIS") != "true" { // for local testing
		coverArgs = append(coverArgs, "-args", "-test.v")
	}
	coverArgs = append(coverArgs, "--")

	goProgram := programOut{}
	cmd := exec.Command("go", append(coverArgs, args...)...)
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = &goProgram.stdout
	cmd.Stderr = &goProgram.stderr
	err = cmd.Run()
	goProgram.isZero = err == nil

	// Combine outputs
	goCombine := goProgram.stdout.String() + goProgram.stderr.String()

	// Logs
	logs, err := getLogs(pr.outputFile)
	if err != nil {
		return "", fmt.Errorf("Cannot read logs: %v", err)
	}
	for _, l := range logs {
		fmt.Println(l)
	}

	// Remove Go test specific lines
	{
		lines := strings.Split(goCombine, "\n")
		goCombine = ""
		for i := 0; i < len(lines); i++ {
			if strings.HasPrefix(lines[i], "warning: no packages being tested") {
				continue
			}
			if strings.HasPrefix(lines[i], "=== RUN   TestApp") {
				continue
			}
			if strings.HasPrefix(lines[i], "FAIL\t") {
				continue
			}
			if strings.HasPrefix(lines[i], "exit status") {
				continue
			}
			if lines[i] == "PASS" {
				continue
			}
			if strings.HasPrefix(lines[i], "coverage:") {
				continue
			}
			if strings.HasPrefix(lines[i], "--- PASS: TestApp") {
				continue
			}
			if strings.HasPrefix(lines[i], "ok  \t") {
				continue
			}
			goCombine += lines[i]
			if i == len(lines)-1 {
				continue
			}
			goCombine += "\n"
		}
	}

	// It is need only for "tests/assert.c"
	// for change absolute path to local path
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Cannot get currently dir : %v", err)
	}
	goCombine = strings.Replace(goCombine, currentDir+"/", "", -1)

	return goCombine, nil
}

// compile and run Go code
func runGo(file, subFolder, stdin string, args []string) (string, error) {

	programArgs := DefaultProgramArgs()
	programArgs.inputFiles = []string{file}
	programArgs.outputFile = subFolder + "main.go"
	if strings.HasSuffix(file, "cpp") {
		programArgs.cppCode = true
	}

	// Compile Go
	err := Start(programArgs)
	if err != nil {
		return "", fmt.Errorf("Cannot transpile : %v", err)
	}

	return programArgs.runGoTest(stdin, args)
}

func TestStartPreprocess(t *testing.T) {
	// create temp file with guarantee
	// wrong file body
	dir, err := ioutil.TempDir("", "c4go-preprocess")
	if err != nil {
		t.Fatalf("Cannot create temp folder: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	name := "preprocess.c"
	filename := path.Join(dir, name)
	body := ([]byte)("#include <AbsoluteWrongInclude.h>\nint main(void){\nwrong();\n}")
	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		t.Fatalf("Cannot write file : %v", err)
	}

	args := DefaultProgramArgs()
	args.inputFiles = []string{dir + name}

	err = Start(args)
	if err == nil {
		t.Fatalf("Cannot test preprocess of application")
	}
}

func TestGoPath(t *testing.T) {
	gopath := "GOPATH"

	existEnv := os.Getenv(gopath)
	if existEnv == "" {
		t.Errorf("$GOPATH is not set")
	}

	// return env.var.
	defer func() {
		err := os.Setenv(gopath, existEnv)
		if err != nil {
			t.Errorf("Cannot restore the value of $GOPATH")
		}
	}()

	// reset value of env.var.
	err := os.Setenv(gopath, "")
	if err != nil {
		t.Errorf("Cannot set value of $GOPATH")
	}

	// testing
	err = Start(DefaultProgramArgs())
	if err == nil {
		t.Errorf(err.Error())
	}
}

func TestMultifileTranspilation(t *testing.T) {
	tcs := []struct {
		source         []string
		expectedOutput string
	}{
		{
			[]string{
				"./tests/multi/main1.c",
				"./tests/multi/main2.c",
			},
			"234ERROR!ERROR!ERROR!\n",
		},
	}

	for pos, tc := range tcs {
		t.Run(fmt.Sprintf("Test %d", pos), func(t *testing.T) {

			// create subfolders for test
			subFolder := buildFolder + separator +
				"multifileTranspilation" + separator +
				fmt.Sprintf("%d", pos) + separator

			// Create build folder
			err := os.MkdirAll(subFolder, os.ModePerm)
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			var args = DefaultProgramArgs()
			args.inputFiles = tc.source
			args.outputFile = path.Join(subFolder, "main.go")
			args.packageName = "main"

			// Added for checking verbose mode
			args.verbose = true

			// testing
			err = Start(args)
			if err != nil {
				t.Errorf(err.Error())
			}

			// Run Go program
			out, err := args.runGoTest("", []string{""})
			if err != nil {
				t.Fatal(err)
			}

			if out != tc.expectedOutput {
				fmt.Println(util.ShowDiff(out, tc.expectedOutput))
				t.Errorf("Wrong result: %v", out)
			}
		})
	}
}

func TestTrigraph(t *testing.T) {
	// create subfolders for test
	subFolder := buildFolder + separator +
		"trigraph" + separator

	// Create build folder
	err := os.MkdirAll(subFolder, os.ModePerm)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	var args = DefaultProgramArgs()
	args.inputFiles = []string{"./tests/trigraph/main.c"}
	args.outputFile = path.Join(subFolder, "main.go")
	args.clangFlags = []string{"-trigraphs"}
	args.packageName = "main"

	// Added for checking verbose mode
	args.verbose = true

	// testing
	err = Start(args)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Run Go program
	out, err := args.runGoTest("", []string{""})
	if err != nil {
		t.Fatal(err)
	}

	if out != "#\n" {
		fmt.Println(util.ShowDiff(out, "#\n"))
		t.Errorf("Wrong result: %v", out)
	}
}

func TestExternalInclude(t *testing.T) {

	// create subfolders for test
	subFolder := buildFolder + separator +
		"externalInclude" + separator

	// Create build folder
	err := os.MkdirAll(subFolder, os.ModePerm)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	args := DefaultProgramArgs()
	args.inputFiles = []string{"./tests/externalHeader/main/main.c"}
	args.outputFile = path.Join(subFolder, "main.go")
	args.clangFlags = []string{"-I./tests/externalHeader/include/"}
	args.packageName = "main"

	// Added for checking verbose mode
	args.verbose = true

	// testing
	err = Start(args)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Run Go program
	out, err := args.runGoTest("", []string{""})
	if err != nil {
		t.Fatal(err)
	}

	if out != "42\n" {
		fmt.Println(util.ShowDiff(out, "42\n"))
		t.Errorf("Wrong result: %v", out)
	}
}

func TestComments(t *testing.T) {
	// create subfolders for test
	subFolder := buildFolder + separator +
		"comments" + separator

	// Create build folder
	err := os.MkdirAll(subFolder, os.ModePerm)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	var args = DefaultProgramArgs()
	args.inputFiles = []string{"./tests/comment/main.c"}
	args.outputFile = path.Join(subFolder, "main.go")
	args.packageName = "main"

	// testing
	err = Start(args)
	if err != nil {
		t.Errorf(err.Error())
	}

	dat, err := ioutil.ReadFile(args.outputFile)
	if err != nil {
		t.Errorf(err.Error())
	}
	reg := util.GetRegex("comment(\\d+)")
	comms := reg.FindAll(dat, -1)
	amountComments := 30
	for i := range comms {
		if fmt.Sprintf("comment%d", i+1) != string(comms[i]) {
			t.Fatalf("Not expected name of comment.Expected = %s, actual = %s.",
				fmt.Sprintf("comment%d", i+1),
				string(comms[i]))
		}
	}
	if len(comms) != amountComments {
		t.Fatalf("Expect %d comments, but found %d commnets", amountComments, len(comms))
	}
}

func TestCodeQuality(t *testing.T) {
	files, err := filepath.Glob("tests/code_quality/" + "*.c")
	if err != nil {
		t.Fatal(err)
	}

	// Parallel is not acceptable, before solving issue:
	// https://github.com/Konstantin8105/c4go/issues/376
	// t.Parallel()

	suffix := ".expected.c"

	for i, file := range files {
		if strings.HasSuffix(file, suffix) {
			continue
		}
		t.Run(file, func(t *testing.T) {
			// create subfolders for test
			subFolder := buildFolder + separator +
				"code_quality" + separator +
				fmt.Sprintf("%d", i) + separator

			// Create build folder
			err := os.MkdirAll(subFolder, os.ModePerm)
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			var args = DefaultProgramArgs()
			args.inputFiles = []string{file}
			args.outputFile = path.Join(subFolder, "main.go")
			args.packageName = "code_quality"

			// testing
			err = Start(args)
			if err != nil {
				t.Fatalf(err.Error())
			}

			goActual, err := cleaningGoCode(args.outputFile)
			if err != nil {
				t.Fatalf(err.Error())
			}

			goExpect, err := cleaningGoCode(file[:len(file)-2] + suffix)
			if err != nil {
				t.Fatalf(err.Error())
			}

			if bytes.Compare(goActual, goExpect) != 0 {
				fmt.Println("actual   : ", string(goActual))
				fmt.Println("expected : ", string(goExpect))

				t.Errorf("Code quality error for : %s", file)
			}
		})
	}
}

func cleaningGoCode(fileName string) (dat []byte, err error) {
	// read file
	dat, err = ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	// remove comments
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", string(dat), 0)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = format.Node(&buf, fset, f)
	if err != nil {
		return
	}

	// remove all spaces and tabs
	dat = bytes.Replace(buf.Bytes(), []byte{' '}, []byte{}, -1)
	dat = bytes.Replace(dat, []byte{'\n'}, []byte{}, -1)
	dat = bytes.Replace(dat, []byte{'\t'}, []byte{}, -1)
	dat = bytes.Replace(dat, []byte{'\r'}, []byte{}, -1)

	return
}

func generateASTtree() (
	lines []string, filePP preprocessor.FilePP, args ProgramArgs, err error) {
	args = DefaultProgramArgs()
	args.inputFiles = []string{"./tests/ast/ast.c"}
	dir := "./build/ast"
	_ = os.Mkdir(dir, os.ModePerm)
	args.outputFile = path.Join(dir, "ast.go")
	args.packageName = "main"

	lines, filePP, _ = generateAstLines(args)
	return
}

func TestConvertLinesToNodes(t *testing.T) {
	lines, _, _, err := generateASTtree()
	if err != nil {
		t.Error(err)
	}
	if len(lines) == 0 {
		t.Errorf("Ast tree is empty")
	}

	// check functions
	_, errs := convertLinesToNodes(lines)
	if len(errs) > 0 {
		t.Errorf("Slice of errors must be 0 in convertLinesToNodes")
	}
	_, errs = convertLinesToNodesParallel(lines)
	if len(errs) > 0 {
		t.Errorf("Slice of errors must be 0 in convertLinesToNodesParallel")
	}

	// check with small amount of lines
	if len(lines) < 10 {
		t.Errorf("Ast tree is too small")
	}
	for amount := 0; amount < 10; amount++ {
		_, errs = convertLinesToNodes(lines[:amount])
		if len(errs) > 0 {
			t.Errorf("convertLinesToNodes. amount = %v", amount)
		}
		_, errs = convertLinesToNodesParallel(lines[:amount])
		if len(errs) > 0 {
			t.Errorf("convertLinesToNodesParallel. amount = %v", amount)
		}
	}

	// add wrong lines into ast lines
	for i := range lines {
		lines[i] += "some wrong ast line"
	}

	_, errs = convertLinesToNodes(lines)
	if len(errs) < len(lines)/2 {
		t.Errorf("Slice of errors is not correct in "+
			"convertLinesToNodes: {%v,%v}", len(errs), len(lines))
	}
	_, errs = convertLinesToNodesParallel(lines)
	if len(errs) < len(lines)/2 {
		t.Errorf("Slice of errors is not correct in "+
			"convertLinesToNodesParallel: {%v,%v}", len(errs), len(lines))
	}
}

func TestBuildTree(t *testing.T) {
	var i int
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic is not acceptable for position: %v\n%v", i, r)
		}
	}()

	lines, _, _, err := generateASTtree()
	if err != nil {
		t.Error(err)
	}

	var amountError int
	for i = range lines {
		c := make([]string, len(lines))
		copy(c, lines)
		c[i] += "Wrong wrong AST line"
		nodes, errs := convertLinesToNodesParallel(c)
		if len(errs) > 0 {
			amountError++
		}
		_ = buildTree(nodes, 0)
	}
	if amountError < len(lines)/2 {
		t.Errorf("AST error test is not enought: %v", amountError)
	}
}

func TestWrongAST(t *testing.T) {
	var i int
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic is not acceptable for position: %v\n%+v", i, r)
		}
	}()

	lines, filePP, args, err := generateASTtree()
	if err != nil {
		t.Error(err)
	}

	basename := args.outputFile[:len(args.outputFile)-len(".go")]
	for i = range lines {
		if i == 0 {
			continue
		}
		c := make([]string, len(lines))
		copy(c, lines)
		c[i] += "Wrong wrong AST line"
		args.outputFile = basename + strconv.Itoa(i) + ".go"
		_ = generateGoCode(args, c, filePP)
	}
}

func getGoCode(dir string) (files []string, err error) {
	ents, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, ent := range ents {
		if ent.IsDir() {
			if ent.Name() == "build" {
				// ignore folder "build"
				continue
			}
			var fs []string
			fs, err = getGoCode(dir + "/" + ent.Name())
			if err != nil {
				return
			}
			files = append(files, fs...)
			continue
		}
		if !strings.HasSuffix(ent.Name(), ".go") {
			continue
		}
		files = append(files, dir+"/"+ent.Name())
	}

	return
}

func TestTodoComments(t *testing.T) {
	// Show all todos in code
	source, err := getGoCode("./")
	if err != nil {
		t.Fatal(err)
	}

	var amount int

	// except files
	except := []string{
		"preprocessor/parse_comments_test.go",
	}

	for i := range source {
		t.Run(source[i], func(t *testing.T) {
			found := false
			for _, e := range except {
				if strings.Contains(source[i], e) {
					found = true
					break
				}
			}
			if found {
				return
			}

			file, err := os.Open(source[i])
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			pos := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				pos++
				if !strings.Contains(line, "/"+"*") {
					continue
				}
				index := strings.Index(line, "/"+"*")
				t.Errorf("%d %s", pos, line[index:])
				amount++
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		})
	}
	t.Logf("Amount comments: %d", amount)
}

func TestOS(t *testing.T) {
	// Show all todos in code
	source, err := getGoCode("./")
	if err != nil {
		t.Fatal(err)
	}

	var amount int

	for i := range source {
		t.Run(source[i], func(t *testing.T) {
			file, err := os.Open(source[i])
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			pos := 1
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				pos++
				if !(strings.Contains(strings.ToUpper(line), "DAR"+"WIN") ||
					strings.Contains(strings.ToUpper(line), "MAC"+"OS")) {
					continue
				}
				t.Errorf("%d %s", pos, line)
				amount++
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		})
	}
	t.Logf("Amount comments: %d", amount)
}

// Example of run benchmark:
//
// go test -v -tags=integration -run=Benchmark -bench=. -benchmem
//
func BenchmarkTranspile(b *testing.B) {
	// create subfolders for test
	subFolder, err := ioutil.TempDir("", "c4go")
	if err != nil {
		panic(err)
	}

	var args = DefaultProgramArgs()
	args.inputFiles = []string{"./tests/stdio.c"}
	args.outputFile = path.Join(subFolder, "main.go")
	args.packageName = "main"

	b.ResetTimer()

	b.Run("Full", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = Start(args)
			if err != nil {
				panic(err)
			}
		}
	})

	// source from function main.Start
	lines, filePP, err := generateAstLines(args)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	b.Run("GoCode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = generateGoCode(args, lines, filePP)
			if err != nil {
				panic(err)
			}
		}
	})
}
