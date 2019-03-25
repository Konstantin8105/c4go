// Tests for ctype.h.

#include "tests.h"
#include <ctype.h>
#include <stdio.h>

#define T is_true
#define F is_false

#define CTYPE(functionName, A, B, C, D, E, F, G, H) \
    diag(#functionName);                            \
    A(functionName('a'));                           \
    B(functionName('B'));                           \
    C(functionName('0'));                           \
    D(functionName('*'));                           \
    E(functionName('\0'));                          \
    F(functionName(' '));                           \
    G(functionName('\n'));                          \
    H(functionName('z'));

// This is just helpful for alignment.
#define _CTYPE CTYPE

char* strnul = "this string has a \0 NUL";
char arrnul[] = "this string has a \0 NUL";

static int clrcomp(char *s, int len)
{
	static char *digs = "0123456789abcdef";
	int i;
	for (i = 0; i < len; i++) { 
		if (strchr(digs, tolower(s[i]))) {
			printf("%d\n", (int)(*strchr(digs, tolower(s[i]))));
			return (strchr(digs, tolower(s[i])) - digs);
		}
	}
	return -1;
}

int main()
{
    plan(114);

    //              . Lower alpha (a)
    //              |  . Upper alpha (B)
    //              |  |  . Digit (0)
    //              |  |  |  . Punctuation (*)
    //              |  |  |  |  . NULL
    //              |  |  |  |  |  . Space
    //              |  |  |  |  |  |  . New line
    //              |  |  |  |  |  |  |  . Non-hex digit (z)
    //              v  v  v  v  v  v  v  v
    _CTYPE(isalnum, T, T, T, F, F, F, F, T);
    _CTYPE(isalpha, T, T, F, F, F, F, F, T);
    _CTYPE(iscntrl, F, F, F, F, T, F, T, F);
    _CTYPE(isdigit, F, F, T, F, F, F, F, F);
    _CTYPE(isgraph, T, T, T, T, F, F, F, T);
    _CTYPE(islower, T, F, F, F, F, F, F, T);
    _CTYPE(isprint, T, T, T, T, F, T, F, T);
    _CTYPE(ispunct, F, F, F, T, F, F, F, F);
    _CTYPE(isspace, F, F, F, F, F, T, T, F);
    _CTYPE(isupper, F, T, F, F, F, F, F, F);
    _CTYPE(isblank, F, F, F, F, F, T, F, F);
    CTYPE(isxdigit, T, T, T, F, F, F, F, F);

    diag("tolower");
    is_eq(tolower('a'), 'a');
    is_eq(tolower('B'), 'b');
    is_eq(tolower('0'), '0');
    is_eq(tolower('*'), '*');
    is_eq(tolower('\0'), '\0');
    is_eq(tolower(' '), ' ');
    is_eq(tolower('\n'), '\n');
    is_eq(tolower('z'), 'z');

    diag("toupper");
    is_eq(toupper('a'), 'A');
    is_eq(toupper('B'), 'B');
    is_eq(toupper('0'), '0');
    is_eq(toupper('*'), '*');
    is_eq(toupper('\0'), '\0');
    is_eq(toupper(' '), ' ');
    is_eq(toupper('\n'), '\n');
    is_eq(toupper('z'), 'Z');

	diag("clrcomp");
	{
		char *word1 = "12";
		is_eq(clrcomp(word1,2),1);
		char *word2 = "qwe12";
		is_eq(clrcomp(word2,5),14);
	}

    done_testing();
}
