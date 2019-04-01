#!/bin/bash

set -e

echo "" > coverage.txt

mkdir -p ./testdata/

# Package list
export PKGS=$(go list ./... | grep -v c4go/testdata | grep -v c4go/examples | grep -v c4go/tests | grep -v /vendor/ | tr '\n' ' ')

# Make comma-separated.
export PKGS_DELIM=$(echo "$PKGS" | tr ' ' ',')

echo "PKGS       : $PKGS"
echo "PKGS_DELIM : $PKGS_DELIM"

go test -v -cover -timeout=30m          \
				  -tags integration     \
	              -coverpkg=$PKGS_DELIM \
				  -coverprofile=./testdata/pkg.coverprofile $PKGS

# Merge coverage profiles.
COVERAGE_FILES=`ls -1 ./testdata/*.coverprofile 2>/dev/null | wc -l`
if [ $COVERAGE_FILES != 0 ]; then
	# check program `gocovmerge` is exist
	if which gocovmerge >/dev/null 2>&1; then
		export FILES=$(ls testdata/*.coverprofile | tr '\n' ' ')
		echo "Combine next coverprofiles : $FILES"
		gocovmerge $FILES > coverage.txt
	fi
fi

echo "End of coverage"

# check race
go test -tags=integration -run=TestIntegrationScripts/tests/ctype.c -race -v

# These steps are from the README to verify it can be installed and run as
# documented.
go build

export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
export C4GO=$C4GO_DIR/c4go

echo "Run: c4go transpile prime.c"
$C4GO transpile -o=./testdata/prime.go $C4GO_DIR/examples/prime.c
echo "47" | go run ./testdata/prime.go
if [ $(cat ./testdata/prime.go | wc -l) -eq 0 ]; then exit 1; fi
if [ $($C4GO ast $C4GO_DIR/examples/prime.c | wc -l) -eq 0 ]; then exit 1; fi

echo "----------------------"

# Run script sqlite
source ./scripts/sqlite.sh
