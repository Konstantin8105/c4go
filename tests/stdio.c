// This program actually still works without including stdio.h but it should be
// here for consistency.

#include "tests.h"
#include <assert.h>
#include <stdarg.h>
#include <stdio.h>
#include <string.h>

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

// size of that file
int filesize = 32;
char* test_file = "tests/stdio.txt";

void test_putchar()
{
    putchar('#');
    char c;
    for (c = 'A'; c <= 'Z'; c++)
        putchar(c);
    putchar('\n');

    pass("%s", "putchar");
}

void test_puts()
{
    puts("#c4go");

    pass("%s", "puts");
}

void test_printf()
{
    // TODO: printf() has a different syntax to Go
    // https://github.com/Konstantin8105/c4go/issues/94

    printf("# Characters: %c %c \n", 'a', 65);
    //printf("Decimals: %d %ld\n", 1977, 650000L);
    printf("# Preceding with blanks: %10d \n", 1977);
    printf("# Preceding with zeros: %010d \n", 1977);
    printf("# Some different radices: %d %x %o %#x %#o \n", 100, 100, 100, 100, 100);
    printf("# floats: %4.2f %+.0e %E \n", 3.1416, 3.1416, 3.1416);
    printf("# Width trick: %*d \n", 5, 10);
    printf("# %s \n", "A string");

    // Type printf
    unsigned long long ull = 42;
    printf("Long long : %u\n", ull);
    unsigned int ui = 42;
    printf("uint : %u\n", ui);
    double d = 42.42;
    printf("double : %f\n", d);
    printf("double : %lf\n", d);
    printf("double : %2.2f\n", d);

    // Literal error
    printf("Символьный литерал \'A\':\t%d\n", sizeof('A'));

    int magnitude = 4;
    char printfFormat[30] = "%0";
    char magnitudeString[10];
    sprintf(magnitudeString, "%d", magnitude);
    strcat(printfFormat, magnitudeString);
    strcat(printfFormat, "d  ");
    printf("# ");
    printf(printfFormat, 120);
    printf(" \n");

    pass("%s", "printf");
}

void test_remove()
{
    // TODO: This does not actually test successfully deleting a file.
    if (remove("myfile.txt") != 0) {
        pass("%s", "error deleting file");
    } else {
        fail("%s", "file successfully deleted");
    }
}

void test_rename()
{
    // TODO: This does not actually test successfully renaming a file.
    int result;
    char oldname[] = "oldname.txt";
    char newname[] = "newname.txt";
    result = rename(oldname, newname);
    if (result == 0) {
        fail("%s", "File successfully renamed");
    } else {
        pass("%s", "Error renaming file");
    }
}

void test_fopen()
{
    FILE* pFile;
    pFile = fopen("./testdata/myfile.txt", "w");
    if (pFile != NULL) {
        is_not_null(pFile);
        fclose(pFile);
    }
}

void test_tmpfile()
{
    char buffer[256];
    FILE* pFile;
    pFile = tmpfile();

    fputs("hello world", pFile);
    rewind(pFile);
    fputs(fgets(buffer, 256, pFile), stdout);
    fclose(pFile);
}

void test_tmpnam()
{
    // TODO: This is a tricky one to test because the output of tmpnam() in C
    // and Go will be different. I will keep the test here so at least it tries
    // to run the code but the test itself is not actually proving anything.

    char* pointer;

    // FIXME: We cannot pass variables by reference yet, which is a legitimate
    // way to use tmpnam(). I have to leave this disabled for now.
    //
    //     char buffer[L_tmpnam];
    //     tmpnam(buffer);
    //     assert(buffer != NULL);

    pointer = tmpnam(NULL);
    is_not_null(pointer);
}

void test_fclose()
{
    remove("./testdata/myfile.txt");
    FILE* pFile;
    pFile = fopen("./testdata/myfile.txt", "w");
    fputs("fclose example", pFile);
    fclose(pFile);
    // remove temp file
    is_eq(remove("./testdata/myfile.txt"), 0)
}

void test_fflush()
{
    char mybuffer[80];
    FILE* pFile;
    pFile = fopen("./testdata/example.txt", "w+");
    is_not_null(pFile) or_return();

    fputs("test", pFile);
    fflush(pFile); // flushing or repositioning required
    fgets(mybuffer, 80, pFile);
    fclose(pFile);
    // remove temp file
    is_eq(remove("./testdata/example.txt"), 0)
}

