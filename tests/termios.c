// termios checking

#include "tests.h"
#include <termios.h>

int main()
{
    plan(0);

    struct termios t;

    (void)t;

    done_testing();
};
