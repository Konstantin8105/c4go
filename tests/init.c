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


// static char *kmap_fa[256] = {
	// [0] = "fa",
	// ['`'] = "‍",
	// ['1'] = "۱",
// };
// 
// static char *digraphs[][2] = {
	// {"cq", "’"},
	// {"pl", "+"},
	// {"hy", "-"},
	// {"sl", "/"},
// };
// 
// static int xai, xaw;
// 
// static struct option {
	// char *abbr;
	// char *name;
	// int *var;
// } options[] = {
	// {"ai", "autoindent", &xai},
	// {"aw", "autowrite", &xaw},
// };
// 
// int ec_print(char *s){
	// printf("%s\n",s);
	// return 32;
// }
// 
// int ec_insert(char *s){
	// printf("%s\n",s);
	// return 42;
// }
// 
// static struct excmd {
	// char *abbr;
	// char *name;
	// int (*ec)(char *s);
// } excmds[] = {
	// {"p", "print", ec_print},
	// {"a", "append", ec_insert},
// };

int main()
{
    plan(13);

    START_TEST(array_float);
    START_TEST(array_char);
    START_TEST(struct_init);
    START_TEST(matrix_double);

    done_testing();
}
