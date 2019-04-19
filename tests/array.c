// Array examples

#include "tests.h"
#include <stdlib.h>

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_intarr()
{
    int a[3];
    a[0] = 5;
    a[1] = 9;
    a[2] = -13;

    is_eq(a[0], 5);
    is_eq(a[1], 9);
    is_eq(a[2], -13);
}

void test_doublearr()
{
    double a[2];
    a[0] = 1.2;
    a[1] = 7; // different type

    is_eq(a[0], 1.2);
    is_eq(a[1], 7.0);
}

void test_intarr_init()
{
    int a[] = { 10, 20, 30 };
    is_eq(a[0], 10);
    is_eq(a[1], 20);
    is_eq(a[2], 30);
}

void test_floatarr_init()
{
    float a[] = { 2.2, 3.3, 4.4 };
    is_eq(a[0], 2.2);
    is_eq(a[1], 3.3);
    is_eq(a[2], 4.4);
}

void test_chararr_init()
{
    char a[] = { 97, 98, 99 };
    is_eq(a[0], 'a');
    is_eq(a[1], 'b');
    is_eq(a[2], 'c');
}

void test_chararr_init2()
{
    char a[] = { 'a', 'b', 'c' };
    is_eq(a[0], 'a');
    is_eq(a[1], 'b');
    is_eq(a[2], 'c');
}

void test_exprarr()
{
    int a[] = { 2 ^ 1, 3 & 1, 4 | 1, (5 + 1) / 2 };
    is_eq(a[0], 3);
    is_eq(a[1], 1);
    is_eq(a[2], 5);
    is_eq(a[3], 3);
}

struct s {
    int i;
    char c;
};

void test_structarr()
{
    struct s a[] = { { 1, 'a' }, { 2, 'b' } };
    is_eq(a[0].i, 1);
    is_eq(a[0].c, 'a');
    is_eq(a[1].i, 2);
    is_eq(a[1].c, 'b');

    struct s b[] = { (struct s){ 1, 'a' }, (struct s){ 2, 'b' } };
    is_eq(b[0].i, 1);
    is_eq(b[0].c, 'a');
    is_eq(b[1].i, 2);
    is_eq(b[1].c, 'b');
}

long dummy(char foo[42])
{
    return sizeof(foo);
}

void test_argarr()
{
    char abc[1];
    is_eq(8, dummy(abc));
}

void test_multidim()
{
    int a[2][3] = { { 5, 6, 7 }, { 50, 60, 70 } };
    is_eq(a[1][2], 70);

    // omit array length
    int b[][3][2] = { { { 1, 2 }, { 3, 4 }, { 5, 6 } },
        { { 6, 5 }, { 4, 3 }, { 2, 1 } } };
    is_eq(b[1][1][0], 4);
    // 2 * 3 * 2 * sizeof(int32)
    is_eq(sizeof(b), 48);

    struct s c[2][3] = { { { 1, 'a' }, { 2, 'b' }, { 3, 'c' } }, { { 4, 'd' }, { 5, 'e' }, { 6, 'f' } } };
    is_eq(c[1][1].i, 5);
    is_eq(c[1][1].c, 'e');
    c[1][1] = c[0][0];
    is_eq(c[1][1].i, 1);
    is_eq(c[1][1].c, 'a');
}

void test_ptrarr()
{
    int b = 22;

    int* d[3];
    d[1] = &b;
    is_eq(*(d[1]), 22);

    int** e[4];
    e[0] = d;
    is_eq(*(e[0][1]), 22);
}

void test_stringarr_init()
{
    char* a[] = { "a", "bc", "def" };
    is_streq(a[0], "a");
    is_streq(a[1], "bc");
    is_streq(a[2], "def");
}

void test_partialarr_init()
{
    // Last 2 values are filled with zeros
    double a[4] = { 1.1, 2.2 };
    is_eq(a[2], 0.0);
    is_eq(a[3], 0.0);

    struct s b[3] = { { 97, 'a' } };
    is_eq(b[0].i, 97);
    is_eq(b[2].i, 0);
    is_eq(b[2].c, 0);
}

extern int arrayEx[];
int arrayEx[4] = { 1, 2, 3, 4 };

int ff() { return 3; }

double rep_double(double a)
{
    return a;
}

int rep_int(int a)
{
    return a;
}

void zero(int* a, int* b, int* c)
{
    *a = *b = *c = 0;
}

