#define MaxMinTest(type) \
{ \
	diag(#type);       \
	type a = 54;       \
	type b = -4;       \
	type c;            \ 
	c = a > b ? a : b ;\
	c = a < b ? a : b; \
	c = b < a ? a : b ;\
	c = b > a ? a : b; \
	                   \
	c = a > b ? a : b + 1 ;\
	c = a < b ? a + 1 : b; \
	c = b < a ? a : b + 1 ;\
	c = b > a ? a + 1 : b; \
} 

void test_min_max()
{
	MaxMinTest(char);
	MaxMinTest(short);
	MaxMinTest(int);
	MaxMinTest(float);
	MaxMinTest(double);
	MaxMinTest(long double);
}

