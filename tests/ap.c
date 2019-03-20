#include "tests.h"
#include <stdio.h>

// input argument - C-pointer
void a(int* v1) { printf("a: %d\n", *v1); }

// input argument - C-array
void b(int v1[], int size)
{
    for (size--; size >= 0; size--) {
        printf("b: %d %d\n", size, v1[size]);
    }
}

long get()
{
    return (long)(0);
}

double global;

double* get_value()
{
    return &global;
}

int main()
{
    plan(4);

    diag("value");
    int i1 = 42;
    a(&i1);
    b(&i1, 1);

    diag("C-array");
    int i2[] = { 11, 22 };
    a(i2);
    b(i2, 2);

    diag("C-pointer from value");
    int* i3 = &i1;
    a(i3);
    b(i3, 1);

    diag("C-pointer from array");
    int* i4 = i2;
    a(i4);
    b(i4, 2);

    diag("C-pointer from array");
    int* i5 = &i2[1];
    a(i5);
    b(i5, 1);

    diag("pointer arithmetic 1");
    i5 = &i2[0];
    int* i6 = i5 + 1;
    a(i6);
    b(i6, 1);

    diag("pointer arithmetic 2");
    int val = 2 - 2;
    int* i7 = 1 + (1 - 1) + val + 0 * (100 - 2) + i5 + 0 - 0 * 0;
    a(i7);
    b(i7, 1);

    diag("pointer arithmetic 3");
    int* i8 = i5 + 1 + 0;
    a(i8);
    b(i8, 1);

    diag("pointer arithmetic 4");
    int i9[] = { *i3, *(i3 + 1) };
    a(i9);
    b(i9, 1);

    diag("pointer arithmetic 5");
    int* i10 = 1 + 0 + i5 + 5 * get() + get() + (12 + 3) * get();
    a(i10);
    b(i10, 1);

    diag("pointer arithmetic 6");
    int* i11 = 1 + 0 + i5 + 5 * get() + get() - (12 + 3) * get();
    a(i11);
    b(i11, 1);

    diag("pointer from function");
    global = 42;
    double* i12 = get_value();
    is_eq(*i12, global);

    diag("pointer for some type");
    typedef int NUMBER;
    typedef int* POINTER;
    NUMBER num;
    POINTER pnt;
    int num_value = 56;
    num = num_value;
    pnt = &num;
    is_eq(*pnt, num_value);
    *pnt = 123;
    is_eq(*pnt, 123);
    is_eq(num, 123);

    done_testing();
}
