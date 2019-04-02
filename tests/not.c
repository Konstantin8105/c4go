#include "tests.h"
#include <stdio.h>

#define START_TEST(t) \
    diag(#t);         \
    test_##t();


#define not_c_type(type)										\
{																\
	diag(#type);												\
	type a;														\
	int p;														\
	diag("not for C type : zero value");						\
	a = 0;														\
	if (!a) { printf("a1"); } else { printf("a2"); }			\
	diag("not-not for C type : zero value");					\
	a = 0;														\
	if (!(!a)) { printf("a3"); } else { printf("a4"); }			\
	diag("for C type : assign");								\
	p = !a;														\
	if (p) { printf("p1"); } else { printf("p2"); }				\
	diag("not for C type : non-zero positive value");			\
	a = 42;														\
	if (!a) { printf("a5"); } else { printf("a6"); }			\
	diag("for C type : assign");								\
	p = !a;														\
	if (p) { printf("p1"); } else { printf("p2"); }				\
	diag("not-not for C type : non-zero positive value");		\
	a = 42;														\
	if (!(!a)) { printf("a7"); } else { printf("a8"); }			\
	diag("not for C type : non-zero negative value");			\
	a = -42;													\
	if (!a) { printf("a5"); } else { printf("a6"); }			\
	diag("not-not for C type : non-zero negative value");		\
	a = -42;													\
	if (!(!a)) { printf("a7"); } else { printf("a8"); }			\
	diag("for C type : assign");								\
	p = !a;														\
	if (p) { printf("p1"); } else { printf("p2"); }				\
	diag("for C type : zero value");							\
	a = 0;														\
	if (a) { printf("a9"); } else { printf("a10"); }			\
	diag("for C type : positive value");						\
	a = 15;														\
	if (a) { printf("a11"); } else { printf("a12"); }			\
	diag("for C type : negative value");						\
	a = -15;													\
	if (a) { printf("a13"); } else { printf("a14"); }			\
	diag("for C type : assign");								\
	p = !a;														\
	if (p) { printf("p1"); } else { printf("p2"); }				\
}

void test_c_types()
{
	not_c_type(char            );
	not_c_type(double          );
	not_c_type(float           );
	not_c_type(int             );
	not_c_type(long double     );
	not_c_type(long long       );
	not_c_type(signed char     );
	not_c_type(unsigned long   );
}

#define not_c_pointer(type)										\
{																\
	diag(#type);												\
	type * a = NULL;											\
	int p;														\
	diag("for C pointer:  null");								\
	if   ( a)  { printf("a1"); } else { printf("a2"); }			\
	diag("for C pointer: not null");							\
	if   (!a)  { printf("a2"); } else { printf("a3"); }			\
	diag("for C pointer: not-not null");						\
	if (!(!a)) { printf("a4"); } else { printf("a5"); }			\
	p = !a;														\
	if (p) { printf("p1"); } else { printf("p2"); }				\
	type b = 42; a = &b ;										\
	diag("for C pointer:  null");								\
	if   ( a)  { printf("a11"); } else { printf("a12"); }		\
	diag("for C pointer: not null");							\
	if   (!a)  { printf("a12"); } else { printf("a13"); }		\
	diag("for C pointer: not-not null");						\
	if (!(!a)) { printf("a14"); } else { printf("a15"); }		\
	p = !a;														\
	if (p) { printf("p1"); } else { printf("p2"); }				\
}

#define not_c_struct(type)										\
{																\
	diag(#type);												\
	type * a = NULL;											\
	int p;														\
	diag("for C pointer:  null");								\
	if   ( a)  { printf("a1"); } else { printf("a2"); }			\
	diag("for C pointer: not null");							\
	if   (!a)  { printf("a2"); } else { printf("a3"); }			\
	diag("for C pointer: not-not null");						\
	if (!(!a)) { printf("a4"); } else { printf("a5"); }			\
	p = !a;														\
	if (p) { printf("p1"); } else { printf("p2"); }				\
	type b ; a = &b ;											\
	diag("for C pointer:  null");								\
	if   ( a)  { printf("a11"); } else { printf("a12"); }		\
	diag("for C pointer: not null");							\
	if   (!a)  { printf("a12"); } else { printf("a13"); }		\
	diag("for C pointer: not-not null");						\
	if (!(!a)) { printf("a14"); } else { printf("a15"); }		\
	p = !a;														\
	if (p) { printf("p1"); } else { printf("p2"); }				\
}

void test_c_pointers()
{
	not_c_pointer(char            );
	not_c_pointer(double          );
	not_c_pointer(float           );
	not_c_pointer(int             );
	not_c_pointer(long double     );
	not_c_pointer(long long       );
	not_c_pointer(signed char     );
	not_c_pointer(unsigned long   );
}

struct str {int i;};
union un{int i;double d;};

void test_c_struct()
{
	not_c_struct(struct str		);
	not_c_struct(union un		);
}

int main()
{
    plan(0);

	START_TEST(c_types);
	START_TEST(c_pointers);
	START_TEST(c_struct);

    done_testing();
}
