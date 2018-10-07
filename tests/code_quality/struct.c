struct T {
	int     *a;
	float   *f;
	double  *d;
	int      q;
	float    w;
	double   r;
};

struct T2 {
	struct T *t;
	struct T  y;
};

int * f(double * a, struct T2 * t)
{
	int * r;
	(void) a;
	(void) t;
	return r;
}
