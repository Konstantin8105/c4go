#include "tests.h"
#include <string.h>

int main()
{
    plan(36);

    diag("TODO: __builtin_object_size");
    // https://github.com/Konstantin8105/c4go/issues/359

    {
        diag("strcpy");
        char* src = "foo bar\0baz";
        char dest1[40];
        char* dest2;
        dest2 = strcpy(dest1, src);
        is_streq(dest1, "foo bar");
        is_streq(dest2, "foo bar");
    }

    diag("strncpy");

    // If the end of the source C string (which is signaled by a null-character)
    // is found before num characters have been copied, destination is padded
    // with zeros until a total of num characters have been written to it.
    {
        char* src = "foo bar\0baz";
        char dest1[40];
        char* dest2;

        dest1[7] = 'a';
        dest1[15] = 'b';
        dest1[25] = 'c';
        dest2 = strncpy(dest1, src, 20);

        is_eq(dest1[0], 'f');
        is_eq(dest1[7], 0);
        is_eq(dest1[15], 0);
        is_eq(dest1[25], 'c');

        is_eq(dest2[0], 'f');
        is_eq(dest2[7], 0);
        is_eq(dest2[15], 0);
        is_eq(dest2[25], 'c');

        is_streq(dest1, "foo bar");
        is_streq(dest2, "foo bar");
    }

    // No null-character is implicitly appended at the end of destination if
    // source is longer than num. Thus, in this case, destination shall not be
    // considered a null terminated C string (reading it as such would
    // overflow).
    {
        char* src = "foo bar\0baz";
        char dest1[40];
        char* dest2;

        dest1[7] = 'a';
        dest1[15] = 'b';
        dest1[25] = 'c';
        dest2 = strncpy(dest1, src, 5);

        is_eq(dest1[0], 'f');
        is_eq(dest1[7], 'a');
        is_eq(dest1[15], 'b');
        is_eq(dest1[25], 'c');

        is_eq(dest2[0], 'f');
        is_eq(dest2[7], 'a');
        is_eq(dest2[15], 'b');
        is_eq(dest2[25], 'c');
    }

    {
        diag("strlen");
        is_eq(strlen(""), 0);
        is_eq(strlen("a"), 1);
        is_eq(strlen("foo"), 3);
        // NULL causes a seg fault.
        // is_eq(strlen(NULL), 0);
        is_eq(strlen("fo\0o"), 2);
    }
    {
        diag("strcat");
        char str[80];
        for (int i = 0;i<80;i++) str[i] = 0;
        strcpy(str, "these ");
        strcat(str, "strings ");
        strcat(str, "are ");
        strcat(str, "concatenated.");
        is_streq(str, "these strings are concatenated.");
    }
    {
        diag("strncat");
        char str[80];
        for (int i = 0;i<80;i++) str[i] = 0;
        strncpy(str, "these", 3);
        strncat(str, " strings", 7);
        strncat(str, " is", 3);
        strncat(str, " concatenated.", 14);
        is_streq(str, "the string is concatenated.");
    }
    {
        diag("strcmp");
        {
            char* a = "ab";
            char* b = "ab";
            is_true(strcmp(a, b) == 0);
        }
        {
            char* a = "bb";
            char* b = "ab";
            is_true(strcmp(a, b) > 0);
        }
        {
            char* a = "ab";
            char* b = "bb";
            is_true(strcmp(a, b) < 0);
        }
    }
    {
        diag("strchr");
        char str[] = "This is a sample string";
        char* pch;
        int amount = 0;
        pch = strchr(str, 's');
        while (pch != NULL) {
            pch = strchr(pch + 1, 's');
            amount++;
        }
        is_eq(amount, 4);
    }
	{
		diag("memset");
		char str[] = "almost every programmer should know memset!";
		char * ptr = memset (str,'-',6);
		is_streq(str,"------ every programmer should know memset!");
		is_eq(ptr-str,0);
	}
	{
		diag("memmove");
		char str[] = "memmove can be very useful......";
		memmove (str+20,str+15,11);
		is_streq(str,"memmove can be very very useful.");
	}
	{
		diag("memcmp");
		char a1[] = {'a','b','c'};
		char a2[] = "abd";
		is_true(memcmp(a1, a1, 3) == 0);
		is_true(memcmp(a1, a2, 3) < 0);
		is_true(memcmp(a2, a1, 3) > 0);
	}

    done_testing();
}
