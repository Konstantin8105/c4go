#include "tests.h"
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    plan(1);

    if (__errno_location() == NULL) {
        fail("not ok");
    } else {
        pass("ok");
    }

    done_testing();
}
