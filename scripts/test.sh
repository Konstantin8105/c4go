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
				  -tags integration     \
	              -coverpkg=$PKGS_DELIM \
				  -coverprofile=coverage.txt $PKGS

# check race
go test -tags=integration                     \
	-run=TestIntegrationScripts/tests/ctype.c \
	-race -v

echo "End of coverage"
