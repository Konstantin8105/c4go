#!/bin/bash

set -e

echo "" > coverage.txt

# The code below was copied from:
# https://github.com/golang/go/issues/6909#issuecomment-232878416
#
# As in @rodrigocorsi2 comment above (using full path to grep due to 'grep -n'
# alias).
export PKGS=$(go list ./... | grep -v c4go/build | grep -v /vendor/)

# Make comma-separated.
export PKGS_DELIM=$(echo "$PKGS" | paste -sd "," -)

# Run tests and append all output to out.txt. It's important we have "-v" so
# that all the test names are printed. It's also important that the covermode be
# set to "count" so that the coverage profiles can be merged correctly together
# with gocovmerge.
#
# Exit code 123 will be returned if any of the tests fail.
go list -f 'go test -v -tags integration -race -covermode atomic -coverprofile {{.Name}}.coverprofile -coverpkg $PKGS_DELIM {{.ImportPath}}' $PKGS | xargs -I{} bash -c "{}"

# Merge coverage profiles.
COVERAGE_FILES=`ls -1 *.coverprofile 2>/dev/null | wc -l`
if [ $COVERAGE_FILES != 0 ]; then
	# check program `gocovmerge` is exist
	if which gocovmerge >/dev/null 2>&1; then
		gocovmerge `ls *.coverprofile` > coverage.txt
		rm *.coverprofile
	fi
fi

# These steps are from the README to verify it can be installed and run as
# documented.
go build

export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
export C4GO=$C4GO_DIR/c4go

echo "Run: c4go transpile prime.c"
$C4GO transpile -o=/tmp/prime.go $C4GO_DIR/examples/prime.c
echo "47" | go run /tmp/prime.go
if [ $(cat /tmp/prime.go | wc -l) -eq 0 ]; then exit 1; fi
if [ $($C4GO ast $C4GO_DIR/examples/prime.c | wc -l) -eq 0 ]; then exit 1; fi

echo "----------------------"

# Run script sqlite
source ./travis/sqlite.sh

# Run script triangle
source ./travis/triangle.sh

# Run script Deitel
source ./travis/deitel.sh
