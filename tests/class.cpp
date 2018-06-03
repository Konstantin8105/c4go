#include "tests.h"

// Class Declaration
class person {
    //Access - Specifier
public:

    //Variable Declaration
    float name;
    int number;
};

int main()
{
    plan(2);

	person obj;
	obj.name = 2.3;
	obj.number = 4;

	is_eq(obj.name,2.3);
	is_eq(obj.number,4);

    done_testing();
}
