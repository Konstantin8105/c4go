package noarch

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

// Programs generated by c4go will reference noarch.Stdin instead of os.Stdin
// directly so that under test these can be replaced. This is required because
// "go test" does not redirect the stdin to the executable it is testing.
var (
	Stdin  = NewFile(os.Stdin)
	Stdout = NewFile(os.Stdout)
	Stderr = NewFile(os.Stderr)
)

// File represents the definition has been translated from the original
// definition for __sFILE, which is an alias for FILE. Not all of the attributes
// have been translated. They should be turned on as needed.
type File struct {
	// This is not part of the original struct but it is needed for internal
	// calls in Go.
	OsFile *os.File

	// unsigned char *_p;
	// int _r;
	// int _w;
	// short _flags;
	// short _file;
	// struct __sbuf _bf;
	// int _lbfsize;
	// void *_cookie;
	// int (* _Nullable _close)(void *);
	// int (* _Nullable _read) (void *, char *, int);
	// fpos_t (* _Nullable _seek) (void *, fpos_t, int);
	// int (* _Nullable _write)(void *, const char *, int);
	// struct __sbuf _ub;
	// struct __sFILEX *_extra;
	// int _ur;
	// unsigned char _ubuf[3];
	// unsigned char _nbuf[1];
	// struct __sbuf _lb;
	// int _blksize;
	// fpos_t _offset;
}

// Fopen handles fopen().
//
// Opens the file whose name is specified in the parameter filePath and
// associates it with a stream that can be identified in future operations by
// the File pointer returned.
//
// The operations that are allowed on the stream and how these are performed are
// defined by the mode parameter.
//
// The returned pointer can be disassociated from the file by calling fclose()
// or freopen(). All opened files are automatically closed on normal program
// termination.
func Fopen(filePath, mode []byte) *File {
	var file *os.File
	var err error

	sFilePath := CStringToString(filePath)

	// TODO: Only some modes are supported by fopen()
	// https://github.com/Konstantin8105/c4go/issues/89
	switch CStringToString(mode) {
	case "r":
		file, err = os.OpenFile(sFilePath, os.O_RDONLY, 0655)
	case "r+":
		file, err = os.OpenFile(sFilePath, os.O_RDWR, 0655)
	case "a":
		file, err = os.OpenFile(sFilePath, os.O_WRONLY|os.O_APPEND, 0655)
	case "a+":
		file, err = os.OpenFile(sFilePath, os.O_RDWR|os.O_APPEND, 0655)
	case "w":
		file, err = os.OpenFile(sFilePath, os.O_RDWR|os.O_CREATE, 0655)
	case "w+":
		file, err = os.OpenFile(sFilePath, os.O_RDWR|os.O_CREATE, 0655)
	default:
		panic(fmt.Sprintf("unsupported file mode: %s", mode))
	}

	if err != nil {
		return nil
	}

	return NewFile(file)
}

// Fclose handles fclose().
//
// Closes the file associated with the stream and disassociates it.
//
// All internal buffers associated with the stream are disassociated from it and
// flushed: the content of any unwritten output buffer is written and the
// content of any unread input buffer is discarded.
//
// Even if the call fails, the stream passed as parameter will no longer be
// associated with the file nor its buffers.
func Fclose(f *File) int32 {
	err := f.OsFile.Close()
	if err != nil {
		// Is this the correct error code?
		return 1
	}

	return 0
}

// Remove handles remove().
//
// Deletes the file whose name is specified in filePath.
//
// This is an operation performed directly on a file identified by its filePath;
// No streams are involved in the operation.
//
// Proper file access shall be available.
func Remove(filePath []byte) int32 {
	if os.Remove(CStringToString(filePath)) != nil {
		return -1
	}

	return 0
}

// Rename handles rename().
//
// Changes the name of the file or directory specified by oldName to newName.
//
// This is an operation performed directly on a file; No streams are involved in
// the operation.
//
// If oldName and newName specify different paths and this is supported by the
// system, the file is moved to the new location.
//
// If newName names an existing file, the function may either fail or override
// the existing file, depending on the specific system and library
// implementation.
//
// Proper file access shall be available.
func Rename(oldName, newName []byte) int32 {
	from := CStringToString(oldName)
	to := CStringToString(newName)

	if os.Rename(from, to) != nil {
		return -1
	}

	return 0
}

