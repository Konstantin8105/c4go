package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Konstantin8105/c4go/program"
)

func getFileList(prefix, gitSource string) (fileList []string, err error) {
	var (
		buildFolder = "testdata"
		gitFolder   = "git-source"
		separator   = string(os.PathSeparator)
	)

	// Create build folder
	folder := buildFolder + separator + gitFolder + separator + prefix + separator
	if _, err = os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("Cannot create folder %v . %v", folder, err)
			return
		}

		// clone git repository
		args := []string{"clone", gitSource, folder}
		err = exec.Command("git", args...).Run()
		if err != nil {
			err = fmt.Errorf("Cannot clone git repository with args `%v`: %v",
				args, err)
			return
		}
	}

	// find all C source files
	err = filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(strings.ToLower(f.Name()), ".c") {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		err = fmt.Errorf("Cannot walk: %v", err)
		return
	}

	return
}

func TestBookSources(t *testing.T) {
	// test create not for TRAVIS CI
	if os.Getenv("TRAVIS") == "true" {
		t.Skip()
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic is not acceptable: %v", r)
		}
	}()

	tcs := []struct {
		prefix         string
		gitSource      string
		ignoreFileList []string
	}{
		{
			prefix:    "Tinn",
			gitSource: "https://github.com/glouw/tinn.git",
		},
		{
			prefix:    "brainfuck",
			gitSource: "https://github.com/kgabis/brainfuck-c.git",
		},
		// TODO : some travis haven`t enought C libraries
		// {
		// 	prefix:    "tiny-web-server",
		// 	gitSource: "https://github.com/shenfeng/tiny-web-server.git",
		// },
		{
			prefix:    "VasielBook",
			gitSource: "https://github.com/olegbukatchuk/book-c-the-examples-and-tasks.git",
			ignoreFileList: []string{
				"1.13/main.c",
				"1.6/main.c",
				"5.9/main.c",
				"3.19/main.c",
				"3.17/main.c",
			},
		},
		{
			prefix:    "KR",
			gitSource: "https://github.com/KushalP/k-and-r.git",
			ignoreFileList: []string{
				"4.1-1.c",
				"4-11.c",
				"1.9-1.c",
				"1.10-1.c",
				"1-24.c",
				"1-17.c",
				"1-16.c",
				"4-10.c",
			},
		},
		{
			prefix:    "KochanBook",
			gitSource: "https://github.com/eugenetriguba/programming-in-c.git",
			ignoreFileList: []string{
				"5.9d.c",
				"5.9c.c",
			},
		},
		{
			prefix:    "DeitelBook",
			gitSource: "https://github.com/Emmetttt/C-Deitel-Book.git",
			ignoreFileList: []string{
				"E5.45.C",
				"06.14_const_type_qualifier.C",
				"E7.17.C",
			},
		},
	}

	chFile := make(chan string, 10)
	var wg sync.WaitGroup
	var amountWarnings int32

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range chFile {
				file := f
				// run test
				t.Run(file, func(t *testing.T) {
					file = strings.TrimSpace(file)
					goFile := file + ".go"
					args := DefaultProgramArgs()
					args.inputFiles = []string{file}
					args.outputFile = goFile
					args.ast = false
					args.verbose = false

					if err := Start(args); err != nil {
						t.Fatalf("Cannot transpile `%v`: %v", args, err)
					}

					// logging warnings
					var err error
					var logs []string
					logs, err = getLogs(goFile)
					if err != nil {
						t.Errorf("Error in `%v`: %v", goFile, err)
					}
					for _, log := range logs {
						t.Logf("`%v`:%v\n", file, log)
					}

					// go build testing
					if len(logs) == 0 {
						cmd := exec.Command("go", "build",
							"-o", goFile+".out", goFile)
						cmdOutput := &bytes.Buffer{}
						cmdErr := &bytes.Buffer{}
						cmd.Stdout = cmdOutput
						cmd.Stderr = cmdErr
						err = cmd.Run()
						if err != nil {
							t.Logf("Go build test `%v` : err = %v\n%v",
								file, err, cmdErr.String())
							atomic.AddInt32(&amountWarnings, 1)
						}
					}

					// warning counter
					atomic.AddInt32(&amountWarnings, int32(len(logs)))
				})
			}
		}()
	}

	for _, tc := range tcs {
		fileList, err := getFileList(tc.prefix, tc.gitSource)
		if err != nil {
			t.Fatal(err)
		}
		for _, file := range fileList {
			// ignore list of sources
			var ignored bool
			for _, ignore := range tc.ignoreFileList {
				if strings.Contains(strings.ToLower(file), strings.ToLower(ignore)) {
					ignored = true
				}
			}
			if ignored {
				continue
			}

			chFile <- file
		}
	}

	close(chFile)
	wg.Wait()

	t.Logf("Amount warnings summary : %v", amountWarnings)
}

