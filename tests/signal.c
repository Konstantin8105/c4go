#include "tests.h"
#include <signal.h> /* signal, raise, sig_atomic_t */
#include <stdio.h> /* printf */

sig_atomic_t signaled = 0;

void my_handler(int param)
{
    signaled = 1;
}

int main()
{
    plan(0);

    signal(SIGINT, my_handler);
    raise(SIGINT);
    printf("signaled is %d.\n", signaled);

    done_testing();
}
