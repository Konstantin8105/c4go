#include "tests.h"
#include <stdio.h>

// TODO: More comprehensive operator tests
// https://github.com/Konstantin8105/c4go/issues/143

void empty() { ; }

int sAdd(char* opt)
{
    int l = strlen(opt) + 12;
    return l;
}

int sMul(char* opt)
{
    int l = strlen(opt) * 12;
    return l;
}

int sMin(char* opt)
{
    int l = strlen(opt) - 12;
    return l;
}

int sDiv(char* opt)
{
    int l = strlen(opt) / 12;
    return l;
}

int simple_repeat(int a)
{
    return a;
}

double* return_null()
{
    return NULL;
}

int f_sizeof(int i)
{
    return i;
}

int reteg(int a)
{
    int arr[5];
    for (int i = 0; i < 5; i++) {
        arr[i] = i;
    }
    int* ptr;
    ptr = &arr[1];
    (void)(ptr);
    return *ptr = a + 1;
}

void view(int c)
{
    printf("%d\n", c);
}

enum Bool { false = 0, true = 1 };
typedef enum Bool bool;
static bool valGlobBool     = true;
static int valGlobInt       = 42;
static double valGlobDouble = 45.0;
static bool restricted      = true;

typedef unsigned int u00;
struct UUU000{
	u00 u[2];
};

