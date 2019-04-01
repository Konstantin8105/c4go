// Tests for structures.

#include "tests.h"
#include <stdio.h>

typedef void function_t(int args);

struct programming {
    float constant;
    char* pointer;
};

void pass_by_ref(struct programming* addr)
{
    char* s = "Show string member.";
    float v = 1.23 + 4.56;

    addr->constant += 4.56;
    addr->pointer = s;

    is_eq(addr->constant, v);
    is_streq(addr->pointer, "Show string member.");
}

void pass_by_val(struct programming value)
{
    value.constant++;

    is_eq(value.constant, 2.23);
    is_streq(value.pointer, "Programming in Software Development.");
}

/**
 * text
 */
typedef struct mainStruct {
    double constant;
} secondStruct;

/*
 * Text
 */
typedef struct {
    double t;
} ts_c;

// Text
typedef struct ff {
    int v1, v2;
} tt1, tt2;

// Text1
// Text2
struct outer {
    int i;
    struct z {
        int j;
    } inner;
};

struct xx {
    int i;
    /**
	 * Text
	 */
    struct yy {
        int j;
        struct zz {
            int k;
        } deep;
    } inner;
};

/**
 * Some function
 */
int summator(int i, float f)
{
    return i + (int)(f);
}

typedef struct J J;
struct J {
    float f;
    int (*fu)(J* j, float i);
};

int j_function(J* j, float i)
{
    if (j != NULL) {
        return (int)(i + (*j).f);
    }
    return -1;
};

void struct_with_rec_fuction()
{
    J j;
    j.f = 5.0;
    j.fu = j_function;
    is_eq(j.fu(&j, 4.0), 9);
    is_eq(j_function(NULL, 4.0), -1);
}

struct FinFinS {
    double d;
    int (*f)(int (*)(int));
};

int FinF1(int a)
{
    return a + 1;
}

int FinF2(int (*f)(int))
{
    int g = 45;
    return f(g);
}

void func_in_func_in_struct()
{
    diag("function in function in struct");
    struct FinFinS ffs;
    ffs.f = FinF2;
    int res = ffs.f(FinF1);
    is_eq(res, 46);
};

struct info {
    struct deep_info {
        int a, b, c;
    } con;
    struct star_deep_info {
        int sa, sb, sc;
    } * star_con;
};

void struct_in_struct_with_star()
{
    diag("struct in struct with star");
    struct info in;

    in.con.a = 45;
    is_eq(in.con.a, 45);

    struct star_deep_info st;
    in.star_con = &st;

    in.star_con->sa = 45;
    is_eq(in.star_con->sa, 45);

    struct info in2;

    struct star_deep_info st2;
    in2.star_con = &st2;

    in2.star_con->sa = 46;
    is_eq(in2.star_con->sa, 46);
    is_eq(in.star_con->sa, 45);
}

struct {
    int a;
} m;
typedef double** y;

struct {
    int aa;
} mm;
typedef double** yy;

struct mmm {
    int aaa;
};
typedef double*** yyy;

struct mmmm {
    int aaaa;
};
typedef double**** yyyy;

struct mmmmm0 {
    int aaaaa;
} mmmmm;
typedef double** yyyyy;

typedef struct {
    int st1;
} st2;
typedef double** st3;

typedef struct st4 {
    int st5;
} st6;
typedef double** st7;
struct st4 st7a;

typedef struct st4a {
    int st5a;
} * st6a;

typedef struct st4b {
    int st5b;
} * const* st6b;

struct st8 {
    int st9;
    struct st10 {
        int st11;
    };
};
typedef double** st12;

struct st13 {
    int st14;
    struct st16 {
        int st17;
    } st18;
} st19;
typedef double** st20;

typedef struct st21 {
    int st22;
    struct st23 {
        int st24;
    } st25;
} st26;
typedef double** st27;

static struct unix_syscall {
    const char* zName;
} aSyscall[] = {
    { "open" },
    { "close" }
};

struct memory {
    int* one;
    float** oop;
    double* two;
    struct memory* mm;
};

typedef double** subseg;

struct mesh {
    subseg* dummysub;
};

double* returner(int* const* i, double* d)
{
    (void)(i);
    return d;
}