float* next_pointer(float* v)
{
    long l = 1;
    long p = 2;
    (void)(l);
    (void)(p);
    return p - p + v + l;
}

double* dvector(long nl, long nh)
{
    double* v;
    v = (double*)malloc((size_t)((nh - nl + 1 + 1) * sizeof(double)));
    for (int i = 0; i < nh - nl; i++) {
        *(v + i) = 42.0;
    }
    return v - nl + 1;
}

typedef struct s structs;
void test_pointer_arith_size_t()
{
    size_t size = 1;
    char* left_ptr;
    char arr[3];
    arr[0] = 'a';
    arr[1] = 'b';
    arr[2] = 'c';
    left_ptr = &arr;
    is_eq(*left_ptr, arr[0]);
    left_ptr = left_ptr + size;
    is_eq(*left_ptr, arr[1]);
    left_ptr += size;
    is_eq(*left_ptr, arr[2]);

    // tests for pointer to struct with size > 1
    structs a[] = { { 1, 'a' }, { 2, 'b' }, { 3, 'c' } };
    is_eq(a[0].i, 1);
    is_eq(a[0].c, 'a');
    is_eq(a[1].i, 2);
    is_eq(a[1].c, 'b');
    is_eq(a[2].i, 3);
    is_eq(a[2].c, 'c');
    structs* ps = &a;
    structs* ps2;
    is_eq(ps->i, 1);
    ps2 = ps + size;
    is_eq(ps2->i, 2);
    ps2 += size;
    is_eq(ps2->i, 3);
    ps2 -= size;
    is_eq(ps2->i, 2);
}

void test_pointer_minus_pointer()
{
    {
        diag("char type");
        char* left_ptr;
        char* right_ptr;
        char arr[300];
        left_ptr = &arr[0];
        right_ptr = &arr[200];
        is_eq(right_ptr - left_ptr, 200);
    }
    {
        diag("long long type");
        long long* left_ptr;
        long long* right_ptr;
        long long arr[300];
        left_ptr = &arr[0];
        right_ptr = &arr[200];
        is_eq(right_ptr - left_ptr, 200);
    }
}

typedef unsigned char pcre_uchar;
typedef unsigned char pcre_uint8;
typedef unsigned short pcre_uint16;
typedef unsigned int pcre_uint32;

#define PT_ANY 0 /* Any property - matches all chars */
#define PT_SC 4 /* Script (e.g. Han) */

#define CHAR_B 'b'

void test_array_to_value_int()
{
    int aqq[1][1] = { { 5 } };
    int** pz;
    int c;
    ///////////////
    pz = aqq;
    c = 99999;
    c = *pz;
    is_eq(c, 5.0);
    ///////////////
    (void)pz;
    (void)c;
    ///////////////
}

void test_array_to_value_unsigned()
{
    unsigned aqq[1][1] = { { 5 } };
    unsigned** pz;
    unsigned c;
    ///////////////
    pz = aqq;
    c = 99999;
    c = *pz;
    is_eq(c, 5.0);
    ///////////////
    (void)pz;
    (void)c;
    ///////////////
}

void test_array_to_value_long()
{
    long aqq[1][1] = { { 5 } };
    long** pz;
    long c;
    ///////////////
    pz = aqq;
    c = 99999;
    c = *pz;
    is_eq(c, 5.0);
    ///////////////
    (void)pz;
    (void)c;
    ///////////////
}

void test_array_to_value_char()
{
    long aqq[1][1] = { { 5 } };
    long** pz;
    long c;
    ///////////////
    pz = aqq;
    c = 99;
    c = *pz;
    is_eq(c, 5.0);
    ///////////////
    (void)pz;
    (void)c;
    ///////////////
}

void test_array_to_value_double()
{
    long aqq[1][1] = { { 5 } };
    long** pz;
    long c;
    ///////////////
    pz = aqq;
    c = 99999;
    c = *pz;
    is_eq(c, 5.0);
    ///////////////
    (void)pz;
    (void)c;
    ///////////////
}

#define SIZE 3
void test_size_pointer()
{
    int A[2][SIZE] = { { 10, 20, 30 }, { 40, 50, 60 } };
    int B[4][SIZE] = { { 0, 1, 2 }, { 3, 0, 4 }, { 5, 6, 0 }, { 7, 8, 9 } };
    int(*pnt)[SIZE];
    pnt = A;
    is_eq(pnt[1][2], A[1][2]);
    is_eq(pnt[0][2], A[0][2]);
    pnt = B;
    is_eq(pnt[1][2], B[1][2]);
    is_eq(pnt[0][2], B[0][2]);
}

