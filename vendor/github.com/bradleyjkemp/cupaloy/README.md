<h1 align="center">
    <img src="https://github.com/bradleyjkemp/cupaloy/blob/master/mascot.png" alt="Mascot" width="200">
    <br>
    <a href="https://travis-ci.org/bradleyjkemp/cupaloy"><img src="https://travis-ci.org/bradleyjkemp/cupaloy.svg?branch=master" alt="Build Status" /></a>
    <a href="https://coveralls.io/github/bradleyjkemp/cupaloy?branch=master"><img src="https://coveralls.io/repos/github/bradleyjkemp/cupaloy/badge.svg" alt="Coverage Status" /></a>
    <a href="https://goreportcard.com/report/github.com/bradleyjkemp/cupaloy"><img src="https://goreportcard.com/badge/github.com/bradleyjkemp/cupaloy" alt="Go Report Card" /></a>
    <a href="https://godoc.org/github.com/bradleyjkemp/cupaloy"><img src="https://godoc.org/github.com/bradleyjkemp/cupaloy?status.svg" alt="GoDoc" /></a>
    <a href="https://sourcegraph.com/github.com/bradleyjkemp/cupaloy?badge"><img src="https://sourcegraph.com/github.com/bradleyjkemp/cupaloy/-/badge.svg" alt="Number of users" /></a>
</h1>

Incredibly simple Go snapshot testing: `cupaloy` takes a snapshot of your test output and compares it to a snapshot committed alongside your tests. If the values don't match then the test will be failed.

There's no need to manually manage snapshot files: just use the `cupaloy.SnapshotT(t, value)` function in your tests and `cupaloy` will automatically find the relevant snapshot file (based on the test name) and compare it with the given value.

## Usage
### Write a test
Firstly, write a test case generating some output and pass this output to `cupaloy.SnapshotT`:
```golang
func TestParsing(t *testing.T) {
    ast := ParseFile("test_input")

    // check that the result is the same as the last time the snapshot was updated
    // if the result has changed (e.g. because the behaviour of the parser has changed)
    // then the test will be failed with an error containing a diff of the changes
    cupaloy.SnapshotT(t, ast)
}
```
The first time this test is run, a snapshot will be automatically created (using the [github.com/davecgh/go-spew](https://github.com/davecgh/go-spew) package).

### Update a snapshot
When the behaviour of your software changes causing the snapshot to change, this test will begin to fail with an error showing the difference between the old and new snapshots. Once you are happy that the new snapshot is correct (and hasn't just changed unexpectedly), you can save the new snapshot by setting the ```UPDATE_SNAPSHOTS``` environment and re-running your tests:
```bash
UPDATE_SNAPSHOTS=true go test ./...
```
This will fail all tests where the snapshot was updated (to stop you accidentally updating snapshots in CI) but your snapshot files will now have been updated to reflect the current output of your code.

### Supported formats
Snapshots of test output are generated using the [github.com/davecgh/go-spew](https://github.com/davecgh/go-spew) package which uses reflection to deep pretty-print your test result and so will support almost all the basic types (from simple strings, slices, and maps to deeply nested structs) without issue. The only types whose contents cannot be fully pretty-printed are functions and channels.

The most important property of your test output is that it is deterministic: if your output contains timestamps or other fields which will change on every run, then `cupaloy` will detect this as a change and so fail the test.


### Further Examples
#### Table driven tests
```golang
var testCases = map[string][]string{
    "TestCaseOne": []string{......},
    "AnotherTestCase": []string{......},
    ....
}

func TestCases(t *testing.T) {
    for testName, args := range testCases {
        t.Run(testName, func(t *testing.T) {
            result := functionUnderTest(args...)
            cupaloy.SnapshotT(t, result)
        })
    }
}
```
#### Changing output directory
```golang
func TestSubdirectory(t *testing.T) {
    result := someFunction()
    snapshotter := cupaloy.New(cupaloy.SnapshotSubdirectory("testdata"))
    err := snapshotter.Snapshot(result)
    if err != nil {
        t.Fatalf("error: %s", err)
    }
}
```
For further usage examples see basic_test.go and advanced_test.go in the examples/ directory which are both kept up to date and run on CI.
