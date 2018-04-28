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
	is_eq(c , 5);
	(void)(pz);
	(void)(c);
}

/*
void test_int_pointer()
{
	unsigned char str[3][2] = {{43,46},{12,78},{33,66}};
	unsigned char **pz = str;
	unsigned char c = 56;
	is_not_null(pz);
	is_eq(c,56);
	c = (*(pz));
	// is_eq(c, 0);
	(void)(pz);
	(void)(c);
}

void test_uint_to_pointer()
{
	double value = 42;
	double *p;
	p = &value;
	unsigned int loc = p;
	is_true(loc >= 0);

	double * y;
	y = (double *)(loc);
	// is_true(*y == 42);
	(void)(y);

	p = NULL;
	loc = p;
	is_true(loc == 0);
	(void)(p);

	typedef struct row{
		double * t;
	} row;
	row r;
	unsigned poi = r.t;
	is_true(poi >= 0);
	(void)(r);
}
*/

int main()
{
    plan(1);

	test_paren_pointer();
	// test_int_pointer();
	// test_uint_to_pointer();

    done_testing();
}
