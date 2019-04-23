#include "tests.h"
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

#define test_strto0(actual, func, end) \
    /* strtod */                       \
    func(strtod(actual, &endptr));     \
    func(strtod(actual, NULL));        \
    is_streq(endptr, end);             \
    /* strtof */                       \
    func(strtof(actual, &endptr));     \
    func(strtof(actual, NULL));        \
    is_streq(endptr, end);             \
    /* strtold */                      \
    func(strtold(actual, &endptr));    \
    func(strtold(actual, NULL));       \
    is_streq(endptr, end);

#define test_strto1(actual, func, expected, end) \
    /* strtod */                                 \
    func(strtod(actual, &endptr), expected);     \
    func(strtod(actual, NULL), expected);        \
    is_streq(endptr, end);                       \
    /* strtof */                                 \
    func(strtof(actual, &endptr), expected);     \
    func(strtof(actual, NULL), expected);        \
    is_streq(endptr, end);                       \
    /* strtold */                                \
    func(strtold(actual, &endptr), expected);    \
    func(strtold(actual, NULL), expected);       \
    is_streq(endptr, end);

#define test_ato(actual, expected, end)                 \
    /* atoi */                                          \
    is_eq(atoi(actual), expected);                      \
    /* atol */                                          \
    is_eq(atol(actual), expected);                      \
    /* atoll */                                         \
    is_eq(atoll(actual), expected);                     \
    /* strtol */                                        \
    is_eq(strtol(actual, &endptr, 10), expected);       \
    is_streq(endptr, end);                              \
    is_eq(strtol(actual, NULL, 10), expected);          \
    /* strtoll */                                       \
    is_eq(strtoll(actual, &endptr, 10), expected);      \
    is_streq(endptr, end);                              \
    is_eq(strtoll(actual, NULL, 10), expected);         \
    /* strtoul */                                       \
    if (expected >= 0) {                                \
        is_eq(strtoul(actual, &endptr, 10), expected);  \
        is_streq(endptr, end);                          \
        is_eq(strtoul(actual, NULL, 10), expected);     \
    }                                                   \
    /* strtoull */                                      \
    if (expected >= 0) {                                \
        is_eq(strtoull(actual, &endptr, 10), expected); \
        is_streq(endptr, end);                          \
        is_eq(strtoull(actual, NULL, 10), expected);    \
    }

#define test_strtol(actual, radix, expected, end)      \
    /* strtol */                                       \
    is_eq(strtol(actual, &endptr, radix), expected);   \
    is_streq(endptr, end);                             \
    /* strtoll */                                      \
    is_eq(strtoll(actual, &endptr, radix), expected);  \
    is_streq(endptr, end);                             \
    /* strtoul */                                      \
    is_eq(strtoul(actual, &endptr, radix), expected);  \
    is_streq(endptr, end);                             \
    /* strtoull */                                     \
    is_eq(strtoull(actual, &endptr, radix), expected); \
    is_streq(endptr, end);

void test_malloc1()
{
    diag("malloc1");

    int i = 16, n;
    char* buffer;

    buffer = (char*)malloc(i + 1);
    is_not_null(buffer) or_return();

    for (n = 0; n < i; n++)
        buffer[n] = i % 26 + 'a';
    buffer[i] = '\0';

    is_streq(buffer, "qqqqqqqqqqqqqqqq");
    free(buffer);
}

void test_malloc2()
{
    diag("malloc2");

    int* p;
    p = (int*)malloc(sizeof(int));
    is_not_null(p) or_return();

    *p = 5;

    is_eq(*p, 5);
}

// Mix around all the types to make sure it is still actually allocating the
// correct size.
void test_malloc3()
{
    diag("malloc3");

    is_eq(sizeof(int), 4);
    is_eq(sizeof(double), 8);

    // 10 ints, should be 5 doubles. Also use a bad cast to make sure that it
    // doesn't interfere with the types.
    double* d;
    d = (char*)malloc(sizeof(int) * 10);
    is_not_null(d) or_return();

    // We can't test how much memory has been allocated by Go, but we can test
    // that there are 5 items.
    *d = 123;
    d[4] = 456;

    is_eq(*d, 123);
    is_eq(d[4], 456);
}