int main()
{
    plan(160);

	is_eq(valGlobInt, 42);
	is_eq(valGlobDouble, 45);
	is_eq(valGlobBool, 1);
	is_eq(restricted, 1);

    int i = 10;
    signed char j = 1;
    float f = 3.14159f;
    double d = 0.0;
    char c = 'A';

    i %= 10;
    is_eq(i, 0);

    i += 10;
    is_eq(i, 10);

    i -= 2;
    is_eq(i, 8);

    i *= 2;
    is_eq(i, 16);

    i /= 4;
    is_eq(i, 4);

    i <<= 2;
    is_eq(i, 16);

    i >>= 2;
    is_eq(i, 4);

    i ^= 0xCFCF;
    is_eq(i, 53195);

    i |= 0xFFFF;
    is_eq(i, 65535);

    i &= 0x0000;
    is_eq(i, 0);

    diag("Other types");

    f += 1.0f;
    is_eq(f, 4.14159);

    d += 1.25f;
    is_eq(d, 1.25);

    i -= 255l;
    is_eq(i, -255);

    i += 'A';
    is_eq(i, -190);

    c += 11;
    is_eq(c, 76);

    diag("Shift with signed int");

    i = 4 << j;
    is_eq(i, 8);

    i = 8 >> j;
    is_eq(i, 4);

    i <<= j;
    is_eq(i, 8);

    i >>= j;
    is_eq(i, 4);

    diag("Operator equal for 1 variables");
    int x;
    x = 42;
    is_eq(x, 42);

    diag("Operator equal for 2 variables");
    int y;
    x = y = 1;
    is_eq(x, 1);
    is_eq(y, 1);

    diag("Operator comma in initialization");
    int x2 = 0, y2 = 1;
    is_eq(x2, 0);
    is_eq(y2, 1);

    diag("Operator equal for 3 variables");
    int a, b, c2;
    a = b = c2 = 3;
    is_eq(a, 3);
    is_eq(b, 3);
    is_eq(c2, 3);

    diag("Huge comma problem for Equal operator");
    int q, w, e;
    q = 7, w = q + 3, e = q + w;
    is_eq(q, 7);
    is_eq(w, 10);
    is_eq(e, 17);

    diag("Huge comma problem for Equal operator with Multiplication");
    float qF, wF, eF;
    qF = 7., wF = qF * 3., eF = qF * wF;
    float expectedQ = 7.;
    float expectedW = 7. * 3.;
    float expectedE = 7. * (7. * 3.);
    is_eq(qF, expectedQ);
    is_eq(wF, expectedW);
    is_eq(eF, expectedE);

    diag("Statement expressions");
    int s1 = ({ 2; });
    is_eq(s1, 2);
    is_eq(({ int foo = s1 * 3; foo + 1; }), 7);

    diag("Not allowable var name for Go");
    int type = 42;
    is_eq(type, 42);

    diag("Go keywords inside C code");
    {
        int chan = 42;
        is_eq(chan, 42);
    }
    {
        int defer = 42;
        is_eq(defer, 42);
    }
    {
        int fallthrough = 42;
        is_eq(fallthrough, 42);
    }
    {
        int func = 42;
        is_eq(func, 42);
    }
    {
        int go = 42;
        is_eq(go, 42);
    }
    {
        int import = 42;
        is_eq(import, 42);
    }
    {
        int interface = 42;
        is_eq(interface, 42);
    }
    {
        int map = 42;
        is_eq(map, 42);
    }
    {
        int package = 42;
        is_eq(package, 42);
    }
    {
        int range = 42;
        is_eq(range, 42);
    }
    {
        int select = 42;
        is_eq(select, 42);
    }
    {
        int type = 42;
        is_eq(type, 42);
    }
    {
        int var = 42;
        is_eq(var, 42);
    }
    {
        int _ = 42;
        is_eq(_, 42);
    }

    // checking is_eq is no need, because if "(void)(az)" not transpile,
    // then go build return fail - value is not used
    diag("CStyleCast <ToVoid>");
    {
        char** az;
        (void)(az);
    }
    {
        double* const* az;
        (void)(az);
    }
    {
        int** az;
        (void)(az);
    }
    {
        float* volatile* az;
        (void)(az);
    }

    diag("CStyleCast <ToVoid> with comma");
    {
        unsigned int* ui;
        (void)(empty(), ui);
    }
    {
        long int* li;
        int counter_li = 0;
        (void)(counter_li++, empty(), li);
        is_eq(counter_li, 1);
    }

    diag("switch with initialization");
    switch (0) {
        int ii;
    case 0: {
        ii = 42;
        is_eq(ii, 42);
    }
    case 1: {
        ii = 50;
        is_eq(ii, 50);
    }
    }
    switch (1) {
        int ia;
    case 0: {
        ia = 42;
        is_eq(ia, 42);
    }
    case 1: {
        ia = 60;
        is_eq(ia, 60);
    }
    }

    diag("Binary operators for definition function");
    is_eq(sAdd("rrr"), 15);
    is_eq(sMul("rrr"), 36);
    is_eq(sMin("rrrrrrrrrrrrr"), 1);
    is_eq(sDiv("rrrrrrrrrrrr"), 1);

    diag("Operators +=, -=, *= , /= ... inside []");
    {
        int a[3];
        a[0] = 5;
        a[1] = 9;
        a[2] = -13;
        int iterator = 0;
        is_eq(a[iterator++], 5);
        is_eq(a[iterator], 9);
        is_eq(a[++iterator], -13);
        is_eq(a[iterator -= 2], 5);
        is_eq(a[iterator += 1], 9);
        is_eq(a[(iterator = 0, iterator)], 5);
        is_eq(simple_repeat((iterator = 42, iterator)), 42);
        is_eq(simple_repeat((iterator = 42, ++iterator, iterator)), 43);
        int b = 0;
        for (iterator = 0; b++, iterator < 2; iterator++, iterator--, iterator++) {
            pass("iterator in for");
        }
        is_eq(b, 3);
        iterator = 0;
        if (i++ > 0) {
            pass("i++ > 0 is pass");
        }
    }
    diag("Equals a=b=c=...");
    {
        int a, b, c, d;
        a = b = c = d = 42;
        is_eq(a, 42);
        is_eq(d, 42);
    }
    {
        double a, b, c, d;
        a = b = c = d = 42;
        is_eq(a, 42);
        is_eq(d, 42);
    }
    {
        int a, b, c, d = a = b = c = 42;
        is_eq(a, 42);
        is_eq(d, 42);
    }
    {
        double a, b, c, d = a = b = c = 42;
        is_eq(a, 42);
        is_eq(d, 42);
    }
    {
        double a[3];
        a[0] = a[1] = a[2] = -13;
        is_eq(a[0], -13);
        is_eq(a[2], -13);
    }
    {
        double a[3];
        a[0] = a[1] = a[2] = -13;
        double b[3];
        b[0] = b[1] = b[2] = 5;

        b[0] = a[0] = 42;
        is_eq(a[0], 42);
        is_eq(b[0], 42);
    }
    {
        double v1 = 12;
        int v2 = -6;
        double* b = &v1;
        int* a = &v2;
        *b = *a = 42;
        is_eq(*a, 42);
        is_eq(*b, 42);
    }
    {
        int yy = 0;
        if ((yy = simple_repeat(42)) > 3) {
            pass("ok")
        }
    }
    diag("pointer in IF");
    double* cd;
    if ((cd = return_null()) == NULL) {
        pass("ok");
    }
    (void)(cd);

    diag("increment for char");
    {
        char N = 'g';
        int aaa = 0;
        if ((aaa++, N--, aaa += 3, N) == 102) {
            pass("ok");
        }
        (void)(aaa);
    }
    diag("Comma with operations");
    {
        int x, y, z;
        x = y = z = 1;
        x <<= y <<= z <<= 1;
        is_eq(x, 16);
        is_eq(y, 4);
        is_eq(z, 2);
    }
    {
        int x, y, z;
        x = y = z = 1000;
        x /= y /= z /= 2;
        is_eq(x, 500);
        is_eq(y, 2);
        is_eq(z, 500);
    }
    {
        int x, y, z;
        x = y = z = 3;
        x *= y *= z *= 2;
        is_eq(x, 54);
        is_eq(y, 18);
        is_eq(z, 6);
    }
    diag("char + bool");
    {
        char prefix = 'W';
        char* buf = "text";
        char* v;
        v = buf + (prefix != 0);
        is_not_null(v);
        is_streq(v, "ext");
    }

    diag("Bitwise complement operator ~");
    {
        int i = 35;
        int o = ~(i);
        is_eq(o, -36);
        is_eq(~ - 12, 11);
    }

    diag("summ of bools");
    {
        int u = 0;
        is_true(u == 0);
        u += (1 != 0);
        is_true(u == 1);
    }

    diag("summ of sizeof");
    {
        int x = sizeof(char);
        is_true(x == 1);
        x = x + sizeof(char);
        is_true(x == 2);
        x += sizeof(char) + sizeof(char);
        is_true(x == 4);
        x = sizeof(char) * 5 + sizeof(char);
        is_true(x == 6);
        x = f_sizeof(sizeof(int));
        printf("%d\n", x);
        int y[2];
        y[0] = 2;
        is_true(y[0] == 2);
        is_true(y[sizeof(char) - 1] == 2);
        y[1] = 5;
        is_true(y[1] == 5);
        is_true(y[sizeof(char)] == 5);
    }
    diag("function with equal in return");
    {
        int a = 42;
        a = reteg(a);
        is_eq(a, 43);
    }
    diag("equal in function");
    {
        int a[2];
        a[0] = -1;
        a[1] = 42;
        int b = a[0];
        b += reteg((*a)++);
        is_eq(a[1], 42);
    }
    diag("operation Not in if");
    {
        int addr = 0;
        if (!addr++) {
            is_eq(addr, 1);
        }
    }
    diag("compare char pointer");
    {
        char* b = "happy new code";
		is_true(&b[3] >  &b[0]);
		is_true(&b[3] == &b[3]);
		is_true(&b[3] <  &b[4]);
    }
    diag("kilo.c example");
    {
        unsigned int flag = 100;
        flag &= ~(2 | 256 | 1024);
        is_eq(flag, 100);
    }
    diag("unary + - ");
    {
        int c = 90;
        view(+c);
        view(-c);
        is_eq(c, 90);
    }
    diag("operation |= for enum");
    {
        enum Sflags {
            SGG = 0x01,
            SGP = 0x02,
            SGR = 0x04,
            SGF = 0x08
        } sflags
            = 0;
        is_eq(sflags, 0);
        sflags |= SGG;
        is_eq(sflags, 1);
        sflags |= SGR;
        is_eq(sflags, 5);
    }
	diag("equal paren");
	{
		int a,b;
		a = b = 42;
		is_eq(a,42);
		is_eq(b,42);
		a = (b = 42);
		is_eq(a,42);
		is_eq(b,42);
	}
	diag("equal paren pointer");
	{
		int * a;
		int * b;
		int val = 45;
		a = b = &val;
		is_eq(*a,val);
		is_eq(*b,val);
		val = 42;
		a = (b = &val);
		is_eq(*a,val);
		is_eq(*b,val);
	}
	diag("equal paren u00");
	{
		u00 a,b;
		a = b = 44;
		is_eq(a,44);
		is_eq(b,44);
		a = (b = 42);
		is_eq(a,42);
		is_eq(b,42);
	}
	diag("equal paren UUU000");
	{
		u00 a = 150;
		struct UUU000 bs;
		bs.u[0] = 100;
		a = bs.u[0] = 44;
		is_eq(a,44);
		is_eq(bs.u[0],44);
		a = (bs.u[0] = 42);
		is_eq(a,42);
		is_eq(bs.u[0],42);
	}
	diag("equal paren UUU000 pointer");
	{
		u00 a = 150;
		u00 b = 100;
		u00 *pb = &b;
		a = (*pb = 42);
		is_eq(a,42);
		is_eq(b,42);

		// check more complex
		a = (*(pb) = (45));
		is_eq(a,45);
		is_eq(b,45);
	}
	diag("equal paren member");
	{
		struct sd1 {
			int * mark1;
		};
		int a1 = 900;
		struct sd1 c1  ;
		int F = 76;
		c1.mark1 = &F;
		struct sd1 * b1 = &c1;
		a1 = (*(b1->mark1) = (43));
		is_eq(a1,43);
		is_eq(*(b1->mark1),43);
	}
	diag("equal paren member 2");
	{
		struct sd2 {
			int * mark;
		};
		int a = 900;
		struct sd2 c  ;
		int F = 98;
		c.mark = &F;
		struct sd2 * b = &c;
		a = (*(b->mark + 0) = (43));
		is_eq(a,43);
		is_eq(*(b->mark),43);
	}


	diag("+= paren member");
	{
		struct sd10 {
			int * mark10;
		};
		int a10 = 900;
		struct sd10 c10  ;
		int F = 76;
		c10.mark10 = &F;
		struct sd10 * b10 = &c10;
		a10 = (*(b10->mark10) += (43));
		is_eq(a10,43+76);
		is_eq(*(b10->mark10),43+76);
	}
	diag("+= paren member 2");
	{
		struct sd20 {
			int * mark;
		};
		int a = 900;
		struct sd20 c  ;
		int F = 98;
		c.mark = &F;
		struct sd20 * b = &c;
		a = (*(b->mark + 0) += (43));
		is_eq(a,43+98);
		is_eq(*(b->mark),43+98);
	}
	diag("++ conv");
	{
		int a = 90;
		long b = (long)++a;
		is_eq(a,91);
		is_eq(b,91);
	}
	diag("chars compare");
	{
		char *s = "building\x00";
		char *r = "build\x00";
		if ( *r ) {
			pass("r is valid");
		}
		if ( *s ) {
			pass("s is valid");
		}
		char qwe[4];
		if ( *r && s + 2 < s + sizeof(qwe)) {
			pass("expression is valid");
		}
		(void)(qwe);
		if (*r == ' '){
			fail("not valid comparing");
		}
		if (*r == 'b'){
			pass("valid comparing");
		}
		while ( *r && s + 2 < s + 4){
			pass("while pass")
			break;
		}
	}
	diag("pointer with & ");
	{
		double a = 43;
		double * ptr = &a;
		int pos = (int)((unsigned long) (ptr) & (unsigned long) 3l);
		(void) pos;
		(void) ptr;
	}

    done_testing();
}

struct wordStr {
	char * w;
};

static const struct wordStr * wordQImpl(const char * ch, int y) {
	return NULL;
}

static const struct wordStr * ( * const wordQ)(const char *, int) = wordQImpl;


int v2(
  int *db,
  unsigned mTrace,
  int(*xTrace)(unsigned,void*,void*,void*),
  void *pArg 
){
	return 0;
}
