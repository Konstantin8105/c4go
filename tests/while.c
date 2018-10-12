#include "tests.h"
#include <assert.h>
#include <ctype.h>
#include <stdio.h>

typedef float* triangle;

int getNumberToken(char a[], int* cursor)
{
    int value = 0;
    assert((*cursor) >= 0 && (*cursor) < strlen(a));
    do {
        value = value * 10 + (int)(a[*cursor]) - (int)('0');
        (*cursor)++;
    } while (isdigit(a[*cursor]));
    return value;
}

int main()
{
    plan(10);

    int value = 1;

    while (value <= 3) {
        pass("value is %d", value);
        value++;
    }

    // continue
    value = 0;
    while (value < 3) {
        value++;
        if (value < 3)
            continue;
        pass("%d", value);
    }

    diag("while without body");
    while (0)
        ;
    pass("%s", "while without body");

    value = 1;
    while ((value--, value))
        ;
    is_eq(value, 0);

    diag("while with star");
    {
        int* ok;
        int value2;
        ok = &value2;
        *ok = 1;
        int counter = 0;
        do {
            counter++;
            if (counter > 10) {
                break;
            }
            if (counter < 10) {
                break;
            }
            *ok = *ok - 1;
        } while (*ok);
        is_eq(*ok, 1);
    }

    diag("while with --");
    {
        int T = 2;
        int counter = 0;
        while (T--) {
            if (counter > 50) {
                break;
            }
        };
        is_eq(T, -1);
    }
    diag("while in triangle");
    {
        float f = 87.76;
        triangle* newtriangle;
        triangle value[50];
        for (int y; y < 50; y++) {
            value[y] = (triangle)(&f);
        }
        newtriangle = &value;
        int counter = 0;
        do {
            counter++;
            if (counter > 10) {
                break;
            }
            if (counter < 10) {
                break;
            }
        } while ((newtriangle)[counter] == (triangle)NULL);
        is_eq(counter, 1);
    }

    int v = 3;
    int res = getNumberToken("my test string", &v);
    is_eq(res, 68);

    done_testing();
}