#define STACK_SIZE 512
#define STACK_PUSH(A) (STACK[SP++] = A)
static unsigned short STACK[STACK_SIZE];
static unsigned int SP = 0;
void test_array_increment()
{
    double a[10];
    for (int i = 0; i < 10; i++)
        a[i] = i;
    is_eq(a[4], 4);
    int pc = 4;
    is_eq(pc, 4);
    is_eq(a[pc++], 4);
    is_eq(pc, 5);
    unsigned short pc2 = 0;
    STACK_PUSH(pc2);
    is_eq(SP, 1);
}

struct MyStruct {
    int number;
    char symbol;
};
void test_struct_init()
{
    struct MyStruct objA = {.symbol = 'A', .number = 100 };
    struct MyStruct objB = {.number = 200 };
    is_eq(objA.symbol, 'A');
    is_eq(objA.number, 100);
    is_eq(objB.number, 200);
}

struct parg_option {
    const char* name;
    int has_arg;
    int* flag;
    int val;
};

static const struct parg_option po_def[] = {
    { "noarg", 1, NULL, 'n' },
    { "optarg", 2, NULL, 'o' },
    { "reqarg", 3, NULL, 'r' },
    { "foo", 4, NULL, 'f' },
    { "foobar", 5, NULL, 'b' },
    { 0, 6, 0, 0 }
};

void test_parg_struct()
{
    is_eq(po_def[0].has_arg, 1);
    is_eq(po_def[5].has_arg, 6);
    is_streq(po_def[1].name, "optarg");
}

int function_array_field(int a)
{
    return a + 1;
}
void test_function_array()
{
    struct fa {
        int (*pf)(int);
    };
    struct fa f[10];
    int i = 0;
    for (i = 0; i < 10; i++) {
        f[i].pf = function_array_field;
    }
    int y = 42;
    for (i = 0; i < 10; i++) {
        y = ((f[i]).pf)(y);
    }
    is_eq(y, 52);
}

void test_string_array()
{
    {
        diag("point 0");
        struct line_t {
            struct line_t* last;
            struct line_t* next;
            int pos;
        };
        struct line_t l1;
        l1.last = NULL;
        l1.next = NULL;
        struct line_t l2;
        l2.last = &l1;
        l2.next = NULL;
        is_true(l2.last == &l1);
    }
    {
        diag("point 1");
        char ch_arr[3][10] = { "spike", "tom", "jerry" };
        printf("%s\n", (*(ch_arr + 0) + 0));
        printf("%s\n", (*(ch_arr + 0) + 1));
        printf("%s\n", (*(ch_arr + 1) + 2));
    }
    // TODO
    // {
    // 	diag("point 2");
    // 	// see https://stackoverflow.com/questions/6812242/defining-and-iterating-through-array-of-strings-in-c
    // 	char *numbers[] = {"One", "Two", "Three", ""}, **n;
    // 	n = numbers;
    // 	while (*n != "") {
    // 	  printf ("%s\n",  *n++);
    // 	}
    // }
    {
        diag("point 3");
        // see https://stackoverflow.com/questions/6812242/defining-and-iterating-through-array-of-strings-in-c
        static const char* strings[] = { "asdf", "asdfasdf", 0 };
        const char** ptr = strings;
        while (*ptr != 0) {
            printf("%s \n", *ptr);
            ++ptr;
        }
    }
    {
        diag("point 4");
        // see https://codereview.stackexchange.com/questions/71119/printing-the-contents-of-a-string-array-using-pointers
        char* names[] = { "John", "Mona", "Lisa", "Frank" };
        for (int i = 0; i < 4; ++i) {
            char* pos = names[i];
            while (*pos != '\0') {
                printf("%c", *(pos++));
            }
            printf("\n");
        }
    }
    {
        diag("point 5");
        const char* names[] = { "John", "Mona", "Lisa", "Frank", NULL };
        for (int i = 0; names[i]; ++i) {
            const char* ch = names[i];
            while (*ch) {
                putchar(*ch++);
            }
            putchar('\n');
        }
    }
    {
        diag("point 6");
        const char* names[] = { "John", "Mona", "Lisa", "Frank", NULL };
        for (const char** pNames = names; *pNames; pNames++) {
            const char* pName = *pNames;
            while (*pName) {
                putchar(*pName++);
            }
            putchar('\n');
        }
    }
    {
        diag("point 7");
        char* names[] = { "John", "Mona", "Lisa", "Frank" };
        int elements = sizeof(names) / sizeof(names[0]);
        for (int i = 0; i < elements; i++) {
            char* p = names[i];
            while (*p)
                putchar(*p++);
            putchar('\n');
        }
    }
    {
        diag("point 8");
        int array[] = { 5, 2, 9, 7, 15 };
        int i = 0;
        array[i]++;
        printf("%d %d\n", i, array[i]);
        array[i]++;
        printf("%d %d\n", i, array[i]);
        array[i++];
        printf("%d %d\n", i, array[i]);
        array[i++];
        printf("%d %d\n", i, array[i]);
    }
}