void struct_null()
{
    struct memory dm;
    float o = 99;
    float* oo = &o;
    float** ooo = &oo;
    dm.oop = ooo;
    struct memory* m = &dm;
    m->one = (int*)(NULL);
    m->two = (double*)(NULL);
    m->mm = (struct memory*)(NULL);
    m->one = (void*)(NULL);
    m->two = (void*)(NULL);
    m->mm = (void*)(NULL);
    *(m->oop) = (int*)NULL;
    (m->oop) = (int*)NULL;
    (void)(dm);
    (void)(m);

    (void)summator(1, 34.4);
    (void)returner(0, 0);
    double fd = 56;
    returner(0, &fd);
    (void)(fd);

    static const struct {
        const char* zPattern;
        const char* zDesc;
    } aTrans[] = {
        { "rchar: ", "Bytes received by read():" },
        { "wchar: ", "Bytes sent to write():" },
        { "syscr: ", "Read() system calls:" },
        { "syscw: ", "Write() system calls:" },
        { "read_bytes: ", "Bytes read from storage:" },
        { "write_bytes: ", "Bytes written to storage:" },
        { "cancelled_write_bytes: ", "Cancelled write bytes:" },
    };

    is_eq(strlen(aTrans[3].zPattern), 7);
    is_streq(aTrans[2].zPattern, "syscr: ");
    is_streq(aTrans[1].zDesc, "Bytes sent to write():");

    double d = 99;
    double* dd = &d;
    double** ddd = &dd;
    *(ddd) = (int*)NULL;
    (void)(ddd);

    struct memorypool {
        int** nowblock;
    };
    struct memorypool Vpool;
    int nowblock;
    int* s_nowblock = &nowblock;
    Vpool.nowblock = &s_nowblock;
    is_not_null(*Vpool.nowblock);
    *(Vpool.nowblock) = NULL;
    struct memorypool* pool = &Vpool;
    if (*(pool->nowblock) == (int*)NULL) {
        pass("ok");
    }

    mm.aa = 42;
    is_eq(mm.aa, 42);

    struct mesh msh;
    subseg sub[10];
    for (int i = 0; i < 10; i++) {
        sub[i] = (subseg)(ddd);
    }
    msh.dummysub = sub;
    struct mesh* ms = &msh;
    ms->dummysub[2] = (subseg)NULL;
    (void)(ms);
    pass("ok");
}

union STRS {
    double d;
    struct {
        double d;
    } T;
};

void struct_inside_union()
{
    union STRS s;
    s.T.d = 10.0;
    is_true(s.d != 0);
}

struct FFS {
    void (*xDlSym)(int*, void*, const char* zSymbol);
};
int global_ffs = 0;
void ffs_i1(int* i, void* v, const char* ch)
{
    global_ffs++;
}

void (*ffs_i2(int* i, void* d, const char* zSymbol))(void)
{
    return ffs_i1;
}

void struct_func_func()
{
    struct FFS f;
    f.xDlSym = ffs_i1;
    is_eq(global_ffs, 0);
    f.xDlSym(NULL, NULL, NULL);
    is_eq(global_ffs, 1);
}

struct empty_str {
};
typedef struct sqlite3_file sqlite3_file;
struct sqlite3_file {
    const struct sqlite3_io_methods* pMethods; /* Methods for an open file */
};

typedef struct sqlite3_io_methods sqlite3_io_methods;
struct sqlite3_io_methods {
    int iVersion;
    int (*xClose)(sqlite3_file*);
};

void struct_after_struct()
{
    sqlite3_file sFile;
    sqlite3_io_methods io;
    sFile.pMethods = &io;
    is_not_null(sFile.pMethods);
}

struct RRR {
    struct sColMap { /* Mapping of columns in pFrom to columns in zTo */
        int iFrom; /* Index of column in pFrom */
        char* zCol; /* Name of column in zTo.  If NULL use PRIMARY KEY */
    } aCol[1]; /* One entry for each of nCol columns */
};

void struct_array()
{
    struct RRR rrr;
    rrr.aCol[0].iFrom = 10;
    is_eq(rrr.aCol[0].iFrom, 10);
}

