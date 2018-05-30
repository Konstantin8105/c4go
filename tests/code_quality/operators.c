void f(unsigned d){
	(void)(d);
}

void operators_equals()
{
	int a,b,c,d;
	a=b=c=d=42;

	int        q1 = 42;
	unsigned   q2 = 42;
	long       q3 = 42;
	long long  q4 = 42;
	short      q5 = 42;

	int        w1 = (long long)(42);
	unsigned   w2 = (long long)(42);
	long       w3 = (long long)(42);
	long long  w4 = (long long)(42);
	short      w5 = (long long)(42);

	// zero value
	int        s1 = 0;
	unsigned   s2 = 0;
	long       s3 = 0;
	long long  s4 = 0;
	short      s5 = 0;

	double     d1 = 0;
	float      d2 = 0;

	double     v1 = 0.0;
	float      v2 = 0.0;
	
	int        x1 =(0);
	unsigned   x2 =(0);
	long       x3 =(0);
	long long  x4 =(0);
	short      x5 =(0);

	double     z1 =(0);
	float      z2 =(0);

	double     t1 =(0.0);
	float      t2 =(0.0);

	// function
	f(42);
}