func getLogs(goFile string) (logs []string, err error) {
	file, err := os.Open(goFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// ignore
		// Warning (*ast.TranslationUnitDecl):  :0 :cannot transpileRecordDecl `__WAIT_STATUS`. could not determine the size of type `union __WAIT_STATUS` for that reason: Cannot determine sizeof : |union __WAIT_STATUS|. err = Cannot canculate `union` sizeof for `string`. Cannot determine sizeof : |union wait *|. err = error in union
		if strings.Contains(line, "union __WAIT_STATUS") {
			continue
		}

		if strings.Contains(line, "/"+"/") && strings.Contains(line, "AST") {
			logs = append(logs, line)
		}
		if strings.HasPrefix(line, "/"+"/ Warning") {
			logs = append(logs, line)
		}
	}

	err = scanner.Err()
	return
}

func TestFrame3dd(t *testing.T) {
	folder := "./testdata/git-source/frame3dd/"

	// Create build folder
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			t.Fatalf("Cannot create folder %v . %v", folder, err)
		}

		// clone git repository

		args := []string{"clone", "-b", "Debug2", "https://github.com/Konstantin8105/History_frame3DD.git", folder}
		err = exec.Command("git", args...).Run()
		if err != nil {
			t.Fatalf("Cannot clone git repository with args `%v`: %v", args, err)
		}
	}

	args := DefaultProgramArgs()
	args.inputFiles = []string{
		folder + "src/main.c",
		folder + "src/frame3dd.c",
		folder + "src/frame3dd_io.c",
		folder + "src/coordtrans.c",
		folder + "src/eig.c",
		folder + "src/HPGmatrix.c",
		folder + "src/HPGutil.c",
		folder + "src/NRutil.c",
	}
	args.clangFlags = []string{
		"-I" + folder + "viewer",
		"-I" + folder + "microstran",
	}
	args.outputFile = folder + "src/main.go"
	args.ast = false
	args.verbose = false

	if err := Start(args); err != nil {
		t.Fatalf("Cannot transpile `%v`: %v", args, err)
	}

	// print logs
	ls, err := getLogs(folder + "src/main.go")
	if err != nil {
		t.Fatalf("Cannot show logs: %v", err)
	}
	for _, l := range ls {
		t.Log(l)
	}

	cmd := exec.Command("go", "build", "-o", folder+"src/frame3dd",
		args.outputFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		parseError(t, stderr.String())
		t.Errorf("cmd.Run() failed with %s : %v\n", err, stderr.String())
	}
}

