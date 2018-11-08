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
    diag("ptrdiff_t");
    int numbers[100];
    int *p1 = &numbers[18], *p2 = &numbers[19];
    ptrdiff_t diff = p2 - p1;
    is_eq(diff, 1);
}

int main()
{
    plan(6);

    test_offset();
    test_ptrdiff_t();

    done_testing();
}
