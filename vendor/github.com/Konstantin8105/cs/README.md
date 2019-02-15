# cs
my golang code style

```
PACKAGE DOCUMENTATION

package cs
    import "."


FUNCTIONS

func All(t *testing.T)
    All run all codestyle test for golang sources

	Ignore data from folder "testdata"

func Debug(t *testing.T)
    Debug test source for avoid debug printing

	Ignore data from folder "testdata"

func Os(t *testing.T)
    Os test source for avoid words "darwin", "macos"

	Ignore data from folder "testdata"

func Todo(t *testing.T)
    Todo calculate amount comments with TODO, FIX, BUG in golang sources.

	Ignore data from folder "testdata"

```
