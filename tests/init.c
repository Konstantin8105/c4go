// Initialization structs, arrays, ...

#include "tests.h"
#include <stdio.h>

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_array_float()
{
    float a[] = { 2.2, 3.3, 4.4 };
    is_eq(a[0], 2.2);
    is_eq(a[1], 3.3);
    is_eq(a[2], 4.4);
}

void test_array_char()
{
    char a[] = { 'a', 'b', 'c' };
    is_eq(a[0], 'a');
    is_eq(a[1], 'b');
    is_eq(a[2], 'c');
}

void test_struct_init()
{
	struct s {
		int i;
		char c;
	};

    struct s a[] = { { 1, 'a' }, { 2, 'b' } };
    is_eq(a[0].i, 1);
    is_eq(a[0].c, 'a');
    is_eq(a[1].i, 2);
    is_eq(a[1].c, 'b');

    struct s c[2][3] = { { { 1, 'a' }, { 2, 'b' }, { 3, 'c' } }, { { 4, 'd' }, { 5, 'e' }, { 6, 'f' } } };
    is_eq(c[1][1].i, 5);
    is_eq(c[1][1].c, 'e');
}

void test_matrix_double()
{
    int a[2][3] = { { 5, 6, 7 }, { 50, 60, 70 } };
    is_eq(a[1][2], 70);
}

static char *kmap_fa[256] = {
	[ 0 ] = "fa",
	['`'] = "@",
	['1'] = "abc",
};
void test_equals_chars()
{
	is_streq(kmap_fa[ 0 ], "fa" );
	is_streq(kmap_fa['`'], "@"  );
	is_streq(kmap_fa['1'], "abc");
	for (int i=0;i < 256;i++) {
			kmap_fa[i] = "Y";
	}
}

static char *di[][2] = {
	{"cq", ";"},
	{"pl", "+"},
	{"hy", "-"},
	{"sl", "/"},
};
void test_di()
{
	is_streq(di[0][0], "cq");
	is_streq(di[1][0], "pl");
	is_streq(di[2][0], "hy");
	is_streq(di[3][0], "sl");
	is_streq(di[0][1], ";");
	is_streq(di[1][1], "+");
	is_streq(di[2][1], "-");
	is_streq(di[3][1], "/");
}


int ec_print(char *s){
	printf("%s\n",s);
	return 32;
}
int ec_insert(char *s){
	printf("%s\n",s);
	return 42;
}
static struct excmd {
	char *abbr;
	char *name;
	int (*ec)(char *s);
} excmds[] = {
	{"p", "print", ec_print},
	{"a", "append", ec_insert},
	{},
};
void test_ex()
{
	is_streq(excmds[0].abbr, "p")
	is_streq(excmds[0].name, "print")
	is_eq(excmds[0].ec("0"), 32);
	
	is_streq(excmds[1].abbr, "a")
	is_streq(excmds[1].name, "append")
	is_eq(excmds[1].ec("1"), 42);

	is_true(excmds[2].abbr == NULL);
	is_true(excmds[2].name == NULL);
	is_true(excmds[2].ec   == NULL);
}

