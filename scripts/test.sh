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

# echo "" > coverage.out
# for d in $(go list ./... | grep -v vendor); do
#     go test -v -race -coverprofile=profile.out -covermode=atomic $d
#     if [ -f profile.out ]; then
#         cat profile.out >> coverage.out
#         rm profile.out
#     fi
# done

echo "End of coverage"