typedef struct ss {
    char id;
} ss;
int struct_sizeof()
{
    ss v;
    v.id = 'e';
    ss* p = &v;
    is_eq(sizeof p->id, 1);
    (void)(p);
    return -1;
}

static const struct {
    int eType; /* Transformation type code */
    int nName; /* Length of th name */
    char* zName; /* Name of the transformation */
    double rLimit; /* Maximum NNN value for this transform */
    double rXform; /* Constant used for this transform */
} aXformType[] = {
    { 0, 6, "second", 464269060800.0, 86400000.0 / (24.0 * 60.0 * 60.0) },
    { 0, 6, "minute", 7737817680.0, 86400000.0 / (24.0 * 60.0) },
    { 0, 4, "hour", 128963628.0, 86400000.0 / 24.0 },
    { 0, 3, "day", 5373485.0, 86400000.0 },
    { 1, 5, "month", 176546.0, 30.0 * 86400000.0 },
    { 2, 4, "year", 14713.0, 365.0 * 86400000.0 },
};
void test_sizeofArray()
{
    is_eq((int)(sizeof(aXformType) / (sizeof(aXformType[0]))), 6);
}

struct StructBase {
    union {
        struct StructUsed* pStr;
        int aaa;
    } InsideUnion;
};
struct StructUsed {
    int vars;
};
void test_structUsed()
{
    struct StructBase sb;
    struct StructUsed ss;
    sb.InsideUnion.pStr = &ss;
    (*sb.InsideUnion.pStr).vars = 10;
    is_eq((*sb.InsideUnion.pStr).vars, 10);
}

struct EmptyName {
    union {
        long l1;
        long l2;
    };
};
void test_emptyname()
{
    struct EmptyName en;
    en.l1 = 10;
    is_eq(en.l1, 10);
    is_eq(en.l2, 10);
}

// Link: http://en.cppreference.com/w/c/language/typedef
typedef int A[]; // A is int[]
A a = { 1, 2 }, b = { 3, 4, 5 }; // type of a is int[2], type of b is int[3]
void test_typedef1()
{
    is_eq(a[1], 2);
    is_eq(b[1], 4);
}
void test_typedef2()
{
    typedef float A[];
    A a = { 1., 2. }, b = { 3., 4., 5. };
    is_eq(a[1], 2.);
    is_eq(b[1], 4.);
}

void test_pointer_member()
{
    struct tttt {
        int r[10];
        int* p;
    };
    struct tttt t;
    t.p = t.r;
    for (int i = 0; i < 10; i++)
        t.r[i] = i;
    is_eq(*t.p, 0);
    t.p++;
    is_eq(*t.p, 1);
    t.p += 2;
    is_eq(*t.p, 3);
}

struct SBA {
    int t;
    char name[100];
};

void struct_byte_array()
{
    struct SBA sba = { 10, "qwe" };
    is_eq(sba.t, 10);
    is_streq(sba.name, "qwe");
}

struct AaA;
typedef struct AaA tAaA;
struct AaA {
    int i;
};
void struct_typ2()
{
    tAaA a;
    a.i = 10;
    is_eq(a.i, 10);
    struct AaA b = a;
    is_eq(b.i, 10);
}

typedef struct map_node_s map_node_t;

struct map_node_s {
    unsigned hash;
    map_node_t* next;
};

map_node_t* ret()
{
    map_node_t mt;
    mt.hash = 42;
    return &mt;
}

void test_map_resize()
{
    map_node_t mt;
    mt.hash = 12;
    map_node_t* n = &mt;
    is_not_null(n);
    is_true(n->hash == 12);
    map_node_t* n2 = ret();
    is_not_null(n2);
    is_true(n2->hash == 42);
    n2->hash = 15;
    is_true(n2->hash == 15);
}

typedef float ext_vec;
extern ext_vec Re;
void test_extern_vec()
{
    ext_vec e;
    Re = 12.0;
    e = 5.0;
    is_eq(e, 5.0);
    is_eq(Re, 12.0);
}
ext_vec Re;

struct wb {
    int i;
};

int wb_test(struct wb* wb)
{
    return wb->i;
}

