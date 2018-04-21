#include "tests.h"
#include <stdio.h>

int d(int v)
{
    return 2 * v;
}

int main()
{
    plan(12);

    int x = 1;

    // Without else
    if (x == 1)
        pass("%s", "x is equal to one");

    if (x == 1) {
        pass("%s", "x is equal to one");
    }

    // With else
    if (x != 1) {
        fail("%s", "x is not equal to one");
    } else {
        pass("%s", "x is equal to one");
    }

    if (NULL) {
        fail("%s", "NULL is zero");
    } else {
        pass("%s", "NULL is not zero");
    }

    if (!NULL) {
        pass("%s", "Invert : ! NULL is zero");
    } else {
        fail("%s", "Invert : ! NULL is not zero");
    }

    diag("Operation inside function");
    int ii = 5;
    if ((ii = d(ii)) != (-1)) {
        is_eq(ii, 10)
    }

    diag("if - else");
    {
        int a = 10;
        int b = 5;
        int c = 0;
        if (a < b) {
        } else if (b > c) {
            pass("ok");
        }
    }

	diag("null");
	{
		char c = 'r';
		char * ch = &c;
		if ( ch != NULL) {
			pass("null test 1")
		}
		if ( NULL != ch) {
			pass("null test 2")
		}
		if ( NULL != NULL) {
			fail("null test 3")
		}
		ch = NULL;
		if ( ch == NULL ) {
			pass("null test 4")
		}
		if ( NULL == ch ) {
			pass("null test 5")
		}
		if ( NULL == NULL ) {
			pass("null test 6")
		}
	}

	/*
	 * TODO strange error for different gcc version
	diag("pointer in if");
	{
		typedef struct rowA{
			unsigned int * p;
		} rowB;
		rowB vv;
		rowB *r = &vv;
		if (r->p == NULL){
			pass("pointer test 1");
		}
		if (!r->p){
			pass("pointer test 2");
		}
		if (r->p){
			fail("pointer 3")
		}
	}
	*/

    done_testing();
}
