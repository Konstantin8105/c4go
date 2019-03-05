#include "tests.h"
#include <stdio.h>

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_cast()
{
    int c[] = { (int)'a', (int)'b' };
    is_eq(c[0], 97);

    double d = (double)1;
    is_eq(d, 1.0);
}

void test_castbool()
{
    char i1 = (1 == 1);
    short i2 = (1 == 1);
    int i3 = (1 == 1);
    long i4 = (1 == 0);
    long long i5 = (1 == 0);

    is_eq((i1 == 1) && (i2 == 1) && (i3 == 1) && (i4 == 0) && (i5 == 0), 1);
}

void char_overflow()
{
    {
        char c;
        c = -1;
        unsigned char u = c;
        is_eq(u, 256 - 1);
    }
    {
        char c = -1;
        unsigned char u = c;
        is_eq(u, 256 - 1);
    }
    {
        char c = (-1);
        unsigned char u = c;
        is_eq(u, 256 - 1);
    }
    {
        char c = (((-1)));
        unsigned char u = c;
        is_eq(u, 256 - 1);
    }
}

typedef double* vertex;
void test_vertex()
{
    diag("vertex");

    double a[1];
    a[0] = 42;
    double b[1];
    b[0] = 45;

    double dxoa;
    vertex triorg = (vertex)(a);
    vertex triapex = (vertex)(b);
    dxoa = triorg[0] - triapex[0];

    is_eq(dxoa, -3);
}

static int strlenChar(const char* z)
{
    int n = 0;
    while (*z) {
        if ((0xc0 & *(z++)) != 0x80)
            n++;
    }
    return n;
}

void test_strCh()
{
    char* z = "Hello, c4go\0";
    is_eq(strlenChar(z), 11);
}

void test_bool_to_int()
{
    int d = 3;
    int x = (d > 0) * 2 / 2 + (d < 10) * 10 + (d == 4) / 2 * 2 + 0 * (d == d);
    is_eq(x, 11);
}

// TODO:
void test_unsafe_pnt()
{
    // (otri).orient = (int) ((unsigned long) (ptr) & (unsigned long) 3l);
    // {
    // long pnt;
    // double d = 42.0;
    // double *dd = &d;
    // long pnt2 = (long) (dd);
    // pnt = pnt2;
    // double *ddd = pnt;
    // is_eq(*ddd, 42);
    // (void)pnt;
    // }
    // {
    // int pnt;
    // long l = 42;
    // pnt = (int)&l;
    // long *l_pnt = (long *)(pnt);
    // is_eq(*l_pnt,42);
    // }
    // {
    // int pnt;
    // long l = 42;
    // pnt = (int) (&l);
    // long *l_pnt = (long *)(pnt);
    // is_eq(*l_pnt,42);
    // }
    {
        void* pnt;
        long l = 42;
        long* d = &l;
        pnt = d;
        long l_pnt = 24;
        l_pnt = *((long*)(pnt));
        is_eq(l_pnt, 42);
        (void)pnt;
    }
    // {
    // int pnt;
    // long l = 42;
    // long *d = &l;
    // pnt = (int)(d);
    // long *l_pnt = (long *)(pnt);
    // is_eq(*l_pnt,42);
    // }
}

void test_init()
{
	char ch[4][10] = {"\"", "\n", "\\", "a"};
	is_streq(ch[0], "\"");
	is_streq(ch[1], "\n");
	is_streq(ch[2], "\\");
	is_streq(ch[3], "a");
}

int main()
{
    plan(42);

    START_TEST(unsafe_pnt);
    START_TEST(bool_to_int);
    START_TEST(cast);
    START_TEST(castbool);
    START_TEST(vertex);
    START_TEST(strCh);
    START_TEST(init);

    {
        typedef unsigned int u32;
        u32 x = 42;
        is_eq(x, 42);
        u32 a[10];
        a[0] = x;
        is_eq(a[0], 42);
    }

    {
        typedef double u32d;
        u32d x = 42.0;
        is_eq(x, 42.0);
        u32d a[10];
        a[0] = x;
        is_eq(a[0], 42.0);
    }

    {
        typedef int integer;
        typedef int INTEGER;
        integer i = 123;
        INTEGER j = 567;
        j = i;
        i = j;
        is_eq(i, 123);
        is_eq(j, 123);
    }

    double* d = (double*)0;
    is_true(d == NULL);
    int* i = (int*)0;
    is_true(i == NULL);
    float* f = (float*)0;
    is_true(f == NULL);
    char* c = (char*)0;
    is_true(c == NULL);

    double* d2 = 0;
    is_true(d2 == NULL);
    int* i2 = 0;
    is_true(i2 == NULL);
    float* f2 = 0;
    is_true(f2 == NULL);
    char* c2 = 0;
    is_true(c2 == NULL);

    diag("Calloc with type");
    {
        double* ddd = (double*)calloc(2, sizeof(double));
        is_not_null(ddd);
        (void)(ddd);
    }
    {
        double* ddd;
        ddd = (double*)calloc(2, sizeof(double));
        is_not_null(ddd);
        (void)(ddd);
    }

    diag("Type convertion from void* to ...");
    {
        void* ptr2;
        int tInt = 55;
        ptr2 = &tInt;
        is_eq(*(int*)ptr2, 55);
        double tDouble = -13;
        ptr2 = &tDouble;
        is_eq(*(double*)ptr2, -13);
        float tFloat = 67;
        is_eq(*(float*)(&tFloat), 67);
    }
    diag("Type convertion from void* to ... in initialization");
    {
        long tLong = 556;
        void* ptr3 = &tLong;
        is_eq(*(long*)ptr3, 556);
    }
    diag("uint16 to bool");
    {
        unsigned short d1 = 2;
        if (!(!(d1))) {
            pass("unsigned short to bool")
        }
        unsigned int d2 = 2;
        if (!(!(d2))) {
            pass("unsigned int to bool")
        }
        unsigned long d3 = 2;
        if (!(!(d3))) {
            pass("unsigned long to bool")
        }
    }
    diag("char to bool");
    {
        char c = 2;
        if (!(!(c))) {
            pass("char to bool")
        }
    }
    diag("long to bool");
    {
        long c = 2;
        if (!(!(c))) {
            pass("long to bool")
        }
    }
	diag("equal slice");
	{
		{
			char n[10] = "hey";
			char m[10] = "boy";
			is_true( n != m );
		}
		{
			char *n = "hey";
			char *m = "boy";
			is_true( n != m );
		}
	}

    char_overflow();

    done_testing();
}
