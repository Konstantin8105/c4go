package cs

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	// Error tree
	errors "github.com/Konstantin8105/errors"
)

// All run all codestyle test for golang sources
//
//	Ignore data from folder "testdata"
//
func All(t *testing.T) {
	tcs := []struct {
		name string
		f    func(*testing.T)
	}{
		{"Todo", Todo},
		{"Debug", Debug},
		{"Os", Os},
	}
	for _, tc := range tcs {
		t.Run(tc.name, tc.f)
	}
}

// Todo calculate amount comments with TODO, FIX, BUG in golang sources.
//
//	Ignore data from folder "testdata"
//
func Todo(t *testing.T) {
	var amount int
	iterator(t, func(line, source string, pos int) {
		if !strings.Contains(line, "/"+"/") {
			// ignore lines witout comments
			return
		}
		if !(strings.Contains(strings.ToUpper(line), "TODO") ||
			strings.Contains(strings.ToUpper(line), "FIX") ||
			strings.Contains(strings.ToUpper(line), "BUG")) {
			// ignore lines without TODO
			return
		}

		// write result
		t.Logf("%13s:%-4d %s", source, pos, strings.TrimSpace(line))
		amount++
	})
	if amount > 0 {
		t.Logf("Amount comments: %d", amount)
	}
}

// Debug test source for avoid debug printing
//
//	Ignore data from folder "testdata"
//
func Debug(t *testing.T) {
	iterator(t, func(line, source string, pos int) {
		if !strings.Contains(line, "fmt"+"."+"Print") {
			// ignore lines without fmt Print
			return
		}
		t.Errorf("Fail: %13s:%-4d %s", source, pos, strings.TrimSpace(line))
	})

	iterator(t, func(line, source string, pos int) {
		index := strings.Index(line, "/"+"/")
		if index < 0 {
			// ignore lines without comments
			return
		}
		if !strings.Contains(line[index:], "fmt.") {
			// ignore lines without "fmt" in comments
			return
		}
		t.Logf("%13s:%-4d %s", source, pos, line[index:])
	})
}

// Os test source for avoid words "darwin", "macos"
//
//	Ignore data from folder "testdata"
//
func Os(t *testing.T) {
	iterator(t, func(line, source string, pos int) {
		if !(strings.Contains(strings.ToUpper(line), "DAR"+"WIN") ||
			strings.Contains(strings.ToUpper(line), "MAC"+"OS")) {
			return
		}
		t.Errorf("Fail: %13s:%-4d %s", source, pos, line)
	})
}

// getGoCode return all golang sources recursive
func getGoCode(dir string) (files []string, err error) {
	defer func() {
		if err != nil {
			err = errors.New(fmt.Sprintf("Cannot search go code in `%s`", dir)).Add(err)
		}
	}()
	ents, err := ioutil.ReadDir(dir)
	if err != nil {
		err = errors.New("Cannot read directory").Add(err)
		return
	}

	for _, ent := range ents {
		if ent.IsDir() {
			if ent.Name() == "testdata" {
				// ignore folder "testdata"
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

func iterator(t *testing.T, f func(line, source string, linePosition int)) {
	// search all sources
	sources, err := getGoCode("./")
	if err != nil {
		t.Fatal(err)
	}

	for i := range sources {
		// ignore folder vendor
		if strings.Contains(sources[i], "vendor"+string(os.PathSeparator)) {
			continue
		}

		t.Run(sources[i], func(t *testing.T) {
			// open file
			file, err := os.Open(sources[i])
			if err != nil {
				t.Fatal(errors.New(fmt.Sprintf("Cannot open file: %s", sources[i])).Add(err))
			}
			// close file
			defer func() {
				err := file.Close()
				if err != nil {
					t.Fatal(errors.New(fmt.Sprintf("Cannot close file: %s", sources[i])).Add(err))
				}
			}()

			// analyze by line
			pos := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				pos++
				// run function
				f(line, sources[i], pos)
			}

			// close scanner
			if err := scanner.Err(); err != nil {
				t.Fatal(errors.New(fmt.Sprintf("Scanner error in file: %s", sources[i])).Add(err))
			}
		})
	}
}
