#include "tests.h"
#include <stdio.h>

int d(int v)
{
    return 2 * v;
}

int main()
{
    plan(21);

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
        char* ch = &c;
        if (ch != NULL) {
            pass("null test 1")
        }
        if (NULL != ch) {
            pass("null test 2")
        }
        if (NULL != NULL) {
            fail("null test 3")
        }
        ch = NULL;
        if (ch == NULL) {
            pass("null test 4")
        }
        if (NULL == ch) {
            pass("null test 5")
        }
        if (NULL == NULL) {
            pass("null test 6")
        }
    }

    diag("Bool to int");
    if ((9 != 0) * 2) {
        pass("ok");
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

    diag("equal on if paren");
    {
        int p = 45;
        int y = 12;
        if (!(!(y = p, y))) {
            pass("case 1");
        }
        if (!(!(y = p))) {
            pass("case 2");
        }
    }

	diag("if equal");
	{
		int x = 0;
		int y = 5;
		if ((x = y) == 0){
			fail("if equal");
		}
		if ((x = 0) == 0){
			is_true(x == 0);
		}
		int l[5];
		for (int i=0;i<5;i++){
			l[i] = i;
		}
		is_true(x == 0);
		is_true(y == 5);
		int s = 2;
		if ((l[0] = y-s) == 3){
			is_true(l[0] == 3);
		}
		if ((l[1] = l[4] - s) == 2){
			is_true(l[1] == 2);
		}
	}
	
	int esct = 0;
	int esc  = 5;
	if (esct += esc) {
		pass("ok"); 
	} else { 
		fail("fail");
	}
	(void)(esc) ;
	(void)(esct);

    done_testing();
}
