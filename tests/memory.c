#include "tests.h"
#include <stdlib.h>

void test_paren_pointer()
{
	diag("paren pointer");
	int aqq[3][3] = { { 5, 6, 7 }, { 50, 60, 70 } , { 500, 600, 700 } };
	int **pz;
	int c;
	pz = aqq;
	c = -1;
	c = (*pz);
	/* is_eq(c , 0); */
	(void)(pz);
	(void)(c);
}

void test_int_pointer()
{
	unsigned char str[3][2] = {{43,46},{12,78},{33,66}};
	unsigned char **pz = str;
	unsigned char c = 56;
	is_not_null(pz);
	is_eq(c,56);
	c = (*(pz));
	/* is_eq(c, 0); */
	(void)(pz);
	(void)(c);
}

int main()
{
    plan(2);

	test_paren_pointer();
	test_int_pointer();

    done_testing();
}
