// Initialization structs, arrays, ...

#include "tests.h"

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

static int xai, xaw;
static struct option {
	char *abbr;
	char *name;
	int *var;
} options[] = {
	{"ai", "autoindent", &xai},
	{"aw", "autowrite", &xaw},
};
void test_options()
{
	is_streq(options[0].abbr, "ai");
	is_streq(options[0].name, "autoindent");
	is_not_null(options[0].var);

	is_streq(options[1].abbr, "aw");
	is_streq(options[1].name, "autowrite");
	is_not_null(options[1].var);
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
};

void test_ex()
{
	is_streq(excmds[0].abbr, "p")
	is_streq(excmds[0].name, "print")
	is_eq(excmds[0].ec("0"), 32);
	
	is_streq(excmds[1].abbr, "a")
	is_streq(excmds[1].name, "append")
	is_eq(excmds[1].ec("1"), 42);
}

static struct hig{
	char *ft;		/* the filetype of this pattern */
	int att[16];		/* attributes of the matched groups */
	char *pat;		/* regular expression */
	int end;		/* the group ending this pattern */
} higs[] = {
	{"c", {5}, "q"},
	{"2", {4}, "w"},
};
void test_hig()
{
	is_streq(higs[0].ft    , "c");
	is_eq   (higs[0].att[0],  5 );
	is_streq(higs[0].pat   , "q");

	is_streq(higs[1].ft    , "2");
	is_eq   (higs[1].att[0],  4 );
	is_streq(higs[1].pat   , "w");
}


int main()
{
    plan(42);

    START_TEST(array_float);
    START_TEST(array_char);
    START_TEST(struct_init);
    START_TEST(matrix_double);
	START_TEST(equals_chars);
	START_TEST(di);
	START_TEST(options);
	START_TEST(ex);
	START_TEST(hig);

    done_testing();
}
