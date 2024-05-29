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

go test                                 \
				  -cover                \
				  -timeout=30m          \
                  -covermode=atomic     \
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
go test                                       \
	-run=TestIntegrationScripts/tests/ctype.c \
	-race -v