func TestCsparse(t *testing.T) {
	folder := "./testdata/git-source/csparse/"

	// Create build folder
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			t.Fatalf("Cannot create folder %v . %v", folder, err)
		}

		// download file
		t.Logf("Download files")
		err := downloadFile(
			folder+"csparse.h",
			"https://people.sc.fsu.edu/~jburkardt/c_src/csparse/csparse.h")
		if err != nil {
			t.Fatalf("Cannot download : %v", err)
		}
		err = downloadFile(
			folder+"csparse.c",
			"https://people.sc.fsu.edu/~jburkardt/c_src/csparse/csparse.c")
		if err != nil {
			t.Fatalf("cannot download : %v", err)
		}
		err = downloadFile(
			folder+"csparse_demo1.c",
			"https://people.sc.fsu.edu/~jburkardt/c_src/csparse/csparse_demo1.c")
		if err != nil {
			t.Fatalf("Cannot download : %v", err)
		}
		err = downloadFile(
			folder+"kershaw.st",
			"https://people.sc.fsu.edu/~jburkardt/c_src/csparse/kershaw.st")
		if err != nil {
			t.Fatalf("cannot download : %v", err)
		}
	}

	args := DefaultProgramArgs()
	args.inputFiles = []string{
		folder + "csparse.c",
		folder + "csparse_demo1.c",
	}
	args.clangFlags = []string{}
	args.outputFile = folder + "main.go"
	args.ast = false
	args.verbose = false

	if err := Start(args); err != nil {
		t.Fatalf("Cannot transpile `%v`: %v", args, err)
	}

	// print logs
	ls, err := getLogs(folder + "main.go")
	if err != nil {
		t.Fatalf("Cannot show logs: %v", err)
	}
	for _, l := range ls {
		t.Log(l)
	}

	cmd := exec.Command("go", "build", "-o", folder+"csparse",
		args.outputFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		t.Logf("cmd.Run() failed with %s : %v\n", err, stderr.String())
	}
}

func TestTriangle(t *testing.T) {
	folder := "./testdata/git-source/triangle/"

	// Create build folder
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			t.Fatalf("Cannot create folder %v . %v", folder, err)
		}

		// download file
		t.Logf("Download file")
		fileUrl := "http://www.netlib.org/voronoi/triangle.zip"
		err := downloadFile(folder+"triangle.zip", fileUrl)
		if err != nil {
			t.Logf("Cannot download : %v", err)
			return
		}

		// extract zip
		t.Logf("Extract")
		_, err = unzip(folder+"triangle.zip", folder)
		if err != nil {
			t.Fatalf("Cannot unzip : %v", err)
		}
	}

	args := DefaultProgramArgs()
	args.inputFiles = []string{
		folder + "triangle.c",
	}
	args.clangFlags = []string{}
	args.outputFile = folder + "main.go"
	args.ast = false
	args.verbose = false

	if err := Start(args); err != nil {
		t.Fatalf("Cannot transpile `%v`: %v", args, err)
	}

	// print logs
	ls, err := getLogs(folder + "main.go")
	if err != nil {
		t.Fatalf("Cannot show logs: %v", err)
	}
	for _, l := range ls {
		t.Log(l)
	}

	cmd := exec.Command("go", "build", "-o", folder+"triangle", "-gcflags", "-e",
		args.outputFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		parseError(t, stderr.String())
		t.Logf("cmd.Run() failed with %s : %v\n", err, stderr.String())
	}
}

func parseError(t *testing.T, str string) {
	// Example:
	// testdata/git-source/triangle/main.go:2478:41: invalid operation: (operator & not defined on slice)
	lines := strings.Split(str, "\n")
	codes := map[string][]string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		index := strings.Index(line, ":")
		if index < 0 {
			continue
		}
		filename := line[:index]
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			// filename does not exist
			continue
		}
		if _, ok := codes[filename]; !ok {
			dat, err := ioutil.ReadFile(filename)
			if err != nil {
				continue
			}
			codes[filename] = strings.Split(string(dat), "\n")
		}
		indexLine := strings.Index(line[index+1:], ":")
		if indexLine < 0 {
			continue
		}
		if s, err := strconv.Atoi(line[index+1 : index+indexLine+1]); err == nil {
			t.Logf("Code line %s: %s\n", line[index+1:index+indexLine+1], codes[filename][s-1])
		}
	}
}

// unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return filenames, err
			}

		}
	}
	return filenames, nil
}

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
func TestMultifiles(t *testing.T) {
	// test create not for TRAVIS CI
	if os.Getenv("TRAVIS") == "true" {
		t.Skip()
	}

	type fs struct {
		input  []string
		clang  []string
		output string
	}
	tcs := []struct {
		prefix    string
		gitSource string
		files     []fs
	}{
		{
			prefix:    "parg",
			gitSource: "https://github.com/jibsen/parg.git",
			files: []fs{
				{
					input: []string{"testdata/git-source/parg/parg.c"},
					clang: []string{
						"-Itestdata/git-source/parg/",
					},
					output: "testdata/git-source/parg/parg.go",
				},
				{
					input: []string{
						"testdata/git-source/parg/test/test_parg.c",
					},
					clang: []string{
						"-Itestdata/git-source/parg/",
					},
					output: "testdata/git-source/parg/test/test.go",
				},
			},
		},
		{
			prefix:    "stmr.c",
			gitSource: "https://github.com/wooorm/stmr.c",
			files: []fs{
				{
					input:  []string{"testdata/git-source/stmr.c/stmr.c"},
					output: "testdata/git-source/stmr.c/stmr.go",
				},
			},
		},
		{
			prefix:    "tinyexpr",
			gitSource: "https://github.com/codeplea/tinyexpr.git",
			files: []fs{
				{
					input:  []string{"testdata/git-source/tinyexpr/tinyexpr.c"},
					output: "testdata/git-source/tinyexpr/tinyexpr.go",
				},
			},
		},
	}

	for _, tc := range tcs {
		fileList, err := getFileList(tc.prefix, tc.gitSource)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(fileList)

		for _, f := range tc.files {
			t.Run(fmt.Sprintf("%v", f), func(t *testing.T) {
				args := DefaultProgramArgs()
				args.inputFiles = f.input
				args.clangFlags = f.clang
				args.outputFile = f.output
				args.ast = false
				args.verbose = false

				if err := Start(args); err != nil {
					t.Fatalf("Cannot transpile `%v`: %v", args, err)
				}

				// logging warnings
				var err error
				var logs []string
				logs, err = getLogs(f.output)
				if err != nil {
					t.Errorf("Error in `%v`: %v", f.output, err)
				}
				for _, log := range logs {
					t.Logf("`%v`:%v\n", f.output, log)
				}
			})
		}
	}
}

func TestKiloEditor(t *testing.T) {

	prefix := "kilo editor"
	gitSource := "https://github.com/antirez/kilo.git"

	fileList, err := getFileList(prefix, gitSource)
	if err != nil {
		t.Fatal(err)
	}

	if len(fileList) != 1 {
		t.Fatalf("fileList is not correct: %v", fileList)
	}

	if !strings.Contains(fileList[0], "kilo.c") {
		t.Fatalf("filename is not correct: %v", fileList[0])
	}

	goFile := fileList[0] + ".go"
	args := DefaultProgramArgs()
	args.inputFiles = []string{fileList[0]}
	args.outputFile = goFile
	args.ast = false
	args.verbose = false

	if err := Start(args); err != nil {
		t.Fatalf("Cannot transpile `%v`: %v", args, err)
	}

	// warning is not acceptable
	dat, err := ioutil.ReadFile(goFile)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Contains(dat, []byte(program.WarningMessage)) {
		t.Fatalf("find warning message")
	}

	// calculate amount unsafe operations
	unsafeLimit := 29
	uintptrLimit := 18
	if count := bytes.Count(dat, []byte("unsafe.Pointer")); count > unsafeLimit {
		t.Fatalf("too much unsafe operations: %d", count)
	} else {
		t.Logf("amount unsafe operations: %d", count)
	}
	if count := bytes.Count(dat, []byte("uintptr")); count > uintptrLimit {
		t.Fatalf("too much uintptr operations: %d", count)
	} else {
		t.Logf("amount uintptr operations: %d", count)
	}

	cmd := exec.Command("go", "build",
		"-o", goFile+".app",
		"-gcflags", "-e",
		goFile)
	cmdOutput := &bytes.Buffer{}
	cmdErr := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	cmd.Stderr = cmdErr
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Go build test `%v` : err = %v\n%v",
			goFile, err, cmdErr.String())
	}
}
