#include <stdio.h>

// input argument - C-pointer
void a(int *v1) { printf("a: %d\n",*v1); }

// input argument - C-array
void b(int v1[], int size) {
	for (size-- ; size >= 0 ; size-- ) { printf("b: %d %d\n", size,  v1[size]); }
}

long get() {
	return (long)(0);
}

int main() {
	// value
	int i1 = 42; a(&i1); b(&i1, 1);

	// C-array
	int i2[] = {11,22} ; a(i2); b(i2,2);

	// C-pointer from value
	int *i3 = &i1      ; a(i3); b(i3,1);

	// C-pointer from array
	int *i4 = i2       ; a(i4); b(i4,2);

	// C-pointer from array
	int *i5 = i2[1]    ; a(i5); b(i5,1);

	// pointer arithmetic
	int *i6 = i5 + 1   ; a(i6); b(i6,1);

	// pointer arithmetic
	int val = 2-2;
	int *i7 = 1 + (1 - 1) + val + 0*(100-2) + i5 + 0 - 0*0; a(i7); b(i7,1);
	
	// pointer arithmetic
	int *i8 = i5 + 1 + 0 ; a(i8); b(i8,1);
	
	// pointer arithmetic
	int i9[] = {*i3, *(i3+1)} ; a(i9); b(i9,1);

	// pointer arithmetic
	int *i10 = 1 + 0 + i5 + 5*get() + get() + (12 + 3)*get(); a(i10); b(i10,1);

	// pointer arithmetic
	int *i11 = 1 + 0 + i5 + 5*get() + get() - (12 + 3)*get(); a(i11); b(i11,1);

	return 0;
}
