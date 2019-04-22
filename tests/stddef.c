#include "tests.h"
#include <stddef.h>

struct foo {
    char a;
    char b[10];
    char c;
};

void test_offset()
{
    is_eq((int)offsetof(struct foo, a), 0);
    is_eq((int)offsetof(struct foo, b), 1);
    is_eq((int)offsetof(struct foo, c), 11);
    is_eq((int)offsetof(struct
              foo /* 
					   comments */
              ,
              // single comment
              b),
        1);
    is_eq((int)offsetof(struct foo, c), 11);
}

void test_ptrdiff_t()
{
    {
        diag("ptrdiff_t : int");
        int numbers[100];
        int *p1 = &numbers[18], *p2 = &numbers[29];
		if (p1 == NULL || p2 == NULL){
			fail("NULL fail");
		}
        ptrdiff_t diff = p2 - p1;
        is_eq(diff, 11);
    }
    {
        diag("ptrdiff_t: long long");
        long long numbers[100];
        long long *p1 = &numbers[18], *p2 = &numbers[29];
        ptrdiff_t diff = p2 - p1;
        is_eq(diff, 11);
    }
}

int main()
{
    plan(7);

    test_offset();
    test_ptrdiff_t();

    done_testing();
}