void test_malloc4()
{
    diag("malloc4");

    int length = 5;

    char* m = malloc(length * sizeof(char));
    is_not_null(m) or_return();
    (void)(m);

    char* m2 = malloc(sizeof(char) * length);
    is_not_null(m2) or_return();
    (void)(m2);

    char* m3;
    m3 = malloc(sizeof(char) * length);
    is_not_null(m3) or_return();
    (void)(m3);

    char* m4;
    m4 = malloc(length * sizeof(char));
    is_not_null(m4) or_return();
    (void)(m4);
}

void test_malloc5()
{
    int size = 5;
    void* v = malloc(size * size);
    char* c = (char*)v;
    for (int n = 0; n < size * size - 1; n++)
        c[n] = n % 26 + 'a';
    c[size * size - 1] = '\0';
    printf("malloc5: %s\n", c);
}

void test_realloc()
{
    {
        int size = 5;
        void* v = realloc((char*)(NULL), size * size);
        char* c = (char*)v;
        for (int n = 0; n < size * size - 1; n++)
            c[n] = n % 26 + 'a';
        c[size * size - 1] = '\0';
        printf("realloc: %s\n", c);
    }
    {
        // from : https://www.geeksforgeeks.org/g-fact-66/
       int* ptr = (int*)malloc(sizeof(int) * 2);
       int* ptr_new;

       *ptr = 10;
       *(ptr + 1) = 20;

       ptr_new = (int*)realloc(ptr, sizeof(int) * 3);
       *(ptr_new + 2) = 30;

       is_eq(*(ptr_new + 0), 10);
       is_eq(*(ptr_new + 1), 20);
       is_eq(*(ptr_new + 2), 30);
   }
   {
       // from : https://www.tutorialspoint.com/c_standard_library/c_function_realloc.htm
       char* str;

       /* Initial memory allocation */
       str = (char*)malloc(15);
       strcpy(str, "tutorialspoint");
       printf("String = %s\n", str);

       /* Reallocating memory */
       str = (char*)realloc(str, 25);
       strcat(str, ".com");
       printf("String = %s\n", str);
   }
}

// calloc() works exactly the same as malloc() however the memory is zeroed out.
// In Go all allocated memory is zeroed out so they actually are the same thing.
void test_calloc()
{
    diag("calloc");

    is_eq(sizeof(int), 4);
    is_eq(sizeof(double), 8);

    // 10 ints, should be 5 doubles. Also use a bad cast to make sure that it
    // doesn't interfere with the types.
    double* d;
    d = (char*)calloc(10, sizeof(int));
    is_not_null(d) or_return();

    // We can't test how much memory has been allocated by Go, but we can test
    // that there are 5 items.
    *d = 123;
    d[4] = 456;

    is_eq(*d, 123);
    is_eq(d[4], 456);
}

void* cs_calloc(int n)
{
    int sizeT = sizeof(int);
    return (n < 2 ? NULL : calloc(n, sizeT));
    (void)sizeT;
}

void test_calloc2()
{
    diag("calloc2");
    is_null(cs_calloc(0));
    is_true(cs_calloc(5) != NULL);
}

void test_free()
{
    int *buffer1, *buffer2, *buffer3;
    buffer1 = (int*)malloc(100 * sizeof(int));
    buffer2 = (int*)calloc(100, sizeof(int));
    buffer3 = (int*)realloc(buffer2, 500 * sizeof(int));
    int i = 0;
    free((i += 1, buffer1));
    if (buffer2 != NULL) {
        i += 1;
    }
    if (i == 2) {
        free(buffer3);
        i++;
    }
    is_eq(i, 3);
}

int values[] = { 40, 10, 100, 90, 20, 25 };
int compare(const void* a, const void* b)
{
    return (*(int*)a - *(int*)b);
}

void test_system()
{
    diag("system");
    is_eq(system("true"), 0);
    is_not_eq(system("false"), 0);
}

void q_sort()
{
    diag("qsort");
    qsort(values, 6, sizeof(int), compare);
    is_eq(values[0], 10);
    is_eq(values[1], 20);
    is_eq(values[2], 25);
    is_eq(values[3], 40);
    is_eq(values[4], 90);
    is_eq(values[5], 100);
}