void test_typedef_pointer()
{
    typedef double* pd;
    double v[2] = { 42., -42. };
    {
        diag("typedef_pointer : 1");
        pd p = &v[0];
        p++;
        is_eq(*p, v[1]);
    }
    {
        diag("typedef_pointer : 2");
        pd p = v;
        p += 1;
        is_eq(*p, v[1]);
    }
    {
        diag("typedef_pointer : 3");
        pd p = &v[1] - 1;
        p = 0 + p + 0 + 1;
        is_eq(*p, v[1]);
    }
    {
        diag("typedef_pointer : 4");
        pd p = 0 + v + 1 + 0 - 1; // v[0]
        p = 0 + p + 0 + 1 - 0 + 1 - 1; // p = p + 1
        is_eq(*p, v[1]);
    }
}

void view_matrix(int** p, int size1, int size2)
{
    for (int i = 0; i < size1; i++) {
        for (int j = 0; j < size2; j++) {
            printf("      p[%d,%d] = %d\n", i, j, p[i][j]);
        }
    }
}

// TODO : it is not Ok for Debug case 
void test_double_array()
{
    // see https://forums.macrumors.com/threads/understanding-double-pointers-in-c.701091/
    int twod[5][5] = { { 2, 4, 6 , 7, 77}, { 8, 10, 12,13,133 }, { 14, 16, 18, 19,199 }, { 20, 22, 24, 25,255 }, {26,28,30,32,322} };
	printf("%d\n", twod[0][0]);
	int * pp;
	pp = twod[0];
	int ** p = &pp;
	(void)(p);
	printf("%d\n", p[0][0]);

    printf(" p is: %d\n", **p);
    printf("*p + 1 is: %d\n", *(*p + 1));
    {
        diag("cases 1:");
        int* pp = *p;
        printf("    1: %d\n", *(pp++));
        printf("    2: %d\n", *(pp++));
        printf("    3: %d\n", *(pp++));
    }
    // TODO : view_matrix(p,4,3);

	p = &pp;
    {
        diag("cases 1a:");
        int* pp = *p;
        printf("    1: %d\n", (*pp)++);
        printf("    2: %d\n", (*pp)++);
        printf("    3: %d\n", (*pp)++);
    }
    // TODO : view_matrix(p,4,3);

	p = &pp;
    {
        diag("cases 2:");
        int** pp = p;
        printf("    1: %d\n", *((*(pp))++));
        printf("    2: %d\n", *((*(pp))++));
        printf("    3: %d\n", *((*(pp))++));
    }
    // TODO : view_matrix(p,4,3);

	p = &pp;
    // {
        // diag("cases 3:");
        // int** pp = p;
        // printf("    1: %d\n", (*pp)[0]);
        // printf("    2: %d\n", (*pp)[1]);
        // printf("    3: %d\n", (*pp)[2]);
    // }
    // TODO : view_matrix(p,4,3);

	p = &pp;
    // {
        // diag("cases 4:");
        // int** pp = p;
        // printf("    1: %d\n", *(*((pp)++)));
        // printf("    2: %d\n", *(*((pp)++)));
        // printf("    3: %d\n", *(*((pp)++)));
    // }
    // TODO : view_matrix(p,4,3);

	p = &pp;
    // {
        // diag("cases 5:");
        // int** pp = p;
        // printf("    1: %d\n", *((*pp)++));
        // printf("    2: %d\n", *((*pp)++));
        // printf("    3: %d\n", *((*pp)++));
    // }
    // TODO : view_matrix(p,4,3);

	p = &pp;
    // {
        // diag("cases 6:");
        // int** pp = p;
        // printf("    1: %d\n", pp[0][0]);
        // printf("    2: %d\n", pp[0][1]);
        // printf("    3: %d\n", pp[0][2]);
    // }
    // TODO : view_matrix(p,4,3);
}

