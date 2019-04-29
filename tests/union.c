// Tests for unions.

#include "tests.h"
#include <stdio.h>

union programming {
    int constant;
    char* pointer;
};

union programming init_var()
{
    union programming variable;
    char* s = "Programming in Software Development.";

    variable.pointer = s;
    is_streq(variable.pointer, "Programming in Software Development.");

    variable.constant = 123;
    is_eq(variable.constant, 123);

    return variable;
}

void pass_by_ref(union programming* addr)
{
    char* s = "Show string member.";
    int v = 123 + 456;

    addr->constant += 456;
    is_eq(addr->constant, v);

    addr->pointer = s;
    is_streq(addr->pointer, "Show string member.");
}

void var_by_val(union programming value)
{
    value.constant++;

    is_eq(value.constant, 124);
}

struct SHA3 {
    union {
        double iY;
        double dY;
    } uY;
    float ffY;
};

union unknown {
    double i2;
    double d2;
};
struct SHA32 {
    union unknown u2;
    float ff2;
};

void union_inside_struct()
{
    diag("Union inside struct");
    struct SHA3 sha;
    sha.ffY = 12.444;
    sha.uY.iY = 4;
    is_eq(sha.uY.iY, 4);
    is_eq(sha.uY.dY, 4);
    is_eq(sha.ffY, 12.444);

    struct SHA32 sha2;
    sha2.ff2 = 12.444;
    sha2.u2.i2 = 4;
    is_eq(sha2.u2.i2, 4);
    is_eq(sha2.u2.d2, 4);
    is_eq(sha2.ff2, 12.444);
    pass("ok");

    union unknown data[2];
    for (int i = 0; i < 2; i++) {
        data[i].i2 = i + 3;
        is_true(data[i].d2 != 0);
    }
    (void)(data);
}

typedef union myunion myunion;
typedef union myunion {
    double PI;
    int B;
} MYUNION;

typedef union {
    double PI;
    int B;
} MYUNION2;

void union_typedef()
{
    diag("Typedef union");
    union myunion m;
    double v = 3.14;
    m.PI = v;
    is_eq(m.PI, 3.14);
    is_true(m.B != 0);
    is_eq(v, 3.14);
    v += 1.0;
    is_eq(v, 4.14);
    is_eq(m.PI, 3.14);

    MYUNION mm;
    mm.PI = 3.14;
    is_eq(mm.PI, 3.14);
    is_true(mm.B != 0);

    myunion mmm;
    mmm.PI = 3.14;
    is_eq(mmm.PI, 3.14);
    is_true(mmm.B != 0);

    MYUNION2 mmmm;
    mmmm.PI = 3.14;
    is_eq(mmmm.PI, 3.14);
    is_true(mmmm.B != 0);
}

typedef struct FuncDestructor FuncDestructor;
struct FuncDestructor {
    int i;
};
typedef struct FuncDef FuncDef;
struct FuncDef {
    int i;
    union {
        FuncDef* pHash;
        FuncDestructor* pDestructor;
    } u;
};

void union_inside_struct2()
{
    FuncDef f;
    FuncDestructor fd;
    fd.i = 100;
    f.u.pDestructor = &fd;

    FuncDestructor* p_fd = f.u.pDestructor;
    is_eq((*p_fd).i, 100);

    is_true(f.u.pHash != NULL);
    is_true(f.u.pDestructor != NULL);
    int vHash = (*f.u.pHash).i;
    is_eq(vHash, 100);
    is_eq((*f.u.pHash).i, 100);
}

union UPNT {
    int* a;
    int* b;
    int* c;
};

void union_pointers()
{
    union UPNT u;
    int v = 32;
    u.a = &v;
    is_eq(*u.a, 32);
    is_eq(*u.b, 32);
    is_eq(*u.c, 32);
    pass("ok")
}

union UPNTF {
    int (*f1)(int);
    int (*f2)(int);
};

int union_function(int a)
{
    return a + 1;
}

void union_func_pointers()
{
    union UPNTF u;
    u.f1 = union_function;
    is_eq(u.f1(21), 22);
    is_eq(u.f2(21), 22);
}

union array_union {
    float a[2];
    float b[2];
};

void union_array()
{
    union array_union arr;
    arr.a[0] = 12;
    arr.b[1] = 14;
    is_eq(arr.a[0], 12);
    is_eq(arr.a[1], 14);
    is_eq(arr.b[0], 12);
    is_eq(arr.b[1], 14);
}

typedef int ii;
typedef struct SHA SHA;
struct SHA {
    union {
        ii s[25];
        unsigned char x[100];
    } u;
    unsigned uuu;
};

void union_arr_in_str()
{
    SHA sha;
    sha.uuu = 15;
    is_eq(sha.uuu, 15);
    for (int i = 0; i < 25; i++)
        sha.u.s[0] = 0;
    is_eq(sha.u.s[0], 0);
    is_true(sha.u.x[0] == 0);
    for (int i = 0; i < 6; i++) {
        sha.u.s[i] = (ii)(4);
        sha.u.s[i] = (ii)(42) + sha.u.s[i];
    }
    is_eq(sha.u.s[5], 46);
    is_true(sha.u.x[0] != 0);
}

union un_struct {
    struct {
        short a;
        short b;
    } str;
    long l;
};

void union_with_struct()
{
    union un_struct u;
    u.str.a = 12;
    u.str.b = 45;
    is_eq(u.str.a, 12);
    is_eq(u.str.b, 45);
    is_true(u.l > 0);
}

struct suf {
    union {
        int* i;
        void (*sa)(int, double*, void*);
    } uf;
};

void union_with_func()
{
    struct suf s;
    int a = 42;
    s.uf.i = &a;
    is_not_null(s.uf.i);
    is_not_null(s.uf.sa);
    (void)(s.uf);
}


union MyNumber {
    int n;
    char s[200];
} obj;


union MyNumber getNumber(char x, int state)
{
    union MyNumber tmp;
    if (state)
        tmp.n = (int)(x + 10 - 'A');
    else {
        switch (x) {
            case 'A':
                strcpy(tmp.s, "десять");
                break;
            case 'B':
                strcpy(tmp.s, "одиннадцать");
                break;
            case 'C':
                strcpy(tmp.s, "двенадцать");
                break;
            case 'D':
                strcpy(tmp.s, "тринадцать");
                break;
            case 'E':
                strcpy(tmp.s, "четырнадцать");
                break;
            case 'F':
                strcpy(tmp.s, "пятнадцать");
        }
    }
    return tmp;
}


void test_union_function()
{
    char k;
    for (k = 'A'; k <= 'F'; k++)
        printf("%c - %d\n", k, getNumber(k, 1).n);
    for (k = 'A'; k <= 'F'; k++) {
        obj = getNumber(k, 0);
        printf("%c - %s\n", k, obj.s);
    }
}

int main()
{
    plan(50);

    union programming variable;

    variable = init_var();
    var_by_val(variable);
    pass_by_ref(&variable);

    union_inside_struct();
    union_typedef();
    union_inside_struct2();
    union_pointers();
    union_func_pointers();
    union_array();
    union_arr_in_str();
    union_with_struct();
    union_with_func();
	test_union_function();

    done_testing();
}
