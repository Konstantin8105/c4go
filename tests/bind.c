#include "tests.h"
#include <stdio.h>

#define START_TEST(t) \
    diag(#t);         \
    test_##t();


// outside C header
#include "bind/bind.h"

void test_int()
{
	int i = 12;
	is_eq(i,12);
	i = bind_int(i);
	is_eq(i,42);
	bind_int_pnt(&i);
	is_eq(i,12);
}

int main()
{
    plan(3);

	START_TEST(int);

    done_testing();
}
