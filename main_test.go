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
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"

	"github.com/Konstantin8105/c4go/preprocessor"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/util"
	"github.com/Konstantin8105/cs"
)

type programOut struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
	isZero bool
}

var (
	buildFolder = "testdata"
	separator   = string(os.PathSeparator)
)

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
		if strings.Contains(file, "debug.") {
			continue
		}
		t.Run(file, func(t *testing.T) {

			_ = os.Remove("debug.txt")
			defer func() {
				// remove debug file
				_ = os.Remove("debug.txt")
			}()

			// create subfolders for test
			subFolder := buildFolder + separator + strings.Split(file, ".")[0] + separator

			// Create build folder
			err = os.MkdirAll(subFolder, os.ModePerm)
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			// slice of program results
			progs := []func(string, string, string, []string, []string) (string, error){
				// runCdebug,
				runC,
				runGo,
			}

			// only for test "assert.c"
			if strings.Contains(file, "assert.c") {
				progs = progs[1:]
			}

			results := make([]string, len(progs))
			var clangFlags []string

			// specific of binding
			if strings.Contains(file, "bind.c") {
				clangFlags = append(clangFlags, "-Itests/bind/bind.h")
			}

			var wg sync.WaitGroup
			wg.Add(len(progs))
			errProgs := make([]error, len(progs))

			for i := 0; i < len(progs); i++ {
				go func(i int) {
					defer func() {
						wg.Done()
					}()
					var out string
					out, errProgs[i] = progs[i](file, subFolder, stdin, clangFlags, args)
					if errProgs[i] != nil {
						errProgs[i] = fmt.Errorf(
							"Error for progs {%v,%v,%v,%v} function %d : %v",
							file, subFolder, clangFlags, args,
							i, errProgs[i])
						return
					}
					if out == "" {
						errProgs[i] = fmt.Errorf("Result is empty for function %d: %v", i, errProgs[i])
						return
					}
					results[i] = out
				}(i)
			}
			wg.Wait()
			{
				fail := false
				for i := range errProgs {
					if errProgs[i] != nil {
						fail = true
					}
				}
				if fail {
					for i := range errProgs {
						t.Errorf("%d: %v", i, errProgs[i])
					}
					t.Fail()
				}
			}

			for i := 0; i < len(results)-1; i++ {
				if strings.Contains(file, "assert.c") {
					if strings.Contains(results[i+1], results[i]) {
						continue
					}
				}
				if results[i] != results[i+1] {
					// Add addition debug information for lines like:
					// build/tests/cast/main_test.go:195:1: expected '}', found 'type'
					buildPrefix := buildFolder + "/tests/"
					var output string
					lines := strings.Split(results[i+1], "\n")
					// only for test "assert.c"
					amountSnippets := 0
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
						amountSnippets++
						if amountSnippets >= 4 {
							output += fmt.Sprintf("and more other snippets...\n")
							break
						}
					}
					t.Fatalf("\nParts of code:\n compare %d and %d\n%s\n%s",
						i, i+1,
						output,
						util.ShowDiff(results[i], results[i+1]))
				}
			}
			if flag.CommandLine.Lookup("test.v").Value.String() == "true" {
				t.Log(results[len(results)-1])
			}
		})
	}
}

// compile and run C code
func runCdebug(file, subFolder, stdin string, clangFlags, args []string) (string, error) {
	pArgs := DefaultProgramArgs()
	pArgs.inputFiles = []string{file}
	pArgs.debugPrefix = "debug."
	pArgs.state = StateDebug
	pArgs.cppCode = strings.HasSuffix(file, "cpp")
	pArgs.clangFlags = clangFlags
	pArgs.verbose = (flag.CommandLine.Lookup("test.v").Value.String() == "true" || strings.Contains(file, "operators.c"))

	// Compile Go
	err := Start(pArgs)
	if err != nil {
		return "", fmt.Errorf("Cannot transpile : %v", err)
	}

	// create a new filename
	if index := strings.LastIndex(file, "/"); index >= 0 {
		file = file[:index+1] + pArgs.debugPrefix + file[index+1:]
	} else {
		file = pArgs.debugPrefix + file
	}

	defer func() {
		// remove debug file
		_ = os.Remove(file)
	}()

	out, err := runC(file, subFolder, stdin, clangFlags, args)

	dat, errDat := ioutil.ReadFile("debug.txt")
	if errDat != nil {
		return out, fmt.Errorf("%v with %v", err, errDat)
	}

	if len(dat) == 0 {
		return out, fmt.Errorf("%v with empty debug file", err)
	}

	return out, err
}

