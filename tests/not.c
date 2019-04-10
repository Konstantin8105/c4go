#include "tests.h"
#include <stdio.h>

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void print_bool(int tr, int fl)
{
    if (tr) {
        printf("print_bool tr is true\n");
    } else {
        printf("print_bool tr is false\n");
    }
    if (fl) {
        printf("print_bool fl is true\n");
    } else {
        printf("print_bool fl is false\n");
    }
    if (tr != 0 && fl != 0) {
        pass("test != 0");
    }
    if (tr == 0 && fl == 0) {
        pass("test == 0");
    }
}

void ok()
{
    pass("success");
}

void f()
{
    fail("need fix");
}

#define not_c_type(type)                                      \
    diag("------------------");                               \
    {                                                         \
        diag(#type);                                          \
        type a;                                               \
        int p;                                                \
        diag("not for C type : zero value");                  \
        a = 0;                                                \
        if (!a) {                                             \
            printf("a1\n");                                   \
            ok();                                             \
        } else {                                              \
            printf("a2\n");                                   \
            f();                                              \
        }                                                     \
        diag("not-not for C type : zero value");              \
        a = 0;                                                \
        if (!(!a)) {                                          \
            printf("a3\n");                                   \
            f();                                              \
        } else {                                              \
            printf("a4\n");                                   \
            ok();                                             \
        }                                                     \
        diag("for C type : function");                        \
        print_bool(!a, !(!a));                                \
        diag("for C type : assign");                          \
        p = !a;                                               \
        if (p) {                                              \
            printf("p1\n");                                   \
        } else {                                              \
            printf("p2\n");                                   \
        }                                                     \
        diag("not for C type : non-zero positive value");     \
        a = 42;                                               \
        if (!a) {                                             \
            printf("a5\n");                                   \
        } else {                                              \
            printf("a6\n");                                   \
        }                                                     \
        diag("for C type : function");                        \
        print_bool(!a, !(!a));                                \
        diag("for C type : assign");                          \
        p = !a;                                               \
        if (p) {                                              \
            printf("p1\n");                                   \
        } else {                                              \
            printf("p2\n");                                   \
        }                                                     \
        diag("not-not for C type : non-zero positive value"); \
        a = 42;                                               \
        if (!(!a)) {                                          \
            printf("a7\n");                                   \
        } else {                                              \
            printf("a8\n");                                   \
        }                                                     \
        diag("not for C type : non-zero negative value");     \
        a = -42;                                              \
        if (!a) {                                             \
            printf("a5\n");                                   \
        } else {                                              \
            printf("a6\n");                                   \
        }                                                     \
        diag("not-not for C type : non-zero negative value"); \
        a = -42;                                              \
        if (!(!a)) {                                          \
            printf("a7\n");                                   \
        } else {                                              \
            printf("a8\n");                                   \
        }                                                     \
        diag("for C type : function");                        \
        print_bool(!a, !(!a));                                \
        diag("for C type : assign");                          \
        p = !a;                                               \
        if (p) {                                              \
            printf("p1\n");                                   \
        } else {                                              \
            printf("p2\n");                                   \
        }                                                     \
        diag("for C type : zero value");                      \
        a = 0;                                                \
        if (a) {                                              \
            printf("a9\n");                                   \
        } else {                                              \
            printf("a10\n");                                  \
        }                                                     \
        diag("for C type : positive value");                  \
        a = 15;                                               \
        if (a) {                                              \
            printf("a11\n");                                  \
        } else {                                              \
            printf("a12\n");                                  \
        }                                                     \
        diag("for C type : negative value");                  \
        a = -15;                                              \
        if (a) {                                              \
            printf("a13\n");                                  \
        } else {                                              \
            printf("a14\n");                                  \
        }                                                     \
        diag("for C type : assign");                          \
        p = !a;                                               \
        if (p) {                                              \
            printf("p1\n");                                   \
        } else {                                              \
            printf("p2\n");                                   \
        }                                                     \
    }

void test_c_types()
{
    not_c_type(char);
    not_c_type(double);
    not_c_type(float);
    not_c_type(short);
    not_c_type(int);
    not_c_type(long);
    not_c_type(long double);
    not_c_type(long long);
    not_c_type(signed char);
}

#define not_c_pointer(type)                  \
    diag("------------------");              \
    {                                        \
        diag(#type);                         \
        type* a = NULL;                      \
        int p;                               \
        diag("for C pointer:  null");        \
        if (a) {                             \
            printf("a1\n");                  \
        } else {                             \
            printf("a2\n");                  \
        }                                    \
        diag("for C pointer: not null");     \
        if (!a) {                            \
            printf("a2\n");                  \
        } else {                             \
            printf("a3\n");                  \
        }                                    \
        diag("for C pointer: not-not null"); \
        if (!(!a)) {                         \
            printf("a4\n");                  \
        } else {                             \
            printf("a5\n");                  \
        }                                    \
        diag("for C type : function");       \
        print_bool(!a, !(!a));               \
        p = !a;                              \
        if (p) {                             \
            printf("p1\n");                  \
        } else {                             \
            printf("p2\n");                  \
        }                                    \
        type b = 42;                         \
        a = &b;                              \
        diag("for C pointer:  null");        \
        if (a) {                             \
            printf("a11\n");                 \
        } else {                             \
            printf("a12\n");                 \
        }                                    \
        diag("for C pointer: not null");     \
        if (!a) {                            \
            printf("a12\n");                 \
        } else {                             \
            printf("a13\n");                 \
        }                                    \
        diag("for C pointer: not-not null"); \
        if (!(!a)) {                         \
            printf("a14\n");                 \
        } else {                             \
            printf("a15\n");                 \
        }                                    \
        diag("for C type : function");       \
        print_bool(!a, !(!a));               \
        p = !a;                              \
        if (p) {                             \
            printf("p1\n");                  \
        } else {                             \
            printf("p2\n");                  \
        }                                    \
    }

#define not_c_struct(type)                   \
    diag("------------------");              \
    {                                        \
        diag(#type);                         \
        type* a = NULL;                      \
        int p;                               \
        diag("for C pointer:  null");        \
        if (a) {                             \
            printf("a1\n");                  \
        } else {                             \
            printf("a2\n");                  \
        }                                    \
        diag("for C pointer: not null");     \
        if (!a) {                            \
            printf("a2\n");                  \
        } else {                             \
            printf("a3\n");                  \
        }                                    \
        diag("for C pointer: not-not null"); \
        if (!(!a)) {                         \
            printf("a4\n");                  \
        } else {                             \
            printf("a5\n");                  \
        }                                    \
        diag("for C type : function");       \
        print_bool(!a, !(!a));               \
        p = !a;                              \
        if (p) {                             \
            printf("p1\n");                  \
        } else {                             \
            printf("p2\n");                  \
        }                                    \
        type b;                              \
        a = &b;                              \
        diag("for C pointer:  null");        \
        if (a) {                             \
            printf("a11\n");                 \
        } else {                             \
            printf("a12\n");                 \
        }                                    \
        diag("for C pointer: not null");     \
        if (!a) {                            \
            printf("a12\n");                 \
        } else {                             \
            printf("a13\n");                 \
        }                                    \
        diag("for C pointer: not-not null"); \
        if (!(!a)) {                         \
            printf("a14\n");                 \
        } else {                             \
            printf("a15\n");                 \
        }                                    \
        diag("for C type : function");       \
        print_bool(!a, !(!a));               \
        p = !a;                              \
        if (p) {                             \
            printf("p1\n");                  \
        } else {                             \
            printf("p2\n");                  \
        }                                    \
    }

void test_c_pointers()
{
    not_c_pointer(char);
    not_c_pointer(double);
    not_c_pointer(float);
    not_c_pointer(int);
    not_c_pointer(long double);
    not_c_pointer(long long);
    not_c_pointer(signed char);
    not_c_pointer(unsigned long);
}

struct str {
    int i;
};
union un {
    int i;
    double d;
};

void test_c_struct()
{
    not_c_struct(struct str);
    not_c_struct(union un);
}

void test_c_function()
{
    void (*a)(void);
    a = NULL;
    int p;
    diag("for C pointer:  null");
    if (a) {
        printf("a1\n");
    } else {
        printf("a2\n");
    }
    diag("for C pointer: not null");
    if (!a) {
        printf("a2\n");
    } else {
        printf("a3\n");
    }
    diag("for C pointer: not-not null");
    if (!(!a)) {
        printf("a4\n");
    } else {
        printf("a5\n");
    }
    diag("for C type : function");
    print_bool(!a, !(!a));
    p = !a;
    if (p) {
        printf("p1\n");
    } else {
        printf("p2\n");
    }
    a = test_c_pointers;
    diag("for C pointer:  null");
    if (a) {
        printf("a11\n");
    } else {
        printf("a12\n");
    }
    diag("for C pointer: not null");
    if (!a) {
        printf("a12\n");
    } else {
        printf("a13\n");
    }
    diag("for C pointer: not-not null");
    if (!(!a)) {
        printf("a14\n");
    } else {
        printf("a15\n");
    }
    diag("for C type : function");
    print_bool(!a, !(!a));
    p = !a;
    if (p) {
        printf("p1\n");
    } else {
        printf("p2\n");
    }
}

int main()
{
    plan(18);

    START_TEST(c_types);
    START_TEST(c_pointers);
    START_TEST(c_struct);
    START_TEST(c_function);

    done_testing();
}
