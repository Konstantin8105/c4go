package preprocessor

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/scanner"
	"unicode"
	"unicode/utf8"

	"github.com/Konstantin8105/c4go/util"
)

// One simple part of preprocessor code
type entity struct {
	positionInSource int
	include          string
	other            string

	// Zero index of `lines` is look like that:
	// # 11 "/usr/include/x86_64-linux-gnu/gnu/stubs.h" 2 3 4
	// After that 0 or more lines of codes
	lines []*string
}

func (e *entity) parseComments(comments *[]Comment) {
	var source bytes.Buffer
	for i := range e.lines {
		if i == 0 {
			continue
		}
		source.Write([]byte(*e.lines[i]))
		source.Write([]byte{'\n'})
	}

	var s scanner.Scanner
	s.Init(strings.NewReader(source.String()))
	s.Mode = scanner.ScanComments
	s.Filename = e.include
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if scanner.TokenString(tok) == "Comment" {
			// parse multiline comments to single line comment
			var lines []string
			if s.TokenText()[1] == '*' {
				lines = strings.Split(s.TokenText(), "\n")
				lines[0] = strings.TrimLeft(lines[0], "/"+"*")
				lines[len(lines)-1] = strings.TrimRight(lines[len(lines)-1], "*"+"/")
				for i := range lines {
					lines[i] = "/" + "/" + lines[i]
				}
			} else {
				lines = append(lines, s.TokenText())
			}

			// save comments
			for _, l := range lines {
				(*comments) = append(*comments, Comment{
					File:    e.include,
					Line:    s.Position.Line + e.positionInSource - 1,
					Comment: l,
				})
			}
		}
	}
}

// isSame - check is Same entities
func (e *entity) isSame(x *entity) bool {
	if e.include != x.include {
		return false
	}
	if e.positionInSource != x.positionInSource {
		return false
	}
	if e.other != x.other {
		return false
	}
	if len(e.lines) != len(x.lines) {
		return false
	}
	for k := range e.lines {
		is := e.lines[k]
		js := x.lines[k]
		if len(*is) != len(*js) || *is != *js {
			return false
		}
	}
	return true
}

// Comment - position of line comment '//...'
type Comment struct {
	File    string
	Line    int
	Comment string
}

// IncludeHeader - struct for C include header
type IncludeHeader struct {
	HeaderName   string
	IsUserSource bool
}

// FilePP a struct with all information about preprocessor C code
type FilePP struct {
	entities []entity
	pp       []byte
	comments []Comment
	includes []IncludeHeader
}

// NewFilePP create a struct FilePP with results of analyzing
// preprocessor C code
func NewFilePP(inputFiles, clangFlags []string, cppCode bool) (
	f FilePP, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Preprocess error : %v", err)
		}
	}()

	var allItems []entity

	allItems, err = analyzeFiles(inputFiles, clangFlags, cppCode)
	if err != nil {
		return
	}

	// Generate list of user files
	userSource := map[string]bool{}
	var us []string
	us, err = GetIncludeListWithUserSource(inputFiles, clangFlags, cppCode)
	if err != nil {
		return
	}
	var all []string
	all, err = GetIncludeFullList(inputFiles, clangFlags, cppCode)
	if err != nil {
		return
	}
	// Generate C header list
	f.includes = generateIncludeList(us, all)

	for j := range us {
		userSource[us[j]] = true
	}

	// Merge the entities
	var lines []string
	for i := range allItems {
		// If found same part of preprocess code, then
		// don't include in result buffer for transpiling
		// for avoid dublicate of code
		var found bool
		for j := 0; j < i; j++ {
			if allItems[i].isSame(&allItems[j]) {
				found = true
				break
			}
		}
		if found {
			continue
		}
		// Parse comments only for user sources
		var isUserSource bool
		if userSource[allItems[i].include] {
			isUserSource = true
		}
		if allItems[i].include[0] == '.' &&
			allItems[i].include[1] == '/' &&
			userSource[allItems[i].include[2:]] {
			isUserSource = true
		}
		if isUserSource {
			allItems[i].parseComments(&f.comments)
		}

		// Parameter "other" is not included for avoid like:
		// ./tests/multi/head.h:4:28: error: invalid line marker flag '2': \
		// cannot pop empty include stack
		// # 2 "./tests/multi/main.c" 2
		//                            ^
		header := fmt.Sprintf("# %d \"%s\"",
			allItems[i].positionInSource, allItems[i].include)
		lines = append(lines, header)
		if len(allItems[i].lines) > 0 {
			for ii, l := range allItems[i].lines {
				if ii == 0 {
					continue
				}
				lines = append(lines, *l)
			}
		}
		f.entities = append(f.entities, allItems[i])
	}
	f.pp = ([]byte)(strings.Join(lines, "\n"))

	return
}