// Fputs handles fputs().
//
// Writes the C string pointed by str to the stream.
//
// The function begins copying from the address specified (str) until it reaches
// the terminating null character ('\0'). This terminating null-character is not
// copied to the stream.
//
// Notice that fputs not only differs from puts in that the destination stream
// can be specified, but also fputs does not write additional characters, while
// puts appends a newline character at the end automatically.
func Fputs(str []byte, stream *File) int32 {
	if stream == nil {
		return -1
	}

	s := CStringToString(str)

	n, err := stream.OsFile.WriteString(s)
	if err != nil {
		panic(err)
	}

	return int32(n)
}

// Tmpfile handles tmpfile().
//
// Creates a temporary binary file, open for update ("wb+" mode, see fopen for
// details) with a filename guaranteed to be different from any other existing
// file.
//
// The temporary file created is automatically deleted when the stream is closed
// (fclose) or when the program terminates normally. If the program terminates
// abnormally, whether the file is deleted depends on the specific system and
// library implementation.
func Tmpfile() *File {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil
	}

	return NewFile(f)
}

// Fgets handles fgets().
//
// Reads characters from stream and stores them as a C string into str until
// (num-1) characters have been read or either a newline or the end-of-file is
// reached, whichever happens first.
//
// A newline character makes fgets stop reading, but it is considered a valid
// character by the function and included in the string copied to str.
//
// A terminating null character is automatically appended after the characters
// copied to str.
//
// Notice that fgets is quite different from gets: not only fgets accepts a
// stream argument, but also allows to specify the maximum size of str and
// includes in the string any ending newline character.
func Fgets(str []byte, num32 int32, stream *File) []byte {
	num := int(num32)
	buf := make([]byte, num)
	n, err := stream.OsFile.Read(buf)

	// FIXME: Is this the right thing to do in this case?
	if err != nil {
		return []byte{}
	}

	// TODO: Allow arguments to be passed by reference.
	// https://github.com/Konstantin8105/c4go/issues/90
	// This appears in multiple locations.

	// Be careful to crop the buffer to the real number of bytes read.
	//
	// We do not trim off the NULL characters because we do not know if the file
	// we are reading is in binary mode.
	if n == num {
		// If it is the case that we have read the entire buffer with this read
		// we need to make sure we leave room for what would be the NULL
		// character at the end of the string in C.
		return buf[:n-1]
	}

	return buf[:n]
}

// Gets read bytes from stdin
func Gets(str []byte) []byte {
	_, err := Stdin.OsFile.Read(str)

	// FIXME: Is this the right thing to do in this case?
	if err != nil {
		return []byte{}
	}

	return str
}

// Rewind handles rewind().
//
// Sets the position indicator associated with stream to the beginning of the
// file.
//
// The end-of-file and error internal indicators associated to the stream are
// cleared after a successful call to this function, and all effects from
// previous calls to ungetc on this stream are dropped.
//
// On streams open for update (read+write), a call to rewind allows to switch
// between reading and writing.
func Rewind(stream *File) {
	stream.OsFile.Seek(0, 0)
}

// Feof handles feof().
//
// Checks whether the end-of-File indicator associated with stream is set,
// returning a value different from zero if it is.
//
// This indicator is generally set by a previous operation on the stream that
// attempted to read at or past the end-of-file.
//
// Notice that stream's internal position indicator may point to the end-of-file
// for the next operation, but still, the end-of-file indicator may not be set
// until an operation attempts to read at that point.
//
// This indicator is cleared by a call to clearerr, rewind, fseek, fsetpos or
// freopen. Although if the position indicator is not repositioned by such a
// call, the next i/o operation is likely to set the indicator again.
func Feof(stream *File) int32 {
	// FIXME: This is a really bad way of doing this. Basically try and peek
	// ahead to test for EOF.
	buf := make([]byte, 1)
	_, err := stream.OsFile.Read(buf)

	result := 0
	if err == io.EOF {
		result = 1
	}

	// Undo cursor before returning.
	stream.OsFile.Seek(-1, 1)

	return int32(result)
}

// NewFile creates a File pointer from a Go file pointer.
func NewFile(f *os.File) *File {
	return &File{
		OsFile: f,
	}
}

