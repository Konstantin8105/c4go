#include "tests.h"
#include <stddef.h>

struct foo {
  char a;
  char b[10];
  char c;
};

void test_offset(){
	is_eq((int)offsetof(struct foo,a), 0);
	is_eq((int)offsetof(struct foo,b), 1);
	is_eq((int)offsetof(struct foo,c), 11);
	// TODO
	/* is_eq((int)offsetof(struct  */
	/* 			foo #<{(|  */
	/* 				   comments |)}># */
	/* 			, */
	/* 		// single comment	 */
	/* 			b */
	/* 			), 1); is_eq((int)offsetof(struct foo,c),11); */
}

int main()
{
    plan(3);

	test_offset();

    done_testing();
}