// GetSource return source of preprocessor C code
func (f FilePP) GetSource() []byte {
	return f.pp
}

// GetComments return comments in preprocessor C code
func (f FilePP) GetComments() []Comment {
	return f.comments
}

// GetIncludeFiles return list of '#include' file in C sources
func (f FilePP) GetIncludeFiles() []IncludeHeader {
	return f.includes
}

// IsUserSource get is it source from user
func (f FilePP) IsUserSource(in string) bool {
	for i := range f.includes {
		if strings.Contains(in, f.includes[i].HeaderName) &&
			f.includes[i].IsUserSource {
			return true
		}
	}
	return false
}

// GetSnippet return short part of code inside preprocessor C code
func (f FilePP) GetSnippet(file string,
	line, lineEnd int,
	col, colEnd int) (
	buffer []byte, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("GetSnippet error for `%v` {%v,%v}{%v,%v}. %v",
				file,
				line, lineEnd,
				col, colEnd,
				err)
		}
	}()

	if lineEnd == 0 {
		lineEnd = line
	}

	// replace 2,3,4... byte of rune to one byte symbol
	var t string
	for _, r := range file {
		if utf8.RuneLen(r) > 1 {
			t += "_"
			continue
		}
		t += string(r)
	}
	file = t

again:
	for i := range f.entities {
		for j := range f.entities[i].include {
			if f.entities[i].include[j] != '\\' {
				continue
			}
			if j+3 > len(f.entities[i].include)-1 {
				continue
			}
			wrongSymbol := false
			var isSymbol2 bool
			runes := f.entities[i].include[j+1 : j+4]
			for y, r := range runes {
				if !unicode.IsDigit(r) {
					wrongSymbol = true
				}
				if y == 0 && r == '2' {
					isSymbol2 = true
				}
			}
			if !wrongSymbol {
				if isSymbol2 {
					f.entities[i].include = f.entities[i].include[:j] + "_" +
						f.entities[i].include[j+4:]
				} else {
					f.entities[i].include = f.entities[i].include[:j] +
						f.entities[i].include[j+4:]
				}
				goto again
			}
		}
	}

	for i := range f.entities {
		if f.entities[i].include != file {
			continue
		}
		lineEnd := lineEnd
		if len(f.entities[i].lines)+f.entities[i].positionInSource < lineEnd {
			continue
		}
		l := f.entities[i].lines[lineEnd+1-f.entities[i].positionInSource]
		if col == 0 && colEnd == 0 {
			return []byte(*l), nil
		}
		if colEnd == 0 {
			if col-1 < len([]byte(*l)) {
				return []byte((*l)[col-1:]), nil
			}
			err = fmt.Errorf("Empty snippet")
			return
		}
		return []byte((*l)[col-1 : colEnd]), nil
	}

	err = fmt.Errorf("Snippet is not found")
	return
}

// analyzeFiles - analyze single file and separation preprocessor code to part
func analyzeFiles(inputFiles, clangFlags []string, cppCode bool) (
	items []entity, err error) {
	// See : https://clang.llvm.org/docs/CommandGuide/clang.html
	// clang -E <file>    Run the preprocessor stage.
	var out bytes.Buffer
	out, err = getPreprocessSources(inputFiles, clangFlags, cppCode)
	if err != nil {
		return
	}

	// Parsing preprocessor file
	r := bytes.NewReader(out.Bytes())
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	// counter - get position of line
	var counter int
	// item, items - entity of preprocess file
	var item *entity

	reg := util.GetRegex("# (\\d+) \".*\".*")

	for scanner.Scan() {
		line := scanner.Text()
		if reg.MatchString(line) {
			if item != (*entity)(nil) {
				items = append(items, *item)
			}
			item, err = parseIncludePreprocessorLine(line)
			if err != nil {
				err = fmt.Errorf("Cannot parse line : %s with error: %s", line, err)
				return
			}
			if item.positionInSource == 0 {
				// cannot by less 1 for avoid problem with
				// indentification of "0" AST base element
				item.positionInSource = 1
			}
			item.lines = make([]*string, 0)
		}
		counter++
		item.lines = append(item.lines, &line)
	}
	if item != (*entity)(nil) {
		items = append(items, *item)
	}
	return
}