static void trans(char* p)
{
    printf("trans = `%s`\n", p);
}

void test_func_byte()
{
    char* const gameOver = "game over";
    trans(gameOver);
    char* gameOver2 = "game over";
    trans(gameOver2);
}

void test_negative_index()
{
    double ad[5] = { 1., 2., 4., 5., 6.0 };
    is_eq(ad[0], 1.0);
    double* p = ad;
    p += 3;
    is_eq(*p, 5.0);
    is_eq(p[-1], 4.0);
    double* ds = &(p[-1]);
    is_eq(ds[-1], 2.0);
}

void test_matrix_init()
{
    int rows = 2;
    int cols = 3;
    int i, j;
    double** m;

    m = (double**)malloc((unsigned)rows * sizeof(double*));
    for (i = 0; i < rows; i++) {
        m[i] = (double*)malloc((unsigned)cols * sizeof(double));
    }

    for (i = 0; i < rows; i++) {
        for (j = 0; j < cols; j++) {
            printf("init [%d , %d]\n", i, j);
            m[i][j] = i * cols + j;
        }
    }

    for (i = 0; i < rows; i++) {
        for (j = 0; j < cols; j++) {
            is_eq(m[i][j], i * cols + j);
        }
    }
}

struct someR {
    unsigned long* ul;
};

void test_post_pointer()
{
    struct someR R;
    unsigned long ull[6] = { 2, 4, 8, 10, 12, 34 };
    R.ul = ull;
    struct someR* pR = &R;
    for (int i = 0; i < 5; i++) {
        printf("%d\n", (int)(*pR->ul));
        is_eq(ull[i], *pR->ul);
        if (i < 4) {
            pR->ul++;
        }
    }
}