// Tmpnam handles tmpnam().
//
// Returns a string containing a file name different from the name of any
// existing file, and thus suitable to safely create a temporary file without
// risking to overwrite an existing file.
//
// If str is a null pointer, the resulting string is stored in an internal
// static array that can be accessed by the return value. The content of this
// string is preserved at least until a subsequent call to this same function,
// which may overwrite it.
//
// If str is not a null pointer, it shall point to an array of at least L_tmpnam
// characters that will be filled with the proposed temporary file name.
//
// The file name returned by this function can be used to create a regular file
// using fopen to be used as a temporary file. The file created this way, unlike
// those created with tmpfile is not automatically deleted when closed; A
// program shall call remove to delete this file once closed.
func Tmpnam(str []byte) []byte {
	// TODO: Allow arguments to be passed by reference.
	// https://github.com/Konstantin8105/c4go/issues/90
	// This appears in multiple locations.

	// TODO: There must be a better way of doing this. This way allows the same
	// great distinct Go temp file generation (that also checks for existing
	// files), but unfortunately creates the file in the process; even if you
	// don't intend to use it.
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return []byte{}
	}

	f.Close()
	return []byte(f.Name())
}

// Fflush handles fflush().
//
// If the given stream was open for writing (or if it was open for updating and
// the last i/o operation was an output operation) any unwritten data in its
// output buffer is written to the file.
//
// If stream is a null pointer, all such streams should be flushed, but this is
// currently not supported.
//
// The stream remains open after this call.
//
// When a file is closed, either because of a call to fclose or because the
// program terminates, all the buffers associated with it are automatically
// flushed.
func Fflush(stream *File) int32 {
	err := stream.OsFile.Sync()
	if err != nil {
		return 1
	}

	return 0
}

// Fprintf handles fprintf().
//
// Writes the C string pointed by format to the stream. If format includes
// format specifiers (subsequences beginning with %), the additional arguments
// following format are formatted and inserted in the resulting string replacing
// their respective specifiers.
//
// After the format parameter, the function expects at least as many additional
// arguments as specified by format.
func Fprintf(f *File, format []byte, args ...interface{}) int32 {
	realArgs := []interface{}{}

	// Convert any C strings into Go strings.
	typeOfByteSlice := reflect.TypeOf([]byte(nil))
	for _, arg := range args {
		if reflect.TypeOf(arg) == typeOfByteSlice {
			realArgs = append(realArgs, CStringToString(arg.([]byte)))
		} else {
			realArgs = append(realArgs, arg)
		}
	}

	n, err := fmt.Fprintf(f.OsFile, CStringToString(format), realArgs...)
	if err != nil {
		return -1
	}

	return int32(n)
}

// Fscanf handles fscanf().
//
// Reads data from the stream and stores them according to the parameter format
// into the locations pointed by the additional arguments.
//
// The additional arguments should point to already allocated objects of the
// type specified by their corresponding format specifier within the format
// string.
func Fscanf(f *File, format []byte, args ...interface{}) int32 {
	realArgs := prepareArgsForScanf(args)

	goFormat := CStringToString(format)

again:
	n, err := fmt.Fscanf(f.OsFile, goFormat, realArgs...)
	if err != nil {
		if strings.Contains(err.Error(), "unexpected newline") {
			goto again
		}
		return -1
	}

	finalizeArgsForScanf(realArgs, args)

	return int32(n)
}

func finalizeArgsForScanf(realArgs []interface{}, args []interface{}) {
	typeOfStringRef := reflect.TypeOf(new(string))
	for i, arg := range realArgs {
		if reflect.TypeOf(arg) == typeOfStringRef {
			s := *arg.(*string)
			copy(args[i].([]byte), []byte(s+"\x00"))
		} else {
			GoPointerToCPointer(arg, args[i])
		}
	}
}

func prepareArgsForScanf(args []interface{}) []interface{} {
	realArgs := []interface{}{}

	typeOfByteSlice := reflect.TypeOf([]byte(nil))
	for _, arg := range args {
		if reflect.TypeOf(arg) == typeOfByteSlice {
			realArgs = append(realArgs, new(string))
		} else {
			realArgs = append(realArgs, CPointerToGoPointer(arg))
		}
	}

	return realArgs
}

func getc(f *os.File) int32 {
	buffer := make([]byte, 1)
	_, err := f.Read(buffer)
	if err != nil {
		return -1
	}

	return int32(buffer[0])
}