#define COMPLEX struct MyComplex
struct MyComplex {
    double re;
    double im;
};
void struct_with_define()
{
    COMPLEX a = { 45, 12 };
    is_eq(a.im, 12);
}

void test_atoi_post()
{
    char* num[3] = { { "123" }, { "987" }, { "456" } };
    char** w = &num;
    int n;
    n = atoi(*w);
    is_streq(*w, "123");
    is_eq(n, 123);

    n = atoi((*w)++);
    is_streq(*w, "23");
    is_eq(n, 123);

    n = atoi((*w));
    is_streq(*w, "23");
    is_eq(n, 23);
}

// TODO:
//	static void * my_memrchr(void *m, int c, long n)
//	{
//		int i;
//		for (i = 0; i < n; i++)
//			if (*(unsigned char *) (m + n - 1 - i) == c)
//				return m + n - 1 - i;
//		// pointer checking: m + n - 1 - i
//		return NULL;
//	}
//
//	void test_my_memrchr()
//	{
//		unsigned char * word = "Hello my dear world";
//		void * p = my_memrchr(word, (int)('d'), 18);
//		is_not_null(p);
//		p = my_memrchr(word, (int)('z'), 18);
//		is_null(p);
//	}

// TODO:
// void printbuffer (const char* pt, size_t max)
// {
  // int length;
  // wchar_t dest;
// 
  // mblen (NULL, 0);         /* reset mblen */
  // mbtowc (NULL, NULL, 0);  /* reset mbtowc */
// 
  // while (max>0) {
    // length = mblen (pt, max);
    // if (length<1) break;
    // mbtowc(&dest,pt,length);
    // printf ("[%lc]",dest);
    // pt+=length; max-=length;
  // }
// }
// 
// void test_mblen()
// {
  // const char str [] = "test string";
  // printbuffer (str,sizeof(str));
// }

// TODO:
// void test_wctomb()
// {
  // const wchar_t str[] = L"wctomb example";
  // const wchar_t* pt;
  // char buffer [MB_CUR_MAX];
  // int i,length;
// 
  // pt = str;
  // while (*pt) {
    // length = wctomb(buffer,*pt);
    // if (length<1) break;
    // for (i=0;i<length;++i) printf ("[%c]",buffer[i]);
    // ++pt;
  // }
// }

