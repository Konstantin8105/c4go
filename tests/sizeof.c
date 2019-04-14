// This file contains tests for the sizeof() function and operator.

#include "tests.h"
#include <stdio.h>

#define check_sizes(type, size)         \
    is_eq(sizeof(type), size);          \
    is_eq(sizeof(unsigned type), size); \
    is_eq(sizeof(signed type), size);   \
    is_eq(sizeof(const type), size);    \
    is_eq(sizeof(volatile type), size);

#define FLOAT(type, size) \
    is_eq(sizeof(type), size);

#define OTHER(type, size) \
    is_eq(sizeof(type), size);

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
    int idx;            /* Row index in the file, zero-based. */
    int size;           /* Size of the row, excluding the null term. */
    int rsize;          /* Size of the rendered row. */
    char *chars;        /* Row content. */
    char *render;       /* Row content "rendered" for screen (for TABs). */
    unsigned char *hl;  /* Syntax highlight type for each character in render.*/
    int hl_oc;          /* Row had open comment at end in last syntax highlight
                           check. */
} erow;

typedef struct part1_erow{
	int part;
} part1_erow;
typedef struct part1a_erow{
	int part;
	int part2;
} part1a_erow;

typedef struct part2_erow{
	char * part;
} part2_erow;
typedef struct part2a_erow{
	char * part;
	char * part2;
} part2a_erow;

typedef struct part3_erow{
	unsigned char * part;
} part3_erow;

int main()
{
    plan(72);

    diag("Integer types");
    check_sizes(char, 1);
    check_sizes(short, 2);
    check_sizes(int, 4);
    check_sizes(long, 8);
    check_sizes(long int, 8);
    check_sizes(long long, 8);
    check_sizes(long long int, 8);

    diag("Floating-point types");
    is_eq(sizeof(float), 4);
    is_eq(sizeof(double), 8);
    is_eq(sizeof(long double), 16);

    diag("Other types");
    is_eq(sizeof(void), 1);

    diag("Pointers");
    is_eq(sizeof(char*), 8);
    is_eq(sizeof(char*), 8);
    is_eq(sizeof(short**), 8);
    is_eq(sizeof(long double**), 8);

    diag("Variables");
    a = 123;
    b = 456;
    struct MyStruct s1;
    s1.b = 0;
    union MyUnion u1;
    u1.b = 0;

    is_eq(sizeof(a), 2);
    is_eq(sizeof(b), 4);
    is_eq(sizeof(s1), 40);
    is_eq(sizeof(u1), 16);

    diag("Structures");
    is_eq(sizeof(struct MyStruct), 40);
    is_eq(sizeof(struct MyStruct*), 8);

    diag("Unions");
    is_eq(sizeof(union MyUnion), 16);
    is_eq(sizeof(union MyUnion*), 8);

    diag("Function pointers");
    is_eq(sizeof(main), 1);

    diag("Arrays");
    char c[3] = { 'a', 'b', 'c' };
    c[0] = 'a';
    is_eq(sizeof(c), 3);

    int* d[3];
    d[0] = &b;
    is_eq(sizeof(d), 24);

    int** e[4];
    e[0] = d;
    is_eq(sizeof(e), 32);

    const char* const f[] = { "a", "b", "c", "d", "e", "f" };
    is_eq(sizeof(f), 48);
    is_streq(f[1], "b");

    diag("MyNums");
    is_eq(sizeof(struct MyNums), 104);

    diag("FILE *");
    is_eq(sizeof(FILE*), 8);
    is_eq(sizeof(struct s), 8);

	diag("erow from kilo editor");
	is_eq(sizeof(part1_erow ),  4);
	is_eq(sizeof(part1_erow*),  8);
	is_eq(sizeof(part1a_erow ),  8);
	is_eq(sizeof(part1a_erow*),  8);
	is_eq(sizeof(part2_erow ),  8);
	is_eq(sizeof(part2_erow*),  8);
	is_eq(sizeof(part2a_erow ),  16);
	is_eq(sizeof(part2a_erow*),  8);
	is_eq(sizeof(part3_erow ),  8);
	is_eq(sizeof(part3_erow*),  8);
	is_eq(sizeof(erow ), 32);
	is_eq(sizeof(erow*),  8);

    done_testing();
}
