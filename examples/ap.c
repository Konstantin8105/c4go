#include <stdio.h>

void a(int *v1) { printf("a: %d\n",*v1); }

void b(int v1[], int size) {
	for (size-- ; size >= 0 ; size-- ) { printf("b: %d %d\n", size,  v1[size]); }
}

int main() {
	int i1 = 42;
	a(&i1);
	b(&i1, 1);

	int i2[] = {11,22};
	a(i2);
	b(i2,2);

	return 0;
}