int main()
{
    plan(766);

    struct_with_define();

    char* endptr;

    diag("abs");
    is_eq(abs(-5), 5);
    is_eq(abs(7), 7);
    is_eq(abs(0), 0);

    diag("atof");
    is_eq(atof("123"), 123);
    is_eq(atof("1.23"), 1.23);
    is_eq(atof(""), 0);
    is_eq(atof("1.2e6"), 1.2e6);
    is_eq(atof(" \n123"), 123);
    is_eq(atof("\t123foo"), 123);
    is_eq(atof("+1.23"), 1.23);
    is_eq(atof("-1.23"), -1.23);
    is_eq(atof("1.2E-6"), 1.2e-6);
    is_eq(atof("1a2b"), 1);
    is_eq(atof("1a.2b"), 1);
    is_eq(atof("a1.2b"), 0);
    is_eq(atof("1.2Ee-6"), 1.2);
    is_eq(atof("-1..23"), -1);
    is_eq(atof("-1.2.3"), -1.2);
    is_eq(atof("foo"), 0);
    is_eq(atof("+1.2+3"), 1.2);
    is_eq(atof("-1.-23"), -1);
    is_eq(atof("-.23"), -0.23);
    is_eq(atof(".4"), 0.4);
    is_eq(atof("0xabc"), 2748);
    is_eq(atof("0x1b9"), 441);
    is_eq(atof("0x"), 0);
    is_eq(atof("0X1f9"), 505);
    is_eq(atof("-0X1f9"), -505);
    is_eq(atof("+0x1f9"), 505);
    is_eq(atof("0X"), 0);
    is_eq(atof("0xfaz"), 250);
    is_eq(atof("0Xzaf"), 0);
    is_eq(atof("0xabcp2"), 10992);
    is_eq(atof("0xabcP3"), 21984);
    is_eq(atof("0xabcP2z"), 10992);
    is_eq(atof("0xabcp-2"), 687);
    is_eq(atof("0xabcp+2"), 10992);
    is_inf(atof("inf"), 1);
    is_inf(atof("INF"), 1);
    is_inf(atof("Inf"), 1);
    is_inf(atof("-Inf"), -1);
    is_inf(atof("+INF"), 1);
    is_inf(atof("infinity"), 1);
    is_inf(atof("INFINITY"), 1);
    is_inf(atof("Infinity"), 1);
    is_inf(atof("+INFINITY"), 1);
    is_inf(atof("-InfINITY"), -1);
    is_nan(atof("nan"));
    is_nan(atof("NaN"));
    is_nan(atof("+NaN"));
    is_nan(atof("NAN"));
    is_nan(atof("-NAN"));
    is_nan(atof("nanabc123"));
    is_nan(atof("NANz123"));
    is_nan(atof("NaN123z"));
    is_nan(atof("-NANz123"));
    // This causes a segfault in C:
    // is_eq(atof(NULL), 0);

    diag("atoi / atol / atoll / strtol / strtoll / strtoul");
    test_ato("123", 123, "");
    test_ato("+456", 456, "");
    test_ato("-52", -52, "");
    test_ato("foo", 0, "foo");
    test_ato(" \t123", 123, "");
    test_ato("", 0, "");
    test_ato(" \t", 0, " \t");
    test_ato("123abc", 123, "abc");

    diag("div");
    {
        div_t result = div(17, 5);
        is_eq(result.quot, 3)
            is_eq(result.rem, 2)

                result
            = div(-17, 5);
        is_eq(result.quot, -3)
            is_eq(result.rem, -2)

                result
            = div(17, -5);
        is_eq(result.quot, -3)
            is_eq(result.rem, 2)

                result
            = div(-17, -5);
        is_eq(result.quot, 3)
            is_eq(result.rem, -2)
    }

    diag("calloc");
    test_calloc();

    // exit() is handled in tests/exit.c

    // free() is handled with the malloc and calloc tests.
    diag("free");
    test_free();

    diag("getenv");
    is_not_null(getenv("PATH"));
    is_not_null(getenv("HOME"));
    is_null(getenv("FOOBAR"));

    diag("labs");
    is_eq(labs(-5), 5);
    is_eq(labs(7), 7);
    is_eq(labs(0), 0);

    diag("ldiv");
    {
        ldiv_t result = ldiv(17, 5);
        is_eq(result.quot, 3)
            is_eq(result.rem, 2)

                result
            = ldiv(-17, 5);
        is_eq(result.quot, -3)
            is_eq(result.rem, -2)

                result
            = ldiv(17, -5);
        is_eq(result.quot, -3)
            is_eq(result.rem, 2)

                result
            = ldiv(-17, -5);
        is_eq(result.quot, 3)
            is_eq(result.rem, -2)
    }

    diag("llabs");
    is_eq(llabs(-5), 5);
    is_eq(llabs(7), 7);
    is_eq(llabs(0), 0);

    diag("lldiv");
    {
        lldiv_t result = lldiv(17, 5);
        is_eq(result.quot, 3)
            is_eq(result.rem, 2)

                result
            = lldiv(-17, 5);
        is_eq(result.quot, -3)
            is_eq(result.rem, -2)

                result
            = lldiv(17, -5);
        is_eq(result.quot, -3)
            is_eq(result.rem, 2)

                result
            = lldiv(-17, -5);
        is_eq(result.quot, 3)
            is_eq(result.rem, -2)
    }

    diag("malloc");
    test_malloc1();
    test_malloc2();
    test_malloc3();
    test_malloc4();
    test_malloc5();

    diag("realloc");
    test_realloc();

    test_calloc2();

    diag("rand");
    int i, nextRand, lastRand = rand();
    for (i = 0; i < 10; ++i) {
        nextRand = rand();
        is_not_eq(lastRand, nextRand)
    }

    diag("srand");
    srand(0);
    lastRand = rand();
    for (i = 0; i < 10; ++i) {
        nextRand = rand();
        is_not_eq(lastRand, nextRand)
    }

    srand(0);
    int a1 = rand();
    int a2 = rand();
    int a3 = rand();

    srand(0);
    int b1 = rand();
    int b2 = rand();
    int b3 = rand();

    is_eq(a1, b1);
    is_eq(a2, b2);
    is_eq(a3, b3);

    diag("strtod / strtof / strtold");
    test_strto1("123", is_eq, 123, "");
    test_strto1("1.23", is_eq, 1.23, "");
    test_strto1("", is_eq, 0, "");
    test_strto1("1.2e6", is_eq, 1.2e6, "");
    test_strto1(" \n123", is_eq, 123, "");
    test_strto1("\t123foo", is_eq, 123, "foo");
    test_strto1("+1.23", is_eq, 1.23, "");
    test_strto1("-1.23", is_eq, -1.23, "");
    test_strto1("1.2E-6", is_eq, 1.2e-6, "");
    test_strto1("1a2b", is_eq, 1, "a2b");
    test_strto1("1a.2b", is_eq, 1, "a.2b");
    test_strto1("a1.2b", is_eq, 0, "a1.2b");
    test_strto1("1.2Ee-6", is_eq, 1.2, "Ee-6");
    test_strto1("-1..23", is_eq, -1, ".23");
    test_strto1("-1.2.3", is_eq, -1.2, ".3");
    test_strto1("foo", is_eq, 0, "foo");
    test_strto1("+1.2+3", is_eq, 1.2, "+3");
    test_strto1("-1.-23", is_eq, -1, "-23");
    test_strto1("-.23", is_eq, -0.23, "");
    test_strto1(".4", is_eq, 0.4, "");
    test_strto1("0xabc", is_eq, 2748, "");
    test_strto1("0x1b9", is_eq, 441, "");
    test_strto1("0x", is_eq, 0, "x");
    test_strto1("0X1f9", is_eq, 505, "");
    test_strto1("-0X1f9", is_eq, -505, "");
    test_strto1("+0x1f9", is_eq, 505, "");
    test_strto1("0X", is_eq, 0, "X");
    test_strto1("0xfaz", is_eq, 250, "z");
    test_strto1("0Xzaf", is_eq, 0, "Xzaf");
    test_strto1("0xabcp2", is_eq, 10992, "");
    test_strto1("0xabcP3", is_eq, 21984, "");
    test_strto1("0xabcP2z", is_eq, 10992, "z");
    test_strto1("0xabcp-2", is_eq, 687, "");
    test_strto1("0xabcp+2", is_eq, 10992, "");

    test_strto1("inf", is_inf, 1, "");
    test_strto1("INF", is_inf, 1, "");
    test_strto1("Inf", is_inf, 1, "");
    test_strto1("-Inf", is_inf, -1, "");
    test_strto1("+INF", is_inf, 1, "");
    test_strto1("infinity", is_inf, 1, "");
    test_strto1("INFINITY", is_inf, 1, "");
    test_strto1("Infinity", is_inf, 1, "");
    test_strto1("+INFINITY", is_inf, 1, "");
    test_strto1("-InfINITY", is_inf, -1, "");

    test_strto0("nan", is_nan, "");
    test_strto0("NaN", is_nan, "");
    test_strto0("+NaN", is_nan, "");
    test_strto0("NAN", is_nan, "");
    test_strto0("-NAN", is_nan, "");

    test_strto0("nanabc123", is_nan, "abc123");
    test_strto0("NANz123", is_nan, "z123");
    test_strto0("NaN123z", is_nan, "123z");
    test_strto0("-NANz123", is_nan, "z123");

    // This causes a segfault in C:
    //     test_strtod1(NULL, is_eq, 0, "");

    diag("strtol / strtoll / strtoul / strtoull");
    test_strtol("123", 8, 83, "");
    test_strtol("123abc", 16, 1194684, "");
    test_strtol("123abc", 8, 83, "abc");

    diag("abs int");
    {
        int h = -1000;
        printf("%d %d\n", abs(h) / 100, abs(h) % 100);
    }

    diag("atoi_post");
    test_atoi_post();

    diag("system");
    test_system();

    diag("q_sort");
    q_sort();

	// test_mblen();
	// test_wctomb();

    // TODO:
    // diag("my_memrchr");
    // test_my_memrchr();

    done_testing();
}
