#define unsafetype(type,name)	\
void name(){					\
	type    t;					\
	type * pt = &t;				\
	(void)(t);					\
	(void)(pt);					\
}

struct str{
	int i;
};

union un{
	int i;
	double d;
};

// integers
unsafetype(char,test_char)
unsafetype(short,test_short)
unsafetype(int,test_int)
unsafetype(long,test_long)
unsafetype(long int, test_li)
unsafetype(long long,test_ll)
unsafetype(long long int, test_lli)

// floats
unsafetype(float, test_f)
unsafetype(double, test_d)
unsafetype(long double, test_ld)

// struct
unsafetype(struct str, test_struct)

// union
unsafetype(union un, test_un)
