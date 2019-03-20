#include "tests.h"

typedef struct SP SP;

struct SP {
    union SP_UN {
        int i;
        char l;
    };
    double d;
    char c;
};

typedef struct SP Mem;
typedef SP Mem2;

int main()
{
    plan(6);

    Mem m;
    is_true(sizeof(m) >= sizeof(double) + sizeof(char) + sizeof(union SP_UN));
    (void)(m);

    Mem2 m2;
    is_true(sizeof(m2) >= sizeof(double) + sizeof(char) + sizeof(union SP_UN));
    (void)(m2);

    SP s;
    is_true(sizeof(s) >= sizeof(double) + sizeof(char) + sizeof(union SP_UN));
    (void)(s);

    union SP_UN sp;
    is_true(sizeof(sp) >= sizeof(int));
    (void)(sp);

    is_eq(sizeof(double), 8);
    is_eq(sizeof(char), 1);

    done_testing();
}
