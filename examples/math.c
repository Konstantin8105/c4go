#include <math.h>
#include <stdio.h>

int main()
{
    int n;
    double param = 8.0, result;
    result = frexp(param, &n);
    printf("result = %5.2f\n", result);
    printf("n      = %d\n", n);
    return 0;
}
