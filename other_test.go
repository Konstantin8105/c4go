package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSourceVasilevBook(t *testing.T) {
	// repair arguments
	arguments := os.Args
	defer func() {
		os.Args = arguments
	}()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic is not acceptable for position: %v", r)
		}
	}()

	prefix := "VasielBook"
	git_source := "https://github.com/olegbukatchuk/book-c-the-examples-and-tasks.git"

	// create folder if not exist
	temp_folder, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatalf("Cannot create a folder : %v", err)
	}

	// clone git repository
	args := []string{"clone", git_source, temp_folder}
	if err := exec.Command("git", args...).Run(); err != nil {
		t.Fatalf("Cannot clone git repository with args `%v`: %v", args, err)
	}

	// find all C source files
	fileList := []string{}
	if err := filepath.Walk(temp_folder, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(strings.ToLower(f.Name()), ".c") {
			fileList = append(fileList, path)
		}
		return nil
	}); err != nil {
		t.Fatalf("Cannot walk: %v", err)
	}

	for _, file := range fileList {
		// black list of sources
		if strings.Contains(file, "1.13/main.c") ||
			strings.Contains(file, "1.6/main.c") ||
			strings.Contains(file, "5.9/main.c") ||
			strings.Contains(file, "3.19/main.c") ||
			strings.Contains(file, "3.17/main.c") {
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