int main()
{
    plan(199);

    test_parg_struct();
    START_TEST(struct_init);
    START_TEST(array_increment);
    START_TEST(array_to_value_int);
    START_TEST(array_to_value_unsigned);
    START_TEST(array_to_value_long);
    START_TEST(array_to_value_char);
    START_TEST(array_to_value_double);
    START_TEST(size_pointer);
    START_TEST(intarr);
    START_TEST(doublearr);
    START_TEST(intarr_init);
    START_TEST(floatarr_init);
    START_TEST(chararr_init);
    START_TEST(chararr_init2);
    START_TEST(exprarr);
    START_TEST(structarr);
    START_TEST(argarr);
    START_TEST(multidim);
    START_TEST(ptrarr);
    START_TEST(stringarr_init);
    START_TEST(partialarr_init);
    START_TEST(function_array);
    START_TEST(typedef_pointer);

    diag("arrayEx");
    is_eq(arrayEx[1], 2.0);

    diag("Array arithmetic");
    float a[5];
    a[0] = 42.;
    is_eq(a[0], 42.);
    a[0 + 1] = 42.;
    is_eq(a[1], 42);
    a[2] = 42.;
    is_eq(a[2], 42);

    diag("Pointer arithmetic. Part 1");
    float* b;
    b = (float*)calloc(5, sizeof(float));

    *b = 42.;
    is_eq(*(b + 0), 42.);

    *(b + 1) = 42.;
    is_eq(*(b + 1), 42.);
    *(2 + b) = 42.;
    is_eq(*(b + 2), 42.);

    *(b + ff()) = 45.;
    is_eq(*(b + 3), 45.);
    *(ff() + b + 1) = 46.;
    is_eq(*(b + 4), 46.);

    *(b + (0 ? 1 : 2)) = -1.;
    is_eq(*(b + 2), -1);

    *(b + 0) = 1;
    *(b + (int)(*(b + 0)) - 1) = 35;
    is_eq(*(b + 0), 35);

    *(b + (int)((float)(2))) = -45;
    is_eq(*(b + 2), -45);

    *(b + 1 + 3 + 1 - 5 * 1 + ff() - 3) = -4.0;
    is_eq(*(b + 0), -4.0);
    is_eq(*b, -4.0);

    is_eq((*(b + 1 + 3 + 1 - 5 * 1 + ff() - 3 + 1) = -48.0, *(b + 1)), -48.0);
    {
        int rrr;
        (void)(rrr);
    }

    diag("Pointer arithmetic. Part 2");
    {
        float* arr;
        arr = (float*)calloc(1 + 1, sizeof(float));
        is_true(arr != NULL);
        (void)(arr);
    }
    {
        float* arr;
        arr = (float*)calloc(1 + ff(), sizeof(float));
        is_true(arr != NULL);
        (void)(arr);
    }
    {
        float* arr;
        arr = (float*)calloc(ff() + ff(), sizeof(float));
        is_true(arr != NULL);
        (void)(arr);
    }
    {
        float* arr;
        arr = (float*)calloc(ff() + 1 + 0 + 0 + 1 * 0, sizeof(float));
        is_true(arr != NULL);
        (void)(arr);
    }

    diag("Pointer to Pointer. 1");
    {
        double Var = 42;
        double** PPptr1;
        double* PPptr2;
        PPptr2 = &Var;
        PPptr1 = &PPptr2;
        is_eq(**PPptr1, Var)
            Var
            = 43;
        is_eq (**PPptr1, Var)(void)(PPptr1);
        (void)(PPptr2);
    }
    diag("Pointer to Pointer. 2");
    {
        double Var = 42.0, **PPptr1, *PPptr2;
        PPptr2 = &Var;
        PPptr1 = &PPptr2;
        is_eq(**PPptr1, Var)
            Var
            = 43.0;
        is_eq (**PPptr1, Var)(void)(PPptr1);
        (void)(PPptr2);
    }
    diag("Pointer to Pointer. 3");
    {
        int i = 50;
        int** ptr1;
        int* ptr2;
        ptr2 = &i;
        ptr1 = &ptr2;
        is_eq(**ptr1, i);
        is_eq(*ptr2, i);
    }
    diag("Pointer to Pointer. 4");
    {
        double arr[5] = { 10., 20., 30., 40., 50. };
        double* ptr;
        ptr = &arr;
        is_eq(*ptr, 10.);
        ++ptr;
        is_eq(*ptr, 20.);
    }
    diag("Pointer to Pointer. 5");
    {
        double arr[5] = { 10., 20., 30., 40., 50. };
        double* ptr;
        ptr = &arr;
        is_eq(*ptr, 10.);
        ptr += 1;
        is_eq(*ptr, 20.);
    }
    diag("Pointer to Pointer. 6");
    {
        int arr[5] = { 10, 20, 30, 40, 50 };
        int* ptr;
        ptr = &arr;
        is_eq(*ptr, 10);
        ptr = 1 + ptr;
        is_eq(*ptr, 20);
    }
    diag("Pointer to Pointer. 7");
    {
        double arr[5] = { 10., 20., 30., 40., 50. };
        double* ptr;
        ptr = &arr;
        is_eq(*ptr, 10.);
        ptr = 1 + ptr;
        is_eq(*ptr, 20.);
    }
    diag("Pointer to Pointer. 8");
    {
        double arr[5] = { 10., 20., 30., 40., 50. };
        double* ptr;
        ptr = &arr;
        is_eq(*ptr, 10.);
        ptr++;
        is_eq(*ptr, 20.);
    }
    diag("Pointer to Pointer. 9");
    {
        double arr[5] = { 10., 20., 30., 40., 50. };
        double* ptr;
        ptr = &arr[2];
        is_eq(*ptr, 30.);
        ptr = ptr - 1;
        is_eq(*ptr, 20.);
    }
    diag("Pointer to Pointer. 10");
    {
        double arr[5] = { 10., 20., 30., 40., 50. };
        double* ptr;
        ptr = &arr[2];
        is_eq(*ptr, 30.);
        ptr -= 1;
        is_eq(*ptr, 20.);
    }
    diag("Pointer to Pointer. 11");
    {
        double arr[5] = { 10., 20., 30., 40., 50. };
        double* ptr;
        ptr = &arr[2];
        is_eq(*ptr, 30.);
        ptr--;
        is_eq(*ptr, 20.);
    }
    diag("Pointer to Pointer. 12");
    {
        double arr[5] = { 10., 20., 30., 40., 50. };
        double* ptr;
        int i = 0;
        for (ptr = &arr[0]; i < 5; ptr++) {
            is_eq(*ptr, arr[i]);
            i++;
        }
    }
    diag("Pointer to Pointer. 13");
    {
        struct temp_str {
            double* qwe;
        };
        struct temp_str t;
        double a[5] = { 10., 20., 30., 40., 50. };
        t.qwe = &a[0];
        double* ptr;
        int i = 0;
        for (ptr = &t.qwe[0]; i < 5; ptr++) {
            is_eq(*ptr, t.qwe[i]);
            i++;
        }
        for (ptr = t.qwe; i < 5; ptr++) {
            is_eq(*ptr, t.qwe[i]);
            i++;
        }
        ptr = 0 + t.qwe + 1;
        is_eq(*ptr, t.qwe[1]);
    }
    diag("Operation += 1 for double array");
    {
        float** m;
        m = (float**)malloc(5 * sizeof(float*));
        is_not_null(m);
        m[0] = (float*)malloc(10 * sizeof(float));
        m[1] = (float*)malloc(10 * sizeof(float));
        m[0] += 1;
        (void)(m);
        pass("ok");
    }
    diag("*Pointer = 0");
    {
        int a, b, c;
        a = b = c = 10;
        is_eq(a, 10);
        zero(&a, &b, &c);
        is_eq(a, 0);
        is_eq(b, 0);
        is_eq(c, 0);
        pass("ok");
    }
    diag("pointer + long");
    {
        float* v = (float*)malloc(5 * sizeof(float));
        *(v + 0) = 5;
        *(v + 1) = 6;
        is_eq(*(next_pointer(v)), 6);
    }
    diag("create array");
    {
        double* arr = dvector(1, 12);
        is_not_null(arr);
        is_eq(arr[1], 42.0);
        is_eq(arr[9], 42.0);
        (void)(arr);
    }

    diag("Increment inside array 1");
    {
        float f[4] = { 1.2, 2.3, 3.4, 4.5 };
        int iter = 0;
        is_eq(f[iter++], 1.2);
        is_eq(f[iter += 1], 3.4);
        is_eq(f[--iter], 2.3);
    }
    diag("Increment inside array 2");
    {
        struct struct_I_A {
            double* arr;
            int* pos;
        };
        struct struct_I_A siia[2];
        {
            double t_arr[5];
            siia[0].arr = t_arr;
        }
        {
            double t_arr[5];
            siia[1].arr = t_arr;
        }
        {
            int t_pos[1];
            siia[0].pos = t_pos;
        }
        {
            int t_pos[1];
            siia[1].pos = t_pos;
        }
        int t = 0;
        int ii, jj;
        int one = 1;

        siia[0].arr[0] = 45.;
        siia[0].arr[1] = 35.;
        siia[0].arr[2] = 25.;

        siia[0].pos[0] = 0;
        ii = -1;
        jj = -1;
        is_eq(siia[0].arr[(t++, siia[jj += one].pos[ii += one] += one, siia[jj].pos[ii])], 35.);

        siia[0].pos[0] = 0;
        ii = -1;
        jj = -1;
        is_eq(siia[0].arr[(t++, siia[++jj].pos[++ii]++, siia[jj].pos[ii])], 35.);

        siia[0].pos[0] = 2;
        ii = -1;
        jj = -1;
        is_eq(siia[0].arr[(t++, siia[0].pos[ii += 1] -= 1, siia[0].pos[ii])], 35.);

        siia[0].pos[0] = 2;
        ii = -1;
        jj = -1;
        is_eq(siia[0].arr[(t++, siia[0].pos[ii += 1]--, siia[0].pos[ii])], 35.);

        is_eq(t, 4);
        (void)(t);
    }
    diag("Increment inside array 3");
    {
        struct struct_I_A3 {
            double* arr;
            int pos;
        };
        struct struct_I_A3 siia[2];
        {
            double t_arr[5];
            siia[0].arr = t_arr;
        }
        {
            double t_arr[5];
            siia[1].arr = t_arr;
        }

        siia[0].arr[0] = 45.;
        siia[0].arr[1] = 35.;
        siia[0].arr[2] = 25.;

        siia[0].pos = 0;
        is_eq(siia[0].arr[siia[0].pos += 1], 35.);

        siia[0].pos = 0;
        is_eq(siia[0].arr[siia[0].pos++], 45.);

        siia[0].pos = 0;
        is_eq(siia[0].arr[++siia[0].pos], 35.);

        siia[0].pos = 2;
        is_eq(siia[0].arr[siia[0].pos -= 1], 35.);

        siia[0].pos = 2;
        is_eq(siia[0].arr[siia[0].pos--], 25.);
    }
    diag("Increment inside array 4");
    {
        struct struct_I_A4 {
            double* arr;
            int pos;
        };
        struct struct_I_A4 siia[2];
        {
            double t_arr[5];
            siia[0].arr = t_arr;
        }
        {
            double t_arr[5];
            siia[1].arr = t_arr;
        }
        int t = 0;

        siia[0].arr[0] = 45.;
        siia[0].arr[1] = 35.;
        siia[0].arr[2] = 25.;

        siia[0].pos = 0;
        is_eq(siia[0].arr[(t++, siia[0].pos += 1)], 35.);

        siia[0].pos = 0;
        is_eq(siia[0].arr[(t++, siia[0].pos++)], 45.);

        siia[0].pos = 2;
        is_eq(siia[0].arr[(t++, siia[0].pos -= 1)], 35.);

        siia[0].pos = 2;
        is_eq(siia[0].arr[(t++, siia[0].pos--)], 25.);

        is_eq(t, 4);
        (void)(t);
    }
    diag("Increment inside array 5");
    {
        struct struct_I_A5 {
            double* arr;
            int pos;
        };
        struct struct_I_A5 siia[2];
        {
            double t_arr[5];
            siia[0].arr = t_arr;
        }
        {
            double t_arr[5];
            siia[1].arr = t_arr;
        }
        int t = 0;

        siia[0].arr[0] = 45.;
        siia[0].arr[1] = 35.;
        siia[0].arr[2] = 25.;

        siia[0].pos = 0;
        is_eq(siia[0].arr[(t++, siia[0].pos += 1, siia[0].pos)], 35.);

        siia[0].pos = 0;
        is_eq(siia[0].arr[(t++, siia[0].pos++, siia[0].pos)], 35.);

        siia[0].pos = 2;
        is_eq(siia[0].arr[(t++, siia[0].pos -= 1, siia[0].pos)], 35.);

        siia[0].pos = 2;
        is_eq(siia[0].arr[(t++, siia[0].pos--, siia[0].pos)], 35.);

        is_eq(t, 4);
        (void)(t);
    }
    diag("Increment inside array 6");
    {
        struct struct_I_A6 {
            double* arr;
            int* pos;
        };
        struct struct_I_A6 siia[2];
        {
            double t_arr[5];
            siia[0].arr = t_arr;
        }
        {
            double t_arr[5];
            siia[1].arr = t_arr;
        }
        {
            int t_pos[1];
            siia[0].pos = t_pos;
        }
        {
            int t_pos[1];
            siia[1].pos = t_pos;
        }
        int t = 0;

        siia[0].arr[0] = 45.;
        siia[0].arr[1] = 35.;
        siia[0].arr[2] = 25.;

        siia[0].pos[0] = 0;
        is_eq(siia[0].arr[(t++, siia[0].pos[0] += 1)], 35.);

        siia[0].pos[0] = 0;
        is_eq(siia[0].arr[(t++, siia[0].pos[0]++)], 45.);

        siia[0].pos[0] = 2;
        is_eq(siia[0].arr[(t++, siia[0].pos[0] -= 1)], 35.);

        siia[0].pos[0] = 2;
        is_eq(siia[0].arr[(t++, siia[0].pos[0]--)], 25.);

        is_eq(t, 4);
        (void)(t);
    }

    test_pointer_arith_size_t();
    START_TEST(pointer_minus_pointer);

    diag("calloc with struct");
    {
        struct cws {
            float* w;
            int nw;
        };
        struct cws t;
        t.nw = 5;
        t.w = (float*)calloc(t.nw, sizeof(*t.w));
        is_not_null(t.w);
        (void)(t);
    }
    {
        diag("[][]char += 1");
        char w1[] = "hello";
        char w2[] = "world";
        char w3[] = "people";
        char* p1 = w1;
        char* p2 = w2;
        char* p3 = w3;
        char* pa[3] = { p1, p2, p3 };
        char** pnt = pa;
        char** pnt2 = pa;
        *pnt += 1;
        is_streq(*pnt, "ello");
        (*pnt2)++;
        is_streq(*pnt2, "llo");
    }
    {
        diag("pnt of value : size_t");
        size_t len = 42;
        size_t* l = &len;
        is_eq(*l, len);
    }
    {
        diag("pnt of value : ssize_t");
        ssize_t len = 42;
        ssize_t* l = &len;
        is_eq(*l, len);
    }
    START_TEST(string_array);
    START_TEST(double_array);
    START_TEST(func_byte);
    START_TEST(negative_index);
    START_TEST(matrix_init);
    START_TEST(post_pointer);

    done_testing();
}
