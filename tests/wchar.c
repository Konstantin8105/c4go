#include "tests.h"
#include <wchar.h>

int main() {
	plan(3);

	{
		diag("wcscpy");
		wchar_t wcs1[]=L"Sample string";
		wchar_t wcs2[40];
		wchar_t wcs3[40];
		wcscpy (wcs2,wcs1);
		wcscpy (wcs3,L"copy successful");
		is_streq(wcs1,L"Sample string");
		is_streq(wcs2,L"Sample string");
		is_streq(wcs3,L"copy successful");
	}

	done_testing();
}
