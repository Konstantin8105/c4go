#include "tests.h"
#include <stdlib.h>
#include <stdio.h>
#include <errno.h>


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
