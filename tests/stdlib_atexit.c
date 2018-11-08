/* atexit example */
#include "tests.h"
#include <stdlib.h>     /* atexit */
#include <stdio.h>

int r_value = 3;

void fnExit1 (void)
{
	r_value += 2;
}

void fnExit2 (void)
{
	r_value *= 5;
}

void done(void)
{
	printf("%d\n",r_value);
}

int main ()
{
  plan(0);
  atexit (done);
  atexit (fnExit1);
  atexit (fnExit2);
  done_testing();
}
