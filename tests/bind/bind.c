#include "test.h"

int before();

int main(void) 
{
	int q = before();
	view0();
	view1();
	view2();
	return q;
}

int before()
{
	int i = 0;
	int *q = &i;
	return *q;
}