void test_fprintf()
{
    remove("./testdata/myfile1.txt");
    FILE* pFile;
    int n;
    char* name = "John Smith";

    pFile = fopen("./testdata/myfile1.txt", "w");
    is_not_null(pFile);

    for (n = 0; n < 3; n++) {
        fprintf(pFile, "Name %d [%-10.10s]\n", n + 1, name);
    }

    fclose(pFile);
    // remove temp file
    is_eq(remove("./testdata/myfile1.txt"), 0)
}

void test_fscanf()
{
    remove("./testdata/myfile2.txt");

    char str[80];
    char end[80];
    float f;
    int i;
    FILE* pFile;

    pFile = fopen("./testdata/myfile2.txt", "w+");
    is_not_null(pFile);

    fprintf(pFile, "%f \r\n\t\n %s %d %s", 3.1416, "PI", 42, "end");
    rewind(pFile);
    fscanf(pFile, "%f", &f);
    fscanf(pFile, "%s", str);
    fscanf(pFile, "%d", &i);
    fscanf(pFile, "%s", end);
    fclose(pFile);
    pFile = NULL;

    is_eq(f, 3.1416);
    is_streq(str, "PI");
    is_eq(i, 42);
    is_streq(end, "end");

    // read again
    FILE* pFile2;
    pFile2 = fopen("./testdata/myfile2.txt", "r");
    is_not_null(pFile2);

    fscanf(pFile2, "%f", &f);
    fscanf(pFile2, "%s", str);
    fscanf(pFile2, "%d", &i);
    fscanf(pFile2, "%s", end);
    fclose(pFile2);
    pFile2 = NULL;

    is_eq(f, 3.1416);
    is_streq(str, "PI");
    is_eq(i, 42);
    is_streq(end, "end");

    // remove temp file
    is_eq(remove("./testdata/myfile2.txt"), 0)

        // test file fscan.txt
        FILE* in
        = fopen("./tests/stdio_fscan.txt", "r");
    char dummy[128];

    for (int iter = 0; iter < 10; iter++) {
        int e = fscanf(in, "%s", dummy);
        printf("fscan : iter[%d] : err = %d\n", iter, e);
        if (e < 0) {
            break;
        }
        printf("fscan : iter[%d] : str = %s\n", iter, dummy);
    }
}

void test_fgetc()
{
    FILE* pFile;
    int c;
    int n = 0;
    pFile = fopen("tests/stdio.c", "r");
    is_not_null(pFile);

    do {
        c = fgetc(pFile);
        if (c == '$')
            n++;
    } while (c != EOF);
    fclose(pFile);

    is_eq(n, 2);
}

void test_fgets()
{
    FILE* pFile;
    char* mystring;
    char dummy[20];

    pFile = fopen("tests/stdio.c", "r");
    is_not_null(pFile);

    mystring = fgets(dummy, 20, pFile);
    is_not_null(mystring);

    fclose(pFile);
}

void test_fputc()
{
    char c;

    fputc('#', stdout);
    for (c = 'A'; c <= 'Z'; c++)
        fputc(c, stdout);
    fputc('\n', stdout);
}

void test_fputs()
{
    FILE* pFile;
    char* sentence = "Hello, World";

    pFile = fopen("./testdata/mylog.txt", "w");
    fputs(sentence, pFile);
    fclose(pFile);
    // remove temp file
    is_eq(remove("./testdata/mylog.txt"), 0)
}

void test_getc()
{
    FILE* pFile;
    int c;
    int n = 0;
    pFile = fopen("tests/stdio.c", "r");
    is_not_null(pFile);

    do {
        c = getc(pFile);
        if (c == '$')
            n++;
    } while (c != EOF);
    fclose(pFile);

    is_eq(n, 2);
}

void test_putc()
{
    FILE* pFile;
    char c;

    pFile = fopen("./testdata/whatever.txt", "w");
    for (c = 'A'; c <= 'Z'; c++) {
        putc(c, pFile);
    }
    fclose(pFile);
    // remove temp file
    is_eq(remove("./testdata/whatever.txt"), 0)
}

void test_fseek()
{
    FILE* pFile;
    pFile = fopen("./testdata/example.txt", "w");
    fputs("This is an apple.", pFile);
    fseek(pFile, 9, SEEK_SET);
    fputs(" sam", pFile);
    fclose(pFile);
    // remove temp file
    is_eq(remove("./testdata/example.txt"), 0)
}

void test_ftell()
{
    FILE* pFile;
    long size;

    pFile = fopen(test_file, "r");
    is_not_null(pFile);

    fseek(pFile, 0, SEEK_END); // non-portable
    size = ftell(pFile);
    fclose(pFile);

    is_eq(size, filesize);
}

