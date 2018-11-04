#include "tests.h"
#include <unistd.h>
#include <termios.h>
#include <fcntl.h>
#include <sys/ioctl.h>
#include <stdio.h>

int main()
{
	plan(0);

	//-------
	struct winsize sz;
	printf("Call:%d\n",TIOCGWINSZ);
	(void) sz;
	int res = ioctl(STDIN_FILENO, TIOCGWINSZ, &sz);
	printf("Res :%d\n",res);
	if (res > 0){
		printf("Desr:%d\n",STDIN_FILENO);
		is_true(sz.ws_col > 0);
		is_true(sz.ws_row > 0);
		printf("%5d %5d\n",sz.ws_col,sz.ws_row);
		fail("Not acceptable");
	}
	//-------

	done_testing();
}
