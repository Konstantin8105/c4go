int * f(int * s)
{
	return s;
}

struct A {
	int    a;
	int   *ap;

	double    b;
	double   *bp;
};

struct B {
	struct A    a;
	struct A   *ap;
};
