#include "tests.h"
#include <sys/ioctl.h>
#include <stdio.h>
int main()
{
	plan(2);
	struct winsize sz;
	ioctl(0, TIOCGWINSZ, &sz);
	is_true(sz.ws_col > 0);
	is_true(sz.ws_row > 0);
	printf("%5d %5d\n",sz.ws_col,sz.ws_row);
	done_testing();
}
