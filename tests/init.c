#include "tests.h"



static char *kmap_fa[256] = {
	[0] = "fa",
	['`'] = "‍",
	['1'] = "۱",
};

static char *digraphs[][2] = {
	{"cq", "’"},
	{"pl", "+"},
	{"hy", "-"},
	{"sl", "/"},
};

static int xai, xaw;

static struct option {
	char *abbr;
	char *name;
	int *var;
} options[] = {
	{"ai", "autoindent", &xai},
	{"aw", "autowrite", &xaw},
};

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

