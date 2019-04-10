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

int main()
{
    // value
    int i1 = 42;
    a(&i1);
    b(&i1, 1);

    // C-array
    int i2[] = { 11, 22 };
    a(i2);
    b(i2, 2);

    // C-pointer from value
    int* i3 = &i1;
    a(i3);
    b(i3, 1);

    // C-pointer from array
    int* i4 = i2;
    a(i4);
    b(i4, 2);

    // C-pointer from array
    int* i5 = i2[1];
    a(i5);
    b(i5, 1);

    return 0;
}
