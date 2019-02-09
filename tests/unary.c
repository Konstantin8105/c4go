#include "tests.h"
#include <stdio.h>

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_notint()
{
    int i = 0;
    if (!i) {
        pass("good");
    } else {
        fail("fail");
    }

    i = 123;
    if (!i) {
        fail("fail");
    } else {
        pass("good");
    }
}

void test_notptr()
{
    FILE* fp = NULL;
    if (!fp) {
        pass("good");
    } else {
        fail("fail");
    }

    fp = stdin;
    if (!fp) {
        fail("fail");
    } else {
        pass("good");
    }
}

int conv(int i)
{
	return i;
}

void test_sizeoftest()
{
	// example :
	// while ((n = fread(buf, sizeof *buf, sizeof buf, fin)) > 0)
// UnaryExprOrTypeTraitExpr 'unsigned long' sizeof
// `-UnaryOperator 0x2e1a880 'char' lvalue prefix '*'
//   `-ImplicitCastExpr 0x2e1a868 'char *' <ArrayToPointerDecay>
//     `-DeclRefExpr 0x2e1a840 'char [1024]' lvalue Var 0x2e19d78 'buf' 'char [1024]'
	char *buf = NULL;
	if (conv(sizeof *buf) > 0) {
		pass("ok");
	} else {
		fail("fail");
	}
	(void)(buf);
}

int main()
{
    plan(5);

    START_TEST(notint)
    START_TEST(notptr)
	START_TEST(sizeoftest)

    done_testing();
}