// Fgetc handles fgetc().
//
// Returns the character currently pointed by the internal file position
// indicator of the specified stream. The internal file position indicator is
// then advanced to the next character.
//
// If the stream is at the end-of-file when called, the function returns EOF and
// sets the end-of-file indicator for the stream (feof).
//
// If a read error occurs, the function returns EOF and sets the error indicator
// for the stream (ferror).
//
// fgetc and getc are equivalent, except that getc may be implemented as a macro
// in some libraries.
func Fgetc(stream *File) int32 {
	return int32(getc(stream.OsFile))
}

// Fputc handles fputc().
//
// Writes a character to the stream and advances the position indicator.
//
// The character is written at the position indicated by the internal position
// indicator of the stream, which is then automatically advanced by one.
func Fputc(c int32, f *File) int32 {
	if f == nil {
		return -1
	}
	n, err := f.OsFile.Write([]byte{byte(c)})
	if err != nil {
		return 0
	}

	return int32(n)
}

// Getchar handles getchar().
//
// Returns the next character from the standard input (stdin).
//
// It is equivalent to calling getc with stdin as argument.
func Getchar() int32 {
	return int32(getc(Stdin.OsFile))
}

// Fseek handles fseek().
//
// Sets the position indicator associated with the stream to a new position.
//
// For streams open in binary mode, the new position is defined by adding offset
// to a reference position specified by origin.
//
// For streams open in text mode, offset shall either be zero or a value
// returned by a previous call to ftell, and origin shall necessarily be
// SEEK_SET.
//
// If the function is called with other values for these arguments, support
// depends on the particular system and library implementation (non-portable).
//
// The end-of-file internal indicator of the stream is cleared after a
// successful call to this function, and all effects from previous calls to
// ungetc on this stream are dropped.
//
// On streams open for update (read+write), a call to fseek allows to switch
// between reading and writing.
func Fseek(f *File, offset int32, origin int32) int32 {
	n, err := f.OsFile.Seek(int64(offset), int(origin))
	if err != nil {
		return -1
	}

	return int32(n)
}

// Ftell handles ftell().
//
// Returns the current value of the position indicator of the stream.
//
// For binary streams, this is the number of bytes from the beginning of the
// file.
//
// For text streams, the numerical value may not be meaningful but can still be
// used to restore the position to the same position later using fseek (if there
// are characters put back using ungetc still pending of being read, the
// behavior is undefined).
func Ftell(f *File) int32 {
	return int32(Fseek(f, 0, 1))
}

// Fread handles fread().
//
// Reads an array of count elements, each one with a size of size bytes, from
// the stream and stores them in the block of memory specified by ptr.
//
// The position indicator of the stream is advanced by the total amount of bytes
// read.
//
// The total amount of bytes read if successful is (size*count).
func Fread(ptr *[]byte, size1, size2 int32, f *File) int32 {
	// Create a new buffer so that we can ensure we read up to the correct
	// number of bytes from the file.
	newBuffer := make([]byte, size1*size2)
	n, err := f.OsFile.Read(newBuffer)

	// Despite any error we need to make sure the bytes read are copied to the
	// destination buffer.
	for i, b := range newBuffer {
		(*ptr)[i] = b
	}

	// Now we can handle the success or failure.
	if err != nil {
		return -1
	}

	return int32(n)
}

// Fwrite handles fwrite().
//
// Writes an array of count elements, each one with a size of size bytes, from
// the block of memory pointed by ptr to the current position in the stream.
//
// The position indicator of the stream is advanced by the total number of bytes
// written.
//
// Internally, the function interprets the block pointed by ptr as if it was an
// array of (size*count) elements of type unsigned char, and writes them
// sequentially to stream as if fputc was called for each byte.
func Fwrite(str []byte, size1, size2 int32, stream *File) int32 {
	n, err := stream.OsFile.Write(str[:int(size1*size2)])
	if err != nil {
		return -1
	}

	return int32(n)
}

// Fgetpos handles fgetpos().
//
// Retrieves the current position in the stream.
//
// The function fills the fpos_t object pointed by pos with the information
// needed from the stream's position indicator to restore the stream to its
// current position (and multibyte state, if wide-oriented) with a call to
// fsetpos.
//
// The ftell function can be used to retrieve the current position in the stream
//as an integer value.
func Fgetpos(f *File, pos []int32) int32 {
	absolutePos := int32(Fseek(f, 0, 1))
	if pos != nil {
		pos[0] = absolutePos
	}

	return absolutePos
}

