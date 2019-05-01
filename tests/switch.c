// This file tests the various forms of switch statement.
//
// We must be extra sensitive to the fact that switch fallthrough is handled
// differently in C and Go. Break statements are removed and fallthrough
// statements are added when nessesary.
//
// It is worth mentioning that a SwitchStmt has a CompoundStmt item that
// contains all of the cases. However, if the individual case are not enclosed
// in a scope each of the case statements and their childen are part of the same
// CompoundStmt. For example, the first switch statement below contains a
// CompoundStmt with 12 children.

#include "tests.h"
#include <stdio.h>

void match_a_single_case()
{
    switch (1) {
    case 0:
        fail("code should not reach here");
        break;
    case 1:
        pass(__func__);
        break;
    case 2:
        fail("code should not reach here");
        break;
    default:
        fail("code should not reach here");
        break;
    }
}

void fallthrough_to_next_case()
{
    switch (1) {
    case 0:
        fail("code should not reach here");
        break;
    case 1:
        pass(__func__);
    case 2:
        pass(__func__);
        break;
    default:
        fail("code should not reach here");
        break;
    }
}

void match_no_cases()
{
    switch (1) {
    case 5:
        fail("code should not reach here");
        break;
    case 2:
        fail("code should not reach here");
        break;
    }
}

void match_default()
{
    switch (1) {
    case 5:
        fail("code should not reach here");
        break;
    case 2:
        fail("code should not reach here");
        break;
    default:
        pass(__func__);
        break;
    }
}

void fallthrough_several_cases_including_default()
{
    switch (1) {
    case 0:
        fail("code should not reach here");
    case 1:
        pass(__func__);
    case 2:
        pass(__func__);
    default:
        pass(__func__);
    }
}

void scoped_match_a_single_case()
{
    switch (1) {
    case 0: {
        fail("code should not reach here");
        break;
    }
    case 1: {
        pass(__func__);
        break;
    }
    case 2: {
        fail("code should not reach here");
        break;
    }
    default: {
        fail("code should not reach here");
        break;
    }
    }
}

void scoped_fallthrough_to_next_case()
{
    switch (1) {
    case 0: {
        fail("code should not reach here");
        break;
    }
    case 1: {
        pass(__func__);
    }
    case 2: {
        pass(__func__);
        break;
    }
    default: {
        fail("code should not reach here");
        break;
    }
    }
}

void scoped_match_no_cases()
{
    switch (1) {
    case 5: {
        fail("code should not reach here");
        break;
    }
    case 2: {
        fail("code should not reach here");
        break;
    }
    }
}

void scoped_match_default()
{
    switch (1) {
    case 5: {
        fail("code should not reach here");
        break;
    }
    case 2: {
        fail("code should not reach here");
        break;
    }
    default: {
        pass(__func__);
        break;
    }
    }
}

void scoped_fallthrough_several_cases_including_default()
{
    switch (1) {
    case 0: {
        fail("code should not reach here");
    }
    case 1: {
        pass(__func__);
    }
    case 2: {
        pass(__func__);
    }
    default: {
        pass(__func__);
    }
    }
}

typedef struct I67 I67;
struct I67 {
    int x, y;
};

void run(I67* i, int pos)
{
    switch (pos) {
    case 0:
        (*i).x += 1;
        (*i).y += 1;
        break;
    case 1:
        (*i).x -= 1;
        (*i).y -= 1;
        break;
    }
}

void run_with_block(I67* i, int pos)
{
    switch (pos) {
    case 0: {
        (*i).x += 1;
        (*i).y += 1;
        break;
    }
    case 1: {
        (*i).x -= 1;
        (*i).y -= 1;
    } break;
    case 2:
        (*i).x *= 1;
        (*i).y *= 1;
        break;
    default:
        (*i).x *= 10;
        (*i).y *= 10;
    }
}

void switch_issue67()
{
    I67 i;
    i.x = 0;
    i.y = 0;
    run(&i, 0);
    is_eq(i.x, 1);
    is_eq(i.y, 1);
    run(&i, 1);
    is_eq(i.x, 0);
    is_eq(i.y, 0);
    run_with_block(&i, 0);
    is_eq(i.x, 1);
    is_eq(i.y, 1);
    run_with_block(&i, 1);
    is_eq(i.x, 0);
    is_eq(i.y, 0);
}

void empty_switch()
{
    int pos = 0;
    switch (pos) {
    }
    is_eq(pos, 0);
}

void default_only_switch()
{
    int pos = 0;
    switch (pos) {
    case -1: // empty case
    case -1 - 4: // empty case
    case (-1 - 4 - 4): // empty case
    case (-3): // empty case
    case -2: // empty case
    default:
        pos++;
    }
    is_eq(pos, 1);
}

void switch_without_input()
{
    int pos = 0;
    switch (0) {
    default:
        pos++;
    }
    is_eq(pos, 1);
}

void test_switch()
{
    int val = 1;
    switch (0) {
    case -1000:
    default:
        // ignored
        val += 3;
    case -1:
        val += 1;
    case 2: {
        val *= 2;
    }
    case 3: {
        val *= 4;
    }
    case 4:
    case 5:
        val += 5;
    case 6:
    case 7:
        val *= 6;
        break;
    case 8:;
    case 9:;
    case 10:;
    }
    is_eq(val, 270);
}

void switch_bool()
{
    double x = 3.0;
    switch (x < 0.0) {
    case 0:
        pass("ok");
        break;
    case 1:
        fail("switch_bool 1");
        break;
    default:
        fail("switch_bool 3");
    }
    int y = 7;
    const int t = 0;
    switch (y) {
    case 0 == t:
        fail("switch_bool 4");
        break;
    case 1 == t:
        fail("switch_bool 5");
        break;
    default:
        pass("ok");
        break;
    }
}

void switch_char()
{
    char form = '1';
    switch (form) {
    default:
    case '1':
        form = 'S';
        break;
    case 0:;
    };
    is_eq(form, 'S');
}

int value()
{
    return 42;
}

void switch_stat()
{
    int v = 15;
    switch (v = value()) {
    case 15:
        fail("wrong");
        break;
    case 42:
        is_true(v == 42);
        break;
    }
}

int case_inside_block(int a)
{
    switch (a) {
    case 1:
        break;
    case 2: {
        a = 42;
    case 3:
        a += 45;
        break;
    }
    case 4:
        break;
    }
    return a;
}

int main()
{
    plan(37);

    switch_char();
    switch_bool();
    match_a_single_case();
    fallthrough_to_next_case();
    match_no_cases();
    match_default();
    fallthrough_several_cases_including_default();
    test_switch();

    // For each of the tests above there will be identical cases that use scopes
    // for the case statements.
    scoped_match_a_single_case();
    scoped_fallthrough_to_next_case();
    scoped_match_no_cases();
    scoped_match_default();
    scoped_fallthrough_several_cases_including_default();

    switch_issue67();
    empty_switch();
    default_only_switch();
    switch_without_input();
    switch_stat();

    is_eq(case_inside_block(0), 0);
    is_eq(case_inside_block(1), 1);
    is_eq(case_inside_block(2), 45 + 42);
    is_eq(case_inside_block(3), 3 + 45);
    is_eq(case_inside_block(4), 4);
    is_eq(case_inside_block(5), 5);

	diag("without parens");
	{
		int x = 0;
		switch (x) 
			case 0:
				pass("ok");
		(void)(x);
	}

    done_testing();
}
