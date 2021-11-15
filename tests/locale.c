#include "tests.h"
#include <locale.h>

int main()
{
    plan(0);// 2);
	// TODO: commented for github action

   // setlocale(LC_MONETARY, "");
   // struct lconv* lc;
   // lc = localeconv();
   // // Local Currency Symbol
   // is_true(strlen(lc->currency_symbol) > 0);
   // // International Currency Symbol
   // is_true(strlen(lc->int_curr_symbol) > 0);

    done_testing();
}