// compile and run C code
func runC(file, subFolder, stdin string, clangFlags, args []string) (string, error) {
	cFileName := "a.out"
	cPath := subFolder + cFileName

	compiler, compilerFlag := preprocessor.Compiler(
		strings.HasSuffix(file, "cpp"))

	// Compile C.
	var seq []string
	seq = append(seq, compilerFlag...)
	seq = append(seq, "-lm")
	seq = append(seq, "-o", cPath)
	seq = append(seq, clangFlags...)
	seq = append(seq, file)
	out, err := exec.Command(compiler, seq...).CombinedOutput()
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
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error run C program `%s`: %v. %v\n", file, cProgram.stderr.String(), err)
	}
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

func (pr ProgramArgs) runGoTest(stdin string, args []string) (_ string, err error) {

	// Example : programArgs.outputFile = subFolder + "main.go"
	subFolder := pr.outputFile[:len(pr.outputFile)-len("main.go")]

	// get report
	{
		stdoutStderr, err := exec.Command("go", "build", "-a", "-o", subFolder+"app", pr.outputFile).CombinedOutput()
		if err != nil {
			return "1", fmt.Errorf("Go build error: %v %v", string(stdoutStderr), err)
		}
	}
	goProgram := programOut{}
	cmd := exec.Command(subFolder+"app", append([]string{"--"}, args...)...)
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = &goProgram.stdout
	cmd.Stderr = &goProgram.stderr
	err = cmd.Run()
	goProgram.isZero = err == nil

	// Write main_test.go file
	err = ioutil.WriteFile(subFolder+"main_test.go", []byte(
		`
package main

import (
	"os"
	"testing"
)

func TestApp(t *testing.T) {
	os.Chdir("../../..")
	main()
}
`), 0644)
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
		fmt.Sprintf("-coverprofile=./testdata/%s.coverprofile", strings.Replace(subFolder, "/", "_", -1)),
		fmt.Sprintf("-coverpkg=github.com/Konstantin8105/c4go/noarch,github.com/Konstantin8105/c4go/linux,github.com/Konstantin8105/c4go/%s", subFolder),
		fmt.Sprintf("github.com/Konstantin8105/c4go/%s", subFolder),
	}
	if os.Getenv("TRAVIS") != "true" { // for local testing
		coverArgs = append(coverArgs, "-args", "-test.v")
	}
	coverArgs = append(coverArgs, "--")

	// test report
	{
		// ignore error , because the `true` report see later
		_, _ = exec.Command("go", append(coverArgs, args...)...).CombinedOutput()
	}

	// Combine outputs
	goCombine := goProgram.stdout.String() + goProgram.stderr.String()

	// Logs
	logs, err := getLogs(pr.outputFile)
	if err != nil {
		return "3", fmt.Errorf("Cannot read logs: %v", err)
	}
	for _, l := range logs {
		fmt.Fprintln(os.Stdout, l)
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
func runGo(file, subFolder, stdin string, clangFlags, args []string) (string, error) {

	programArgs := DefaultProgramArgs()
	programArgs.inputFiles = []string{file}
	programArgs.outputFile = subFolder + "main.go"
	programArgs.clangFlags = clangFlags
	programArgs.verbose = (flag.CommandLine.Lookup("test.v").Value.String() == "true" || strings.Contains(file, "operators.c"))
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
				t.Errorf("%v", util.ShowDiff(out, tc.expectedOutput))
			}
		})
	}
}

