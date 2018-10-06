#include "tests.h"
#include <stdlib.h>

void test_struct_pointer()
{
    typedef struct row {
        double* t;
    } row;
    row r;
    unsigned poi = r.t;
    is_true(poi >= 0);
    (void)(r);
}

int main()
{
    plan(1);

    test_struct_pointer();

    done_testing();
}
