#include "tests.h"
#include <stdio.h>      /* printf */
#include <signal.h>     /* signal, raise, sig_atomic_t */

sig_atomic_t signaled = 0;

void my_handler (int param)
{
  signaled = 1;
}

int main()
{
    plan(0);

	signal (SIGINT, my_handler);
	raise(SIGINT);
	printf ("signaled is %d.\n",signaled);

    done_testing();
}