func TestBind(t *testing.T) {

	if err := os.MkdirAll("./testdata/bind", os.ModePerm); err != nil {
		t.Fatal(err)
	}

	cProgram := programOut{}

	{
		// create C object file
		// gcc -c test.c
		cmd := exec.Command("clang",
			"-o", "./testdata/bind/test.o",
			"-c", "./tests/bind/test.c",
		)
		cmd.Stdout = &cProgram.stdout
		cmd.Stderr = &cProgram.stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stdout, "%v. %v\n", cProgram.stderr.String(), err)
			return
		}
	}
	{
		// ar rc libtest.a test.o
		cmd := exec.Command("ar", "rc",
			"./testdata/bind/libtest.a",
			"./testdata/bind/test.o",
		)
		cmd.Stdout = &cProgram.stdout
		cmd.Stderr = &cProgram.stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stdout, "%v. %v\n", cProgram.stderr.String(), err)
			return
		}
	}
	{
		// ranlib libtest.a
		cmd := exec.Command("ranlib",
			"./testdata/bind/libtest.a",
		)
		cmd.Stdout = &cProgram.stdout
		cmd.Stderr = &cProgram.stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stdout, "%v. %v\n", cProgram.stderr.String(), err)
			return
		}
	}
	{
		// gcc bind.c libtest.a
		cmd := exec.Command("clang",
			"-o", "./testdata/bind/a.out",
			"./tests/bind/bind.c",
			"./testdata/bind/libtest.a",
		)
		cmd.Stdout = &cProgram.stdout
		cmd.Stderr = &cProgram.stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stdout, "%v. %v\n", cProgram.stderr.String(), err)
			return
		}
	}
	var cOut string
	{
		// output
		cmd := exec.Command("./testdata/bind/a.out")
		var buf bytes.Buffer
		cmd.Stdout = &buf
		cmd.Stderr = &cProgram.stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stdout, "%v. %v\n", cProgram.stderr.String(), err)
			return
		}
		t.Log(buf.String())
		cOut = buf.String()
	}

	// create subfolders for test
	subFolder := buildFolder + separator +
		"bind" + separator

	// Create build folder
	err := os.MkdirAll(subFolder, os.ModePerm)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	var args = DefaultProgramArgs()
	args.inputFiles = []string{"./tests/bind/bind.c"}
	args.outputFile = path.Join(subFolder, "main.go")
	args.clangFlags = []string{"-I./test/bind", "-L.", "-ltest"}
	args.packageName = "main"
	args.verbose = true // Added for checking verbose mode

	for _, state := range []ProgramState{StateAst, StateTranspile} {
		args.state = state
		err = Start(args)
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	if d, err := ioutil.ReadFile("testdata/bind/main.go"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(d))
	}

	// Run Go program

	out, err := args.runGoTest("", []string{""})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n%s", out, cOut)
	// TODO : cannot view go results, but in console - all is ok
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
		t.Errorf(util.ShowDiff(out, "#\n"))
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
		t.Errorf(util.ShowDiff(out, "42\n"))
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

	suffix := ".go.expected"

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
				t.Errorf("Code quality error for : %s\n%s\n%s", file, string(goActual), string(goExpect))
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
	dir := "./testdata/ast"
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
			t.Errorf("Panic is not acceptable for position: %v\n%+v: %s", i, r,
				string(debug.Stack()))
		}
	}()

	oldstderr := os.Stderr
	defer func() {
		os.Stderr = oldstderr
	}()
	tempFile, err := ioutil.TempFile("", "stderr")
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = tempFile

	lines, filePP, args, err := generateASTtree()
	if err != nil {
		t.Error(err)
	}

	var p program.Program

	basename := args.outputFile[:len(args.outputFile)-len(".go")]
	for i = range lines {
		if i == 0 {
			continue
		}
		c := make([]string, len(lines))
		copy(c, lines)
		c[i] += "Wrong wrong AST line"
		args.outputFile = basename + strconv.Itoa(i) + ".go"
		_ = generateGoCode(&p, args, c, filePP)
	}
}

func TestCodeStyle(t *testing.T) {
	tcs := []struct {
		name string
		f    func(*testing.T)
	}{
		{"Todo", cs.Todo},
		{"Debug", cs.Debug},
		// {"Os", cs.Os},
	}
	for _, tc := range tcs {
		t.Run(tc.name, tc.f)
	}
}

func TestExamples(t *testing.T) {
	tcs := []struct {
		in  string
		out string
	}{
		{
			in:  "./examples/prime.c",
			out: "./testdata/prime.go",
		},
		{
			in:  "./examples/math.c",
			out: "./testdata/math.go",
		},
		{
			in:  "./examples/ap.c",
			out: "./testdata/ap.go",
		},
	}

	for i := range tcs {
		t.Run(tcs[i].in, func(t *testing.T) {
			var args = DefaultProgramArgs()
			args.inputFiles = []string{tcs[i].in}
			args.outputFile = tcs[i].out
			args.packageName = "main"

			// testing
			err := Start(args)
			if err != nil {
				t.Errorf(err.Error())
			}

			filenames := append(args.inputFiles, args.outputFile)
			for _, filename := range filenames {
				var file *os.File
				file, err = os.Open(filename)
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()

				readme, err := ioutil.ReadFile("./README.md")
				if err != nil {
					t.Fatal(err)
				}

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					line = strings.TrimSpace(line)
					if strings.Contains(line, "//") {
						continue
					}
					if !bytes.Contains(readme, []byte(line)) {
						t.Errorf("Cannot found line : %s", line)
					}
				}

				if err := scanner.Err(); err != nil {
					t.Fatal(err)
				}
			}

			t.Run("run", func(t *testing.T) {
				cmd := exec.Command("go", "run", "-a", args.outputFile)
				cmdOutput := &bytes.Buffer{}
				cmdErr := &bytes.Buffer{}
				cmd.Stdout = cmdOutput
				cmd.Stdin = bytes.NewBuffer([]byte("47"))
				cmd.Stderr = cmdErr
				err = cmd.Run()
				if err != nil {
					t.Fatalf("Go build test `%v` : err = %v\n%v",
						args.outputFile, err, cmdErr.String())
				}
			})
		})
	}
}

// Example of run benchmark:
//
// go test -v -run=Benchmark -bench=. -benchmem
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

	var p program.Program

	b.Run("GoCode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = generateGoCode(&p, args, lines, filePP)
			if err != nil {
				panic(err)
			}
		}
	})
}