// See : https://clang.llvm.org/docs/CommandGuide/clang.html
// clang -E <file>    Run the preprocessor stage.
func getPreprocessSources(inputFiles, clangFlags []string, cppCode bool) (
	out bytes.Buffer, err error) {
	// get temp dir
	dir, err := ioutil.TempDir("", "c4go-union")
	if err != nil {
		return
	}
	defer func() { _ = os.RemoveAll(dir) }()

	// file name union file
	var unionFileName = dir + "/" + "unionFileName.c"

	// create a body for union file
	var unionBody string
	for i := range inputFiles {
		var absPath string
		absPath, err = filepath.Abs(inputFiles[i])
		if err != nil {
			return
		}
		unionBody += fmt.Sprintf("#include \"%s\"\n", absPath)
	}

	// write a union file
	err = ioutil.WriteFile(unionFileName, []byte(unionBody), 0644)
	if err != nil {
		return
	}

	// Add open source defines
	clangFlags = append(clangFlags, "-D_GNU_SOURCE")

	// preprocessor clang
	var stderr bytes.Buffer

	var args []string
	args = append(args, "-E", "-C")
	args = append(args, clangFlags...)
	args = append(args, unionFileName) // All inputFiles

	var outFile bytes.Buffer
	var cmd *exec.Cmd

	compiler, compilerFlag := Compiler(cppCode)
	args = append(compilerFlag, args...)
	cmd = exec.Command(compiler, args...)

	cmd.Stdout = &outFile
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("preprocess for file: %v\nfailed: %v\nStdErr = %v", inputFiles, err, stderr.String())
		return
	}
	_, err = out.Write(outFile.Bytes())
	if err != nil {
		return
	}

	return
}

func generateIncludeList(userList, allList []string) (
	includes []IncludeHeader) {

	for i := range allList {
		var isUser bool
		for j := range userList {
			if allList[i] == userList[j] {
				isUser = true
				break
			}
		}
		includes = append(includes, IncludeHeader{
			HeaderName:   allList[i],
			IsUserSource: isUser,
		})
	}

	return
}

// GetIncludeListWithUserSource - Get list of include files
// Example:
// $ clang  -MM -c exit.c
// exit.o: exit.c tests.h
func GetIncludeListWithUserSource(inputFiles, clangFlags []string, cppCode bool) (
	lines []string, err error) {
	return getIncludeList(inputFiles, clangFlags, "-MM", cppCode)
}

// GetIncludeFullList - Get full list of include files
// Example:
// $ clang -M -c triangle.c
// triangle.o: triangle.c /usr/include/stdio.h /usr/include/features.h \
//   /usr/include/stdc-predef.h /usr/include/x86_64-linux-gnu/sys/cdefs.h \
//   /usr/include/x86_64-linux-gnu/bits/wordsize.h \
//   /usr/include/x86_64-linux-gnu/gnu/stubs.h \
//   /usr/include/x86_64-linux-gnu/gnu/stubs-64.h \
//   / ........ and other
func GetIncludeFullList(inputFiles, clangFlags []string, cppCode bool) (
	lines []string, err error) {
	return getIncludeList(inputFiles, clangFlags, "-M", cppCode)
}

func getIncludeList(inputFiles, clangFlags []string, flag string, cppCode bool) (
	lines []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot get Include List : %v", err)
		}
	}()
	var out bytes.Buffer
	var stderr bytes.Buffer
	var args []string
	for i := range inputFiles {
		inputFiles[i], err = filepath.Abs(inputFiles[i])
		if err != nil {
			return
		}
	}
	args = append(args, flag, "-c")
	args = append(args, inputFiles...)
	args = append(args, clangFlags...)
	var cmd *exec.Cmd

	compiler, compilerFlag := Compiler(cppCode)
	args = append(compilerFlag, args...)
	cmd = exec.Command(compiler, args...)

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("preprocess failed: %v\nStdErr = %v", err, stderr.String())
		return
	}
	return parseIncludeList(out.String())
}