// Fsetpos handles fsetpos().
//
// Restores the current position in the stream to pos.
//
// The internal file position indicator associated with stream is set to the
// position represented by pos, which is a pointer to an fpos_t object whose
// value shall have been previously obtained by a call to fgetpos.
//
// The end-of-file internal indicator of the stream is cleared after a
// successful call to this function, and all effects from previous calls to
// ungetc on this stream are dropped.
//
// On streams open for update (read+write), a call to fsetpos allows to switch
// between reading and writing.
//
// A similar function, fseek, can be used to set arbitrary positions on streams
// open in binary mode.
func Fsetpos(stream *File, pos []int32) int32 {
	return int32(Fseek(stream, int32(pos[0]), 0))
}

// Printf handles printf().
//
// Writes the C string pointed by format to the standard output (stdout). If
// format includes format specifiers (subsequences beginning with %), the
// additional arguments following format are formatted and inserted in the
// resulting string replacing their respective specifiers.
func Printf(format []byte, args ...interface{}) int32 {
	realArgs := []interface{}{}

	// Convert any C strings into Go strings.
	typeOfByteSlice := reflect.TypeOf([]byte(nil))
	for _, arg := range args {
		if reflect.TypeOf(arg) == typeOfByteSlice {
			realArgs = append(realArgs, CStringToString(arg.([]byte)))
		} else {
			realArgs = append(realArgs, arg)
		}
	}

	n, _ := fmt.Fprintf(os.Stdout, CStringToString(format), realArgs...)

	return int32(n)
}

// Puts handles puts().
//
// Writes the C string pointed by str to the standard output (stdout) and
// appends a newline character ('\n').
//
// The function begins copying from the address specified (str) until it reaches
// the terminating null character ('\0'). This terminating null-character is not
// copied to the stream.
//
// Notice that puts not only differs from fputs in that it uses stdout as
// destination, but it also appends a newline character at the end automatically
// (which fputs does not).
func Puts(str []byte) int32 {
	n, _ := fmt.Fprintln(os.Stdout, CStringToString(str))

	return int32(n)
}

// Scanf handles scanf().
//
// Reads data from stdin and stores them according to the parameter format into
// the locations pointed by the additional arguments.
//
// The additional arguments should point to already allocated objects of the
// type specified by their corresponding format specifier within the format
// string.
func Scanf(format []byte, args ...interface{}) int32 {
	realArgs := prepareArgsForScanf(args)

	// We cannot use fmt Scanf() here because that would use the real stdin
	// which does not work under test. See docs for noarch.Stdin.
	n, _ := fmt.Fscanf(Stdin.OsFile, CStringToString(format), realArgs...)
	finalizeArgsForScanf(realArgs, args)

	return int32(n)
}

func Sscanf(str []byte, format []byte, args ...interface{}) int32 {
	wws := bytes.ReplaceAll(str, []byte("\n"), []byte{' '}) // without newline
	realArgs := prepareArgsForScanf(args)
	n, _ := fmt.Sscanf(CStringToString(wws), CStringToString(format), realArgs...)
	finalizeArgsForScanf(realArgs, args)
	return int32(n)
}

// Putchar handles putchar().
//
// Writes a character to the standard output (stdout).
//
// It is equivalent to calling putc with stdout as second argument.
func Putchar(character int32) {
	fmt.Fprintf(os.Stdout, "%c", int(character))
}

// Sprintf handles sprintf().
//
// Writes the C string pointed by format to the standard output (stdout). If
// format includes format specifiers (subsequences beginning with %), the
// additional arguments following format are formatted and inserted in the
// resulting string replacing their respective specifiers.
func Sprintf(buffer, format []byte, args ...interface{}) int32 {
	realArgs := []interface{}{}

	realArgs = append(realArgs, convert(args)...)

	result := fmt.Sprintf(CStringToString(format), realArgs...)
	for i := range []byte(result) {
		buffer[i] = result[i]
	}
	buffer[len(result)] = '\x00'

	n := len(result)
	return int32(n)
}

