package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func getFileList(prefix, gitSource string) (fileList []string, err error) {
	// create folder if not exist
	var temp_folder string
	temp_folder, err = ioutil.TempDir("", prefix)
	if err != nil {
		err = fmt.Errorf("Cannot create a folder : %v", err)
		return
	}

	// clone git repository
	args := []string{"clone", gitSource, temp_folder}
	err = exec.Command("git", args...).Run()
	if err != nil {
		err = fmt.Errorf("Cannot clone git repository with args `%v`: %v",
			args, err)
		return
	}

	// find all C source files
	err = filepath.Walk(temp_folder, func(path string, f os.FileInfo, err error) error {
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
						t.Fatalf("Cannot transpile `%v`", os.Args)
					}
					// logging warnings
					var err error
					var logs []string
					logs, err = getLogs(goFile)
					if err != nil {
						t.Errorf("Error in `%v`: %v", goFile, err)
					}
					for _, log := range logs {
						t.Log(log)
						fmt.Println(log)
					}
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
