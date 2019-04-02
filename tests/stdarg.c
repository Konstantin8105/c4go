#include "tests.h"
#include <assert.h>
#include <stdarg.h>
#include <stdio.h>

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void simple(const char* fmt, ...)
{
    char buffer[155];
    for (int i = 0; i < 155; i++) {
        buffer[i] = 0;
    }
    is_streq(buffer, "");

    va_list args;
    va_start(args, fmt);

    char temp[100];
    for (int i = 0; i < 100; i++) {
        temp[i] = 0;
    }

    int len = 4;
    for (int i = 0; i < len; i++) {
        char f = fmt[i];
        if (f == 'd') {
            int i = va_arg(args, int);
            sprintf(temp, "%d ", i);
            strcat(buffer, temp);
            is_streq(buffer, "3 ")
        } else if (f == 'c') {
            // note automatic conversion to integral type
            int c = va_arg(args, int);
            sprintf(temp, "%c ", c);
            strcat(buffer, temp);
            is_streq(buffer, "3 a ")
        } else if (f == 'f') {
            double d = va_arg(args, double);
            sprintf(temp, "%.3f ", d);
            strcat(buffer, temp);
        }
    }

    va_end(args);

    is_streq(buffer, "3 a 1.999 42.500 ")
}

void test_va_list()
{
    simple("dcff", 3, 'a', 1.999, 42.5);
}

int sum(int num_args, ...)
{
    int val = 0;
    va_list ap;
    int i;

    va_start(ap, num_args);
    for (i = 0; i < num_args; i++) {
        val += va_arg(ap, int);
    }
    va_end(ap);

    return val;
}

void test_va_list2()
{
    is_eq(sum(3, 10, 20, 30), 60);
}

int strange(int num_args, ...)
{
    int val = 0;
    va_list ap;
    int i;

    va_start(ap, num_args);
    for (i = 0; i < num_args; i++) {
        *va_arg(ap, int*) += 2;
    }
    va_end(ap);

    va_start(ap, num_args);
    for (i = 0; i < num_args; i++) {
        val += *va_arg(ap, int*);
    }
    va_end(ap);

    return val;
}

void test_va_list3()
{
    int v1 = 10;
    int v2 = 23;
    is_eq(strange(2, &v1, &v2), 10 + 2 + 23 + 2);
}

void out(int num_args, va_list ap)
{
    for (int i = 0; i < num_args; i++) {
        int Y = va_arg(ap, int);
        printf("%d -> %d\n", i, Y);
    }
}

void red(int num_args, ...)
{
    va_list ap;
    // base test
    va_start(ap, num_args);
    out(num_args, ap);
    va_end(ap);
    // repeat test
    va_start(ap, num_args);
    out(num_args, ap);
    va_end(ap);
}

void test_va_list4()
{
    red(3, 12, 23, 34);
}

int main()
{
    plan(6);

    START_TEST(va_list)
    START_TEST(va_list2)
    START_TEST(va_list3)
    START_TEST(va_list4)

    done_testing();
}
