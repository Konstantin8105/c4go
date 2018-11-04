#include "tests.h"
#include <sys/ioctl.h>
#include <stdio.h>
int main()
{
	plan(0);
	struct winsize sz;
	printf("%d\n",TIOCGWINSZ);
	(void) sz;
	// TODO : ioctl(0, TIOCGWINSZ, &sz);
	// TODO : is_true(sz.ws_col > 0);
	// TODO : is_true(sz.ws_row > 0);
	// TODO : printf("%5d %5d\n",sz.ws_col,sz.ws_row);
	done_testing();
}
