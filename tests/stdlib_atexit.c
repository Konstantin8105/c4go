/* atexit example */
#include "tests.h"
#include <stdlib.h>     /* atexit */

int r_value = 3;

void fnExit1 (void)
{
	r_value += 2;
}

void fnExit2 (void)
{
	r_value *= 5;
}

void check()
{
	is_eq(r_value,17);
}

void done()
{
    int exit_status = 0;                                                           
    if (total_failures > 0)                                                        
    {                                                                              
        diag("FAILED: There was %d failed tests.", total_failures);                
        exit_status = 101;                                                         
    }                                                                              
    if (current_test != total_tests)                                               
    {                                                                              
        diag("FAILED: Expected %d tests, but ran %d.", total_tests, current_test); 
        exit_status = 102;                                                         
    }                                                                              
	exit(exit_status);
}

int main ()
{
  plan(1);
  atexit (done);
  atexit (check);
  atexit (fnExit1);
  atexit (fnExit2);
  return 0;
}
