int * f(int * s)
{
	int a = *s;
	int * b = s;
	int * c = &a;
	s = f(&a);
	s = f(b);
	s = f(c);
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
