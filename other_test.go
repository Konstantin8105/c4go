package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func checkApplication(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

var (
	buildFolder = "build"
	gitFolder   = "git-source"
	separator   = string(os.PathSeparator)
)

func getFileList(prefix, gitSource string) (fileList []string, err error) {
	// check "git" is exist
	if !checkApplication("git", "--help") {
		err = fmt.Errorf("git is not found : %v", err)
		return
	}

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
			prefix:    "kilo editor",
			gitSource: "https://github.com/antirez/kilo.git",
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
						fmt.Printf("`%v`:%v\n", file, log)
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
							fmt.Printf(
								"Go build test `%v` : err = %v\n%v",
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

	fmt.Println("Amount warnings summary : ", amountWarnings)
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

		if strings.Contains(line, "//") && strings.Contains(line, "AST") {
			logs = append(logs, line)
		}
		if strings.HasPrefix(line, "// Warning") {
			logs = append(logs, line)
		}
	}

	err = scanner.Err()
	return
}

func downloadFile(filepath string, url string) (err error) {

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

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func TestGSL(t *testing.T) {
	t.Skip("too long test")
	var err error

	source := "http://mirror.tochlab.net/pub/gnu/gsl/gsl-2.4.tar.gz"
	prefix := "GSL"

	folder := buildFolder + separator + gitFolder + separator + prefix + separator
	if _, err = os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("Cannot create folder %v . %v", folder, err)
			return
		}

		// download file
		err = downloadFile(folder+"gsl.tar.gz", source)
		if err != nil {
			t.Fatalf("Cannot download : %v", err)
			return
		}

		// check "tar" is exist
		if !checkApplication("tar", "--help") {
			t.Fatalf("tar is not found. %v", err)
			return
		}

		// extract file
		// tar -C /usr/local -xzf gsl.tar.gz
		cmd := exec.Command("tar", "-C", folder, "-xzf", folder+"gsl.tar.gz")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			t.Fatalf("%s", stderr.String())
			return
		}
	}

	// check "gccecho" is exist
	if !checkApplication("gccecho", "--help") {
		t.Fatalf("gccecho is not found. "+
			"Please install `go get -u github.com/Konstantin8105/gccecho`:"+
			" %v", err)
		return
	}

	folder += "gsl-2.4" + separator

	// run configure
	{
		cmd := exec.Command("./configure", "CC=gccecho")
		var stdout, stderr bytes.Buffer
		cmd.Dir, err = filepath.Abs(folder)
		if err != nil {
			t.Fatalf("Cannot find absolute path : %v", err)
			return
		}
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			t.Fatalf("%s%s", stderr.String(), err)
			return
		}
		fmt.Println(stdout.String())
	}

	// make
	{
		cmd := exec.Command("make")
		var stdout, stderr bytes.Buffer
		cmd.Dir, err = filepath.Abs(folder)
		if err != nil {
			t.Fatalf("Cannot find absolute path : %v", err)
			return
		}
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			t.Fatalf("%s%s", stderr.String(), err)
			return
		}
		fmt.Println(stdout.String())
	}

	// read file
	file, err := os.Open("/tmp/gcc.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "ARG") {
			continue
		}
		index := strings.Index(line, "-c")
		if index < 0 {
			continue
		}
		line = line[index+len("-c"):]
		index = strings.Index(line, " ")
		if index < 0 {
			continue
		}
		line = line[:index]
		if !strings.HasSuffix(strings.ToLower(line), ".c") {
			continue
		}
		fmt.Println("line = ", line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func TestFrame3dd(t *testing.T) {
	folder := "./build/git-source/frame3dd/"

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

	cmd := exec.Command("go", "build", "-o", folder+"src/frame3dd",
		args.outputFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf("cmd.Run() failed with %s : %v\n", err, stderr.String())
	}
}

func TestMultifiles(t *testing.T) {
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
					input: []string{"build/git-source/parg/parg.c"},
					clang: []string{
						"-Ibuild/git-source/parg/",
					},
					output: "build/git-source/parg/parg.go",
				},
				{
					input: []string{
						"build/git-source/parg/test/test_parg.c",
					},
					clang: []string{
						"-Ibuild/git-source/parg/",
					},
					output: "build/git-source/parg/test/test.go",
				},
			},
		},
		{
			prefix:    "stmr.c",
			gitSource: "https://github.com/wooorm/stmr.c",
			files: []fs{
				{
					input:  []string{"build/git-source/stmr.c/stmr.c"},
					output: "build/git-source/stmr.c/stmr.go",
				},
			},
		},
		{
			prefix:    "tinyexpr",
			gitSource: "https://github.com/codeplea/tinyexpr.git",
			files: []fs{
				{
					input:  []string{"build/git-source/tinyexpr/tinyexpr.c"},
					output: "build/git-source/tinyexpr/tinyexpr.go",
				},
			},
		},
	}

	for _, tc := range tcs {
		fileList, err := getFileList(tc.prefix, tc.gitSource)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(">", fileList)

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
					fmt.Printf("`%v`:%v\n", f.output, log)
				}

			})
		}
	}
}