void test_fread()
{
    FILE* pFile;
    int lSize;
    char buffer[1024];
    int result;

    pFile = fopen("tests/getchar.c", "r");
    is_not_null(pFile);

    // obtain file size:
    fseek(pFile, 0, SEEK_END);
    lSize = ftell(pFile);
    is_eq(lSize, 216);

    rewind(pFile);

    // copy the file into the buffer:
    result = fread(buffer, 1, lSize, pFile);
    is_eq(result, lSize);

    // See issue #107
    buffer[lSize - 1] = 0;

    is_eq(strlen(buffer), 215);

    // terminate
    fclose(pFile);
}

void test_fwrite()
{
    FILE* pFile;
    pFile = fopen("./testdata/myfile.bin", "w");
    fwrite("xyz", 1, 3, pFile);
    fclose(pFile);
    // remove temp file
    is_eq(remove("./testdata/myfile.bin"), 0)
}

void test_fgetpos()
{
    FILE* pFile;
    int c;
    int n;
    fpos_t pos;

    pFile = fopen("tests/stdio.c", "r");
    is_not_null(pFile);

    c = fgetc(pFile);
    is_eq(c, '/');

    fgetpos(pFile, &pos);
    for (n = 0; n < 3; n++) {
        fsetpos(pFile, &pos);
        c = fgetc(pFile);
        is_eq(c, '/');
    }

    fclose(pFile);
}

void test_fsetpos()
{
    FILE* pFile;
    fpos_t position;

    pFile = fopen("./testdata/myfile.txt", "w");
    fgetpos(pFile, &position);
    fputs("That is a sample", pFile);
    fsetpos(pFile, &position);
    fputs("This", pFile);
    fclose(pFile);
    // remove temp file
    is_eq(remove("./testdata/myfile.txt"), 0)
}

void test_rewind()
{
    int n;
    FILE* pFile;
    char buffer[27];

    pFile = fopen("./testdata/myfile.txt", "w+");
    for (n = 'A'; n <= 'Z'; n++)
        fputc(n, pFile);
    rewind(pFile);
    fread(buffer, 1, 26, pFile);
    fclose(pFile);
    buffer[26] = '\0';

    is_eq(strlen(buffer), 26);

    // remove temp file
    is_eq(remove("./testdata/myfile.txt"), 0)
}

void test_feof()
{
    FILE* pFile;
    int n = 0;
    pFile = fopen(test_file, "r");
    is_not_null(pFile);

    while (fgetc(pFile) != EOF) {
        ++n;
    }
    if (feof(pFile)) {
        pass("%s", "End-of-File reached.");
        is_eq(n, filesize);
    } else {
        fail("%s", "End-of-File was not reached.");
    }

    fclose(pFile);
}

void test_sprintf()
{
    char buffer[100];
    int cx;
    cx = snprintf(buffer, 100, "The half of %d is %d", 60, 60 / 2);
    is_streq(buffer, "The half of 60 is 30");
    is_eq(cx, 20);
}

void test_snprintf()
{
    char buffer[50];
    int n, a = 5, b = 3;
    n = sprintf(buffer, "%d plus %d is %d", a, b, a + b);
    is_streq(buffer, "5 plus 3 is 8");
    is_eq(n, 13);

    char status[80];
    char* filename = "out.txt";
    int numrows = 10;
    int dirty = 1;
    int len = snprintf(status, sizeof(status), "%.20s - %d lines %s",
        filename, numrows, dirty ? "(modified)" : "");
    printf("%s\n", status);
    is_eq(len, 29);
}

int PrintFError(const char* format, ...)
{
    char buffer[256];
    va_list args;
    va_start(args, format);
    int s = vsprintf(buffer, format, args);
    va_end(args);
    printf("vsnprintf buffer: `%s`\n", buffer);
    return s;
}

void test_vsprintf()
{
    int s = PrintFError("Success function '%s' %.2f", "vsprintf", 3.1415926);
    is_true(s == 19 + 8 + 5);
}

int PrintFError2(const char* format, ...)
{
    char buffer[256];
    va_list args;
    va_start(args, format);
    int s = vsnprintf(buffer, 256, format, args);
    va_end(args);
    printf("vsnprintf buffer: `%s`\n", buffer);
    return s;
}

void test_vsnprintf()
{
    int s = PrintFError2("Success function '%s' %.2f", "vsprintf", 3.1415926);
    is_true(s == 19 + 8 + 5);
    s = PrintFError2("HHELP %d", (int)(2));
    is_true(s == 7);
}

