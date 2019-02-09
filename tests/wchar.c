#include "tests.h"
#include <wchar.h>

int main() {
	plan(2);

	{
		diag("wcscmp");
		wchar_t wsKey[] = L"jilia";
		if (wcscmp(wsKey,L"jilia") == 0){
			pass("0k - equal");
		}
		if (wcscmp(wsKey,L"jiiia") != 0){
			pass("0k - not equal");
		}
	}

	// {
	// 	diag("wcscpy");
	// 	wchar_t wcs1[]=L"Sample string";
	// 	wchar_t wcs2[40];
	// 	wchar_t wcs3[40];
	// 	wcscpy (wcs2,wcs1);
	// 	wcscpy (wcs3,L"copy successful");
	// 	is_wchareq(wcs1,L"Sample string");
	// 	is_wchareq(wcs2,L"Sample string");
	// 	is_wchareq(wcs3,L"copy successful");
	// }

	done_testing();
}
