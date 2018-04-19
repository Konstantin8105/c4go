package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func TestSourceVasilevBook(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic is not acceptable: %v", r)
		}
	}()

	prefix := "VasielBook"
	gitSource := "https://github.com/olegbukatchuk/book-c-the-examples-and-tasks.git"

	fileList, err := getFileList(prefix, gitSource)
	if err != nil {
		t.Fatal(err)
	}

	ignoreFileList := []string{
		"1.13/main.c",
		"1.6/main.c",
		"5.9/main.c",
		"3.19/main.c",
		"3.17/main.c",
	}

	for _, file := range fileList {
		// ignore list of sources
		var ignored bool
		for _, ignore := range ignoreFileList {
			if strings.Contains(strings.ToLower(file), strings.ToLower(ignore)) {
				ignored = true
			}
		}
		if ignored {
			continue
		}

		// run test
		t.Run(file, func(t *testing.T) {
			file = strings.TrimSpace(file)
			os.Args = []string{"c4go", "transpile", "-o=" + file + ".go", file}
			code := runCommand()
			if code != 0 {
				t.Fatalf("Cannot transpile `%v`", os.Args)
			}
		})
	}
}

func TestKRSourceBook(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic is not acceptable: %v", r)
		}
	}()

	prefix := "KR"
	gitSource := "https://github.com/KushalP/k-and-r.git"

	fileList, err := getFileList(prefix, gitSource)
	if err != nil {
		t.Fatal(err)
	}

	ignoreFileList := []string{
		"4.1-1.c",
		"4-11.c",
		"1.9-1.c",
		"1.10-1.c",
		"1-24.c",
		"1-17.c",
		"1-16.c",
		"4-10.c",
	}

	for _, file := range fileList {
		// ignore list of sources
		var ignored bool
		for _, ignore := range ignoreFileList {
			if strings.Contains(strings.ToLower(file), strings.ToLower(ignore)) {
				ignored = true
			}
		}
		if ignored {
			continue
		}

		// run test
		t.Run(file, func(t *testing.T) {
			file = strings.TrimSpace(file)
			os.Args = []string{"c4go", "transpile", "-o=" + file + ".go", file}
			code := runCommand()
			if code != 0 {
				t.Fatalf("Cannot transpile `%v`", os.Args)
			}
		})
	}
}

func TestSourceKochanBook(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic is not acceptable: %v", r)
		}
	}()

	prefix := "KochanBook"
	gitSource := "https://github.com/eugenetriguba/programming-in-c.git"

	fileList, err := getFileList(prefix, gitSource)
	if err != nil {
		t.Fatal(err)
	}

	ignoreFileList := []string{
		"5.9d.c",
		"5.9c.c",
	}

	for _, file := range fileList {
		// ignore list of sources
		var ignored bool
		for _, ignore := range ignoreFileList {
			if strings.Contains(strings.ToLower(file), strings.ToLower(ignore)) {
				ignored = true
			}
		}
		if ignored {
			continue
		}

		// run test
		t.Run(file, func(t *testing.T) {
			file = strings.TrimSpace(file)
			os.Args = []string{"c4go", "transpile", "-o=" + file + ".go", file}
			code := runCommand()
			if code != 0 {
				t.Fatalf("Cannot transpile `%v`", os.Args)
			}
		})
	}
}

func TestSourceDeitelBook(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic is not acceptable: %v", r)
		}
	}()

	prefix := "DeitelBook"
	gitSource := "https://github.com/Emmetttt/C-Deitel-Book.git"

	fileList, err := getFileList(prefix, gitSource)
	if err != nil {
		t.Fatal(err)
	}

	ignoreFileList := []string{
		"E5.45.C",
		"06.14_const_type_qualifier.C",
		"E7.17.C",
	}

	for _, file := range fileList {
		// ignore list of sources
		var ignored bool
		for _, ignore := range ignoreFileList {
			if strings.Contains(strings.ToLower(file), strings.ToLower(ignore)) {
				ignored = true
			}
		}
		if ignored {
			continue
		}

		// run test
		t.Run(file, func(t *testing.T) {
			file = strings.TrimSpace(file)
			os.Args = []string{"c4go", "transpile", "-o=" + file + ".go", file}
			code := runCommand()
			if code != 0 {
				t.Fatalf("Cannot transpile `%v`", os.Args)
			}
		})
	}
}
