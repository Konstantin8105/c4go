// This file contains tests for the sizeof() function and operator.

#include "tests.h"
#include <stdio.h>

#define is_not_less(arg1, arg2) \
    is_true(arg1 >= arg2) // printf("arg1 = %d\n",arg1);

#define check_sizes(type, size)               \
    is_not_less(sizeof(type), size);          \
    is_not_less(sizeof(unsigned type), size); \
    is_not_less(sizeof(signed type), size);   \
    is_not_less(sizeof(const type), size);    \
    is_not_less(sizeof(volatile type), size);

#define FLOAT(type, size) \
    is_not_less(sizeof(type), size);

#define OTHER(type, size) \
    is_not_less(sizeof(type), size);

// We print the variable so that the compiler doesn't complain that the variable
// is unused.
#define VARIABLE(v, p) \
    printf("%s = (%d) %d bytes\n", #v, p, sizeof(v));

struct MyStruct {
    double a, aa, aaa, aaaa;
    char b;
    char c;
};

union MyUnion {
    long double a;
    char b;
    int c;
};

short a;
int b;

struct MyNums {
    char name[100];
    int size;
    int numbers[];
};

struct s {
    FILE* p;
};

typedef struct erow {
    int idx; /* Row index in the file, zero-based. */
    int size; /* Size of the row, excluding the null term. */
    int rsize; /* Size of the rendered row. */
    char* chars; /* Row content. */
    char* render; /* Row content "rendered" for screen (for TABs). */
    unsigned char* hl; /* Syntax highlight type for each character in render.*/
    int hl_oc; /* Row had open comment at end in last syntax highlight
                           check. */
} erow;

typedef struct part1_erow {
    int part;
} part1_erow;
typedef struct part1a_erow {
    int part;
    int part2;
} part1a_erow;
typedef struct part1b_erow {
    int part;
    int part2;
    int part3;
} part1b_erow;
typedef struct part1c_erow {
    int part;
    int part2;
    int part3;
    char* part4;
} part1c_erow;
typedef struct part1d_erow {
    int part;
    int part2;
    int part3;
    char* part4;
    char* part5;
} part1d_erow;
typedef struct part1e_erow {
    int part;
    int part2;
    int part3;
    char* part4;
    char* part5;
    unsigned char* part6;
} part1e_erow;

typedef struct part2_erow {
    char* part;
} part2_erow;
typedef struct part2a_erow {
    char* part;
    char* part2;
} part2a_erow;

typedef struct part3_erow {
    unsigned char* part;
} part3_erow;

struct editorSyntax {
    char** filematch;
    char** keywords;
    char singleline_comment_start[2];
    char multiline_comment_start[3];
    char multiline_comment_end[3];
    int flags;
};

/* C / C++ */
char* C_HL_extensions[] = { ".c", ".cpp", NULL };
char* C_HL_keywords[] = {
    /* A few C / C++ keywords */
    "switch", "if", "while", "for", "break", "continue", "return", "else",
    "struct", "union", "typedef", "static", "enum", "class",
    /* C types */
    "int|", "long|", "double|", "float|", "char|", "unsigned|", "signed|",
    "void|", NULL
};

#define HL_HIGHLIGHT_STRINGS (1 << 0)
#define HL_HIGHLIGHT_NUMBERS (1 << 1)

/* Here we define an array of syntax highlights by extensions, keywords,
 * comments delimiters and flags. */
struct editorSyntax HLDB[] = {
    { /* C / C++ */
        C_HL_extensions,
        C_HL_keywords,
        "//", "/*", "*/",
        HL_HIGHLIGHT_STRINGS | HL_HIGHLIGHT_NUMBERS }
};

#define HLDB_ENTRIES ((sizeof(HLDB)) / (sizeof(HLDB[0])))

typedef struct SP SP;

struct SP {
    union SP_UN {
        int i;
        char l;
    };
    double d;
    char c;
};

typedef struct SP Mem;
typedef SP Mem2;

int main()
{
    plan(83);

    diag("Integer types");
    check_sizes(char, 1);
    check_sizes(short, 2);
    check_sizes(int, 4);
    check_sizes(long, 8);
    check_sizes(long int, 8);
    check_sizes(long long, 8);
    check_sizes(long long int, 8);

    diag("Floating-point types");
    is_not_less(sizeof(float), 4);
    is_not_less(sizeof(double), 8);
    is_not_less(sizeof(long double), 16);

    diag("Other types");
    is_not_less(sizeof(void), 1);

    diag("Pointers");
    is_not_less(sizeof(char*), 8);
    is_not_less(sizeof(char*), 8);
    is_not_less(sizeof(short**), 8);
    is_not_less(sizeof(long double**), 8);

    diag("Variables");
    a = 123;
    b = 456;
    struct MyStruct s1;
    s1.b = 0;
    union MyUnion u1;
    u1.b = 0;

    is_not_less(sizeof(a), 2);
    is_not_less(sizeof(b), 4);
    is_not_less(sizeof(s1), 40);
    is_not_less(sizeof(u1), 16);

    diag("Structures");
    is_not_less(sizeof(struct MyStruct), 40);
    is_not_less(sizeof(struct MyStruct*), 8);

    diag("Unions");
    is_not_less(sizeof(union MyUnion), 16);
    is_not_less(sizeof(union MyUnion*), 8);

    diag("Function pointers");
    is_not_less(sizeof(main), 1);

    diag("Arrays");
    char c[3] = { 'a', 'b', 'c' };
    c[0] = 'a';
    is_not_less(sizeof(c), 3);

    int* d[3];
    d[0] = &b;
    is_not_less(sizeof(d), 24);

    int** e[4];
    e[0] = d;
    is_not_less(sizeof(e), 32);

    const char* const f[] = { "a", "b", "c", "d", "e", "f" };
    is_not_less(sizeof(f), 48);
    is_streq(f[1], "b");

    diag("MyNums");
    is_not_less(sizeof(struct MyNums), 104);

    diag("FILE *");
    is_not_less(sizeof(FILE*), 8);
    is_not_less(sizeof(struct s), 8);

    diag("erow from kilo editor");
    is_not_less(sizeof(part1_erow), 4);

    is_not_less(sizeof(part1a_erow), 8);
    is_not_less(sizeof(part1b_erow), 12);
    is_not_less(sizeof(part1c_erow), 24);
    is_not_less(sizeof(part1d_erow), 32);
    is_not_less(sizeof(part1e_erow), 40);
    is_not_less(sizeof(erow), 48);

    is_not_less(sizeof(part2_erow), 8);
    is_not_less(sizeof(part2a_erow), 16);
    is_not_less(sizeof(part3_erow), 8);

    diag("HLDB");
    is_not_less(sizeof(HLDB), 32);
    is_not_less(sizeof(HLDB[0]), 32);
    is_true(sizeof(HLDB) == sizeof(HLDB[0]));
    is_eq((HLDB_ENTRIES), 1);

    diag("sqlite examples");
    is_not_less(sizeof(union SP_UN), 4);
    Mem m;
    is_not_less(sizeof(m), 16);
    is_not_less(sizeof(m), sizeof(double) + sizeof(char));
    (void)(m);
    Mem2 m2;
    is_not_less(sizeof(m2), 16);
    is_not_less(sizeof(m2), sizeof(double) + sizeof(char));
    (void)(m2);
    SP s;
    is_not_less(sizeof(s), 16);
    is_not_less(sizeof(s), sizeof(double) + sizeof(char));
    (void)(s);
    union SP_UN sp;
    is_not_less(sizeof(sp), 4);
    is_not_less(sizeof(sp), sizeof(int));
    (void)(sp);

    done_testing();
}
