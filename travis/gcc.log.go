package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	file, err := os.Open("/tmp/gcc.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	f, err := os.Create("./travis/gsl.list")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cList map[string]bool = map[string]bool{}
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
		line = line[index+len("-c "):]
		index = strings.Index(line, " ")
		if index < 0 {
			continue
		}
		line = line[:index]
		if !strings.HasSuffix(strings.ToLower(line), ".c") {
			continue
		}

		folder := "/tmp/GSL/gsl-2.4/"
		var fileList []string
		// find all C source files
		err = filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
			if strings.HasSuffix(strings.ToLower(f.Name()), ".c") {
				if strings.HasSuffix(path, "/"+line) {
					fileList = append(fileList, path)
				}
			}
			return nil
		})
		if err != nil {
			err = fmt.Errorf("Cannot walk: %v", err)
			return
		}

		for _, f := range fileList {
			cList[f] = true
		}

		fmt.Println("line = ", line, fileList)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for k := range cList {
		fmt.Printf("%s ", k)
		f.WriteString(fmt.Sprintf("%s ", k))
	}
}
