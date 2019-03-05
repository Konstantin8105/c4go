#include "tests.h"
#include <stdio.h>

int main()
{
    plan(1);
    // check function `gets` only
    char str[50];
    int i;
    for (i = 0; i < 50; i++)
        str[i] = '\0';
    gets(str);
    is_streq(str, "7");

    done_testing();
}