void test_eof()
{
    if ((int)(EOF) == -1) {
        pass("ok");
    }
    char c = EOF;
    if (c == (char)(EOF)) {
        pass("ok");
    }
    char a[1];
    a[0] = 's';
    if (a[0] != EOF) {
        pass("ok");
    }
    a[0] = EOF;
    if (a[0] != EOF) {
        fail("EOF == EOF - fail");
    } else {
        pass("ok");
    }
}

void test_perror()
{
    perror("test perror");
}

void test_getline()
{
    {
        diag("getline: not empty file");
        FILE* pFile;
        pFile = fopen(test_file, "r");
        is_not_null(pFile);

        size_t len;
        char* line = NULL;
        char** pnt = &line;
        size_t* l = &len;
        ssize_t pos = getline(pnt, l, pFile);
        for (int i = 0; i < pos; i++) {
            printf("[%d] : `%d`\n", i, line[i]);
        }
        printf("pos [%d] == filesize [%d]\n", pos, filesize);
        is_eq(pos, filesize);
    }
    {
        diag("getline: not empty file");
        FILE* pFile;
        pFile = fopen("./tests/empty.txt", "r");
        is_not_null(pFile);

        size_t len;
        char* line = NULL;
        char** pnt = &line;
        size_t* l = &len;
        ssize_t pos = getline(pnt, l, pFile);
        is_eq(pos, -1);
    }
}

void test_sscanf()
{
    char sentence[] = "Example\nRudolph is 12 years old";
    char header[50];
    char temp[50];
    char str[20];
    int i;
    sscanf(sentence, "%s %s %s %d", header, str, temp, &i);
    printf("Header: %s\n", header);
    is_eq(i, 12);
    is_streq(str, "Rudolph");
}

void test_FILE()
{
    FILE* p = stdout;
    is_true(p != stderr);
    is_true(p == stdout);
    (void)p;
}

void WriteFormatted(const char* format, ...)
{
    va_list args;
    va_start(args, format);
    vprintf(format, args);
    va_end(args);
}

void test_vprintf()
{
    WriteFormatted("Call with %d variable argument.\n", 1);
    WriteFormatted("Call with %d variable %s.\n", 2, "arguments");
}

void FWriteFormatted(FILE* stream, const char* format, ...)
{
    va_list args;
    va_start(args, format);
    vfprintf(stream, format, args);
    va_end(args);
}

void test_vfprintf()
{
    FILE* pFile;
    pFile = fopen("./testdata/vfprintf.txt", "w");
    FWriteFormatted(pFile, "Call with %d variable argument.\n", 1);
    FWriteFormatted(pFile, "Call with %d variable %s.\n", 2, "arguments");
    fclose(pFile);
}

void test_setbuf()
{
    char buffer[BUFSIZ];
    FILE *pFile1, *pFile2;

    pFile1 = fopen("./testdata/setbuf.txt", "w");
    pFile2 = fopen("./testdata/setbuf2.txt", "a");

    setbuf(pFile1, buffer);
    fputs("This is sent to a buffered stream", pFile1);
    fflush(pFile1);

    setbuf(pFile2, NULL);
    fputs("This is sent to an unbuffered stream", pFile2);

    fclose(pFile1);
    fclose(pFile2);
    (void)(buffer);
}

void test_setvbuf()
{
    FILE* pFile;
    pFile = fopen("./testdata/detvbuf.txt", "w");
    setvbuf(pFile, NULL, _IOFBF, 1024);
    fclose(pFile);
}

int main()
{
    plan(71);

    START_TEST(putchar);
    START_TEST(puts);
    START_TEST(printf);
    START_TEST(remove);
    START_TEST(rename);
    START_TEST(fopen);
    START_TEST(tmpfile);
    START_TEST(tmpnam);
    START_TEST(fclose);
    START_TEST(fflush);
    START_TEST(printf);
    START_TEST(fprintf);
    START_TEST(fscanf);
    START_TEST(fgetc);
    START_TEST(fgets);
    START_TEST(fputc);
    START_TEST(fputs);
    START_TEST(getc);
    START_TEST(putc);
    START_TEST(fseek);
    START_TEST(ftell);
    START_TEST(fread);
    START_TEST(fwrite);
    START_TEST(fgetpos);
    START_TEST(fsetpos);
    START_TEST(rewind);
    START_TEST(feof);
    START_TEST(sprintf);
    START_TEST(snprintf);
    START_TEST(vsprintf);
    START_TEST(vsnprintf);
    START_TEST(eof);
    START_TEST(getline);
    START_TEST(sscanf);
    START_TEST(FILE);
    START_TEST(vprintf);
    START_TEST(vfprintf);
    START_TEST(setbuf);
    START_TEST(setvbuf);

    // that test must be last test
    START_TEST(perror);

    done_testing();
}
