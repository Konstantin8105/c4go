#include "tests.h"
#include <stdio.h>

void f_empty()
{
    return;
};

int ret_two(int a)
{
    return 2 + a;
}

int ret_s(int a)
{
    return a > 0;
}

int ret_ter(int a)
{
    int (*v)(int);
    v = ret_s;
    return v ? 1 : ret_two(0);
}

void test_return_ternary()
{
    is_eq(ret_ter(10), 1);
    is_eq(ret_ter(-1), 1);
}

#define MaxMinTest(type) \
{ \
	diag(#type);       \
	type a = 54;       \
	type b = -4;       \
	type c;            \ 
	c = a > b ? a : b ;\
	is_true(c == a);   \
	c = a < b ? a : b; \
	is_true(c == b);   \
	c = b < a ? a : b ;\
	is_true(c == a);   \
	c = b > a ? a : b; \
	is_true(c == b);   \
	                   \
	c = a > b ? a : b + 1 ;\
	is_true(c == a);   \
	c = a < b ? a + 1 : b; \
	is_true(c == b);   \
	c = b < a ? a : b + 1 ;\
	is_true(c == a);   \
	c = b > a ? a + 1 : b; \
	is_true(c == b);   \
} 

void test_min_max()
{
	MaxMinTest(int);
	MaxMinTest(float);
	MaxMinTest(double);
}

int main()
{
    plan(41);

    int a = 'a' == 65 ? 10 : 100;
    float b = 10 == 10 ? 1.0 : 2.0;
    char* c = 'x' == 5 ? "one" : "two";
    char d = a == 100 ? 'x' : 1;

    is_eq(a, 100);
    is_eq(b, 1);
    is_streq(c, "two");
    is_eq(d, 'x');

    is_false(0 ? 1 : 0);
    is_false(NULL ? 1 : 0);
    is_true('x' ? 1 : 0);

    a = a == 10 ? b == 1.0 ? 1 : 2 : 2;

    if (a == (a == 2 ? 5 : 10)) {
        fail(__func__);
    } else {
        pass(__func__);
    }

    diag("CStyleCast <ToVoid>");
    {
        double a, b;
        0 ? (void)(a) : (void)(b);
        (void)(a), (void)(b);
    }
    {
        double a;
        0 ? (void)(a) : f_empty();
        (void)(a);
    }
    {
        double b;
        0 ? f_empty() : (void)(b);
        (void)(b);
    }
    {
        ;
        0 ? f_empty() : f_empty();
    }
    pass("Ok - ToVoid");
    {
        diag("oper ++");
        int b = 42;
        int addr = 0;
        b = addr++ ? 1 : 2;
        is_eq(addr, 1);
        is_eq(b, 2);
    }

    test_return_ternary();

    diag("postfix else");
    {
        int a = 12;
        int d;
        d = (1 == 0) ? --a : ++a;
        is_eq(d, 13);
        d = (1 != 0) ? --a : ++a;
        is_eq(d, 12);
    }

    diag("ternaty without middle");
    {
        int a = 42;
        a = 19 > -9 ?: 23;
        is_eq(a, 1);
        a = 54;
        a = 4 < -1 ?: 25;
        is_eq(a, 25);
    }

	diag("max min ternary");
	test_min_max();

    done_testing();
}