#define test_part_array_pointer(type)		\
{											\
	diag("test_part_array_pointer");		\
    diag(#type);							\
	type a0 = 1;							\
	type a1 = 3;							\
	type aZ = 4;							\
	type * d[5] = { &a0 , &a1 };			\
	is_eq(*d[0], a0);						\
	is_eq(*d[1], a1);						\
	is_true(d[4] == NULL);					\
	for (int i = 0; i < 5; i++) {			\
		d[i] = &aZ;							\
	}										\
	is_eq(*d[0], aZ);						\
	is_eq(*d[1], aZ);						\
}

#define test_part_array(type)		\
{									\
    diag("test_part_array");		\
	diag(#type); 					\
	type d[5] = { 1 , 3 };			\
	is_eq(d[0], 1);					\
	is_eq(d[1], 3);					\
	for (int i = 0; i < 5; i++) {	\
		d[i] = (float) i;			\
	}								\
	is_eq(d[0], 0);					\
	is_eq(d[1], 1);					\
}

void test_partly()
{
	// Test partly initialization of array
	test_part_array(char            );
	test_part_array(double          );
	test_part_array(float           );
	test_part_array(int             );
	test_part_array(long double     );
	test_part_array(long long       );
	test_part_array(signed char     );
	test_part_array(unsigned long   );

	// Test partly initialization of array pointer
	test_part_array_pointer(char               );
	test_part_array_pointer(double             );
	test_part_array_pointer(float              );
	test_part_array_pointer(int                );
	test_part_array_pointer(long double        );
	test_part_array_pointer(long long          );
	test_part_array_pointer(signed char        );
	test_part_array_pointer(unsigned long      );
}

void test_FILE()
{
	FILE * F[3] = {stderr};
	is_true(F[0] == stderr);
	for (int i = 0; i < 3;i++) {
		F[i] = stdout;	
	}
	is_true(F[1] == stdout);
}

void test_void()
{
	void * v[3] = {stderr};
	is_true((FILE *)(v[0]) == stderr);
	double r = 42;
	for (int i = 0; i < 3;i++) {
		v[i] = &r;	
	}
	is_eq(*(double *)(v[1]) , r);
}

static int xai, xaw;
static struct option {
	char *abbr;
	char *name;
	int  *vars;
} options[] = {
	{"ai", "autoindent", &xai},
	{"aw", "autowrite", &xaw},
	{},
};
void test_options()
{
	is_streq(options[0].abbr, "ai");
	is_streq(options[0].name, "autoindent");
	is_not_null(options[0].vars);

	is_streq(options[1].abbr, "aw");
	is_streq(options[1].name, "autowrite");
	is_not_null(options[1].vars);

	is_true(options[2].abbr == NULL);
	is_true(options[2].name == NULL);
	is_true(options[2].vars == NULL);
}

static struct hig{
	char *ft;
	int att[16];
	char *pat;	
	int end;	
} higs[] = {
	{},
	{"c", {5}, "q"},
	{"2", {4}, "w"},
	{},
};
void test_hig()
{
	is_true  (higs[0].ft     == NULL);
	// TODO : is_true  (higs[0].att[0] == NULL);
	is_true  (higs[0].pat    == NULL);

	is_streq(higs[1].ft    , "c");
	is_eq   (higs[1].att[0],  5 );
	is_streq(higs[1].pat   , "q");

	is_streq(higs[2].ft    , "2");
	is_eq   (higs[2].att[0],  4 );
	is_streq(higs[2].pat   , "w");

	is_true  (higs[3].ft     == NULL);
	// TODO : is_true  (higs[3].att[0] == NULL);
	is_true  (higs[3].pat    == NULL);
}

typedef struct hig hug;
void test_hug()
{
	hug hugs[] = {
		{"c", {5}, "q"},
		{"2", {4}, "w"},
	};
	is_streq(hugs[0].ft    , "c");
	is_eq   (hugs[0].att[0],  5 );
	is_streq(hugs[0].pat   , "q");

	is_streq(hugs[1].ft    , "2");
	is_eq   (hugs[1].att[0],  4 );
	is_streq(hugs[1].pat   , "w");
}

struct poz{
	char           y[5];
	char          *c   ;
	double         d[2];
	struct poz    *ppt ;
	struct hig  parr[3];
	hug            h[9];
} pozes[] = {
	{"dream","home", {1, 2}},
	{"hold" ,"a"   , {0,42}},
	{},
};
void test_poz()
{
	is_streq(pozes[0].y   , "dream");
	is_streq(pozes[0].c   , "home" );
	is_eq   (pozes[0].d[0], 1      );
	is_eq   (pozes[0].d[1], 2      );
	
	is_streq(pozes[1].y   , "hold" );
	is_streq(pozes[1].c   , "a"    );
	is_eq   (pozes[1].d[0], 0      );
	is_eq   (pozes[1].d[1], 42     );

	is_not_null(&pozes[2]);
}

void test_ab()
{
	char * num[3] = {{"123"},{"987"},{"456"}};
	is_streq(num[0], "123");
	is_streq(num[1], "987");
	is_streq(num[2], "456");
}

void test_vti()
{
    int aqq[1][1] = { { 5 } };
	is_eq(aqq[0][0], 5);
}

void test_bm()
{
	char brac[4][3][2] = { {"1", "2"} };
	is_streq(brac[0][1], "2");
}

int main()
{
    plan(148);

	START_TEST(partly);

	// Test partly initialization of FILE
	START_TEST(FILE)

	// Test partly initialization of void
	START_TEST(void)

    START_TEST(array_float);
    START_TEST(array_char);
    START_TEST(struct_init);
    START_TEST(matrix_double);
	START_TEST(equals_chars);
	START_TEST(di);
	START_TEST(options);
	START_TEST(ex);
	START_TEST(hig);
	START_TEST(hug);
	START_TEST(poz);
	START_TEST(ab);
	START_TEST(vti);
	START_TEST(bm);
	
    done_testing();
}