void test_same_name()
{
    diag("=== same name ===");
    {
        struct wb wb;
        wb.i = 42;
        is_eq(wb_test(&wb), 42);
    }
    {
        struct wb tt;
        tt.i = 42;
        struct wb* wb;
        wb = &tt;
        is_eq(wb_test(wb), 42);
    }
    {
        struct wb wb[2];
        wb[0].i = 42;
        wb[1].i = 42;
        is_eq(wb_test(&(wb[0])), 42);
    }
    diag("=================");
}

typedef int pointx;
typedef struct {
    pointx x;
    int y;
} Point2;
const Point2 p2[] = { {.y = 4, .x = 5 } };
const Point2* getPoint(int index)
{
    return &(p2[index]);
}
typedef unsigned char pcre_uchar;
typedef unsigned char pcre_uint8;
typedef unsigned int pcre_uint32;
typedef struct spu {
    pcre_uchar* hvm;
} spu;

void pointer_arithm_in_struct()
{
    pcre_uchar str[] = "abcd";
    spu s;
    spu* ps = &s;
    ps->hvm = &str[1];
    is_true(ps->hvm == &str[1]);
    ps->hvm += 2;
    is_true(ps->hvm == &str[3]);
}

int main()
{
    plan(100);

    pointer_arithm_in_struct();
    test_extern_vec();
    test_map_resize();
    struct_typ2();
    struct_byte_array();
    test_pointer_member();
    test_typedef1();
    test_typedef2();
    struct_array();
    struct_func_func();
    struct_after_struct();
    struct_sizeof();
    test_sizeofArray();
    test_structUsed();
    test_emptyname();
    test_same_name();

    struct programming variable;
    char* s = "Programming in Software Development.";

    variable.constant = 1.23;
    variable.pointer = s;

    is_eq(variable.constant, 1.23);
    is_streq(variable.pointer, "Programming in Software Development.");

    pass_by_val(variable);
    pass_by_ref(&variable);

    struct mainStruct s1;
    s1.constant = 42.;
    is_eq(s1.constant, 42.);

    secondStruct s2;
    s2.constant = 42.;
    is_eq(s2.constant, 42.);

    ts_c c;
    c.t = 42.;
    is_eq(c.t, 42.);

    tt1 t1;
    t1.v1 = 42.;
    is_eq(t1.v1, 42.)

        tt2 t2;
    t2.v1 = 42.;
    is_eq(t2.v1, 42.)

        struct ff tf2;
    tf2.v2 = t1.v1;
    is_eq(tf2.v2, t1.v1);

    struct outer o;
    o.i = 12;
    o.inner.j = 34;
    is_eq(o.i + o.inner.j, 46);

    struct xx x;
    x.i = 12;
    x.inner.j = 34;
    x.inner.deep.k = 56;
    is_eq(x.i + x.inner.j + x.inner.deep.k, 102);

    struct u {
        int y;
    };
    struct u yy;
    yy.y = 42;
    is_eq(yy.y, 42);

    diag("Typedef struct with same name");
    {
        typedef struct Uq Uq;
        struct Uq {
            int uq;
        };
        Uq uu;
        uu.uq = 42;
        is_eq(uu.uq, 42);
    }

    diag("Initialization of struct");
    struct Point {
        int x;
        int y;
    };
    struct Point p = {.y = 2, .x = 3 };
    is_eq(p.x, 3);
    is_eq(p.y, 2);

    diag("ImplicitValueInitExpr");
    {
        typedef struct {
            int x2;
            int y2;
        } coord2;

        typedef struct {
            coord2 position2;
            int possibleSteps2;
        } extCoord2;

        extCoord2 followingSteps[2] = {
            {.possibleSteps2 = 1 }, {.possibleSteps2 = 1 },
        };
        is_eq(followingSteps[0].possibleSteps2, 1);
    }
    {
        struct coord {
            int x;
            int y;
        };

        struct extCoord {
            struct coord position;
            int possibleSteps;
        };

        struct extCoord followingSteps[2] = {
            {.possibleSteps = 1 }, {.possibleSteps = 1 },
        };
        is_eq(followingSteps[0].possibleSteps, 1);
    }

    diag("Double typedef type");
    {
        typedef int int2;
        typedef int2 int3;
        typedef int3 int4;

        is_eq((int)((int4)((int3)((int2)(42)))), 42);
    }
    {
        typedef size_t size2;
        is_eq(((size2)((size_t)(56))), 56.0)
    }
    {
        is_eq((size_t)(43), 43);
    }

    diag("Function pointer inside struct");
    {
        struct F1 {
            int x;
            int (*f)(int, float);
        };
        struct F1 f1;
        f1.x = 42;
        f1.f = summator;
        is_eq(f1.x, 42);
        is_eq(f1.f(3, 5), 8);
    }
    {
        typedef struct {
            int x;
            int (*f)(int, float);
        } F2;
        F2 f2;
        f2.x = 42;
        f2.f = summator;
        is_eq(f2.x, 42);
        is_eq(f2.f(3, 5), 8);
    }

    diag("typedef function");
    {
        typedef int ALIAS(int, float);
        ALIAS* f = summator;
        is_eq(f(3, 5), 8);
    }
    {
        typedef int ALIAS2(int, float);
        ALIAS2* f;
        f = summator;
        is_eq(f(3, 5), 8);
    }

    diag("typedef struct C C inside function");
    {
        typedef struct CCC CCC;
        struct CCC {
            float ff;
        };
        CCC c;
        c.ff = 3.14;
        is_eq(c.ff, 3.14);
    }
    typedef struct CP CP;
    struct CP {
        float ff;
    };
    CP cp;
    cp.ff = 3.14;
    is_eq(cp.ff, 3.14);

    diag("struct name from Go keyword");
    {
        struct chan {
            int i;
        };
        struct chan UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct defer {
            int i;
        };
        struct defer UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct fallthrough {
            int i;
        };
        struct fallthrough UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct func {
            int i;
        };
        struct func UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct go {
            int i;
        };
        struct go UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct import {
            int i;
        };
        struct import UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct interface {
            int i;
        };
        struct interface UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct map {
            int i;
        };
        struct map UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct package {
            int i;
        };
        struct package UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct range {
            int i;
        };
        struct range UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct select {
            int i;
        };
        struct select UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct type {
            int i;
        };
        struct type UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct var {
            int i;
        };
        struct var UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct _ {
            int i;
        };
        struct _ UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct init {
            int i;
        };
        struct init UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct len {
            int i;
        };
        struct len UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct copy {
            int i;
        };
        struct copy UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct fmt {
            int i;
        };
        struct fmt UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }
    {
        struct cap {
            int i;
        };
        struct cap UU;
        UU.i = 5;
        is_eq(UU.i, 5);
    }

    // uncomment after success implementation of struct scope
    // https://github.com/Konstantin8105/c4go/issues/368
    /*
	diag("Typedef struct name from Go keyword")
	{ typedef struct {int i;} chan        ;	chan        UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} defer       ;	defer       UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} fallthrough ;	fallthrough UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} func        ;	func        UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} go          ;	go          UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} import      ;	import      UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} interface   ;	interface   UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} map         ;	map         UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} package     ;	package     UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} range       ;	range       UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} select      ;	select      UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} type        ;	type        UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} var         ;	var         UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} _           ;	_           UU; UU.i = 5; is_eq(UU.i,5);}
	{ typedef struct {int i;} init        ;	init        UU; UU.i = 5; is_eq(UU.i,5);}
*/

    struct_with_rec_fuction();

    diag("name of struct inside struct");
    {
        typedef struct TI TI;
        struct TI {
            TI *left, *right;
            double varTI;
        };
        TI t1;
        t1.varTI = 4.3;
        TI t2;
        t2.varTI = 4.1;
        TI tt;
        tt.left = &t1;
        (*tt.left).right = &t2;
        tt.right = &t2;
        is_eq((*tt.left).varTI, 4.3);
        is_eq((*(*tt.left).right).varTI, 4.1);
        is_eq((*tt.right).varTI, 4.1);
    }

    struct_in_struct_with_star();
    struct_null();

    func_in_func_in_struct();

    struct_inside_union();

    done_testing();
}