// Vsprintf handles vsprintf().
//
// Writes the C string pointed by format to the standard output (stdout). If
// format includes format specifiers (subsequences beginning with %), the
// additional arguments following format are formatted and inserted in the
// resulting string replacing their respective specifiers.
func Vsprintf(buffer, format []byte, varList ...interface{}) int32 {
	realArgs := []interface{}{}
	realArgs = append(realArgs, convert(varList)...)

	result := fmt.Sprintf(CStringToString(format), realArgs...)
	for i := range []byte(result) {
		buffer[i] = result[i]
	}
	buffer[len(result)] = '\x00'

	n := int32(len(result))
	return n
}

func Vprintf(format []byte, varList ...interface{}) int32 {
	buffer := make([]byte, 1000)
	n := Vsprintf(buffer, format, varList...)
	buf := CStringToString(buffer)
	fmt.Fprintf(os.Stdout, "%s", buf)
	return n
}

func Vfprintf(f *File, format []byte, varList ...interface{}) int32 {
	buffer := make([]byte, 1000)
	n := Vsprintf(buffer, format, varList...)
	buf := CStringToString(buffer)
	fmt.Fprintf(f.OsFile, "%s", buf)
	return n
}

// Snprintf handles snprintf().
//
// Writes the C string pointed by format to the standard output (stdout). If
// format includes format specifiers (subsequences beginning with %), the
// additional arguments following format are formatted and inserted in the
// resulting string replacing their respective specifiers.
func Snprintf(buffer []byte, n int32, format []byte, args ...interface{}) int32 {
	return Vsnprintf(buffer, n, format, args...)
}

// convert - convert va_list
func convert(args []interface{}) (result []interface{}) {
	for i := range args {
		arg := args[i]

		typeOfByteSlice := reflect.TypeOf([]byte(nil))
		if reflect.TypeOf(arg) == typeOfByteSlice {
			result = append(result, CStringToString(arg.([]byte)))
			continue
		}

		switch v := arg.(type) {
		case int32:
			result = append(result, v)
			continue
		case string:
			result = append(result, v)
			continue
		case float64:
			result = append(result, v)
			continue
		}

		// here come &main.va_list{position:0, slice:[]interface {}{2}}
		if reflect.TypeOf(arg).Kind() == reflect.Ptr {
			val := reflect.ValueOf(arg)
			v := reflect.Indirect(val)
			if v.NumField() == 2 {
				if v.Field(1).Type().String() == "[]interface {}" {
					field := v.Field(1)
					slice := field.Interface().([]interface{})
					for i := range slice {
						switch v := slice[i].(type) {
						case []byte:
							result = append(result, CStringToString(v))

						default:
							result = append(result, slice[i])
						}
					}
				}
			}

			continue
		}

		result = append(result, arg)
	}
	return
}

// Vsnprintf handles vsnprintf().
//
// Writes the C string pointed by format to the standard output (stdout). If
// format includes format specifiers (subsequences beginning with %), the
// additional arguments following format are formatted and inserted in the
// resulting string replacing their respective specifiers.
func Vsnprintf(buffer []byte, n int32, format []byte, varList ...interface{}) int32 {
	realArgs := []interface{}{}
	realArgs = append(realArgs, convert(varList)...)

	var result string
	if len(realArgs) > 0 {
		result = fmt.Sprintf(CStringToString(format), realArgs...)
	} else {
		result = fmt.Sprintf(CStringToString(format))
	}
	if int32(len(result)) > n {
		result = result[:n]
	}
	for i := range []byte(result) {
		buffer[i] = result[i]
	}
	if len(buffer) <= len(result) {
		buffer = append(buffer, '\x00')
	} else {
		buffer[len(result)] = '\x00'
	}

	n = int32(len(result))
	return n
}

func Perror(msg []byte) {
	m := CStringToString(msg)
	fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", m)
}

func Getline(line [][]byte, len []uint32, f *File) SsizeT {
	counter := 0
	for {
		buf := make([]byte, 1)
		n, err := f.OsFile.Read(buf)
		buf = buf[:n]
		if err == io.EOF || err != nil {
			break
		}
		for i := range buf {
			line[0] = append(line[0], []byte{buf[i]}...)
		}
		counter += n
		if buf[n-1] == '\n' {
			break
		}
	}
	if counter == 0 {
		return SsizeT(-1)
	}
	line[0] = append(line[0], '\x00')
	return SsizeT(counter)
}
