#!/bin/bash

# See that example from https://github.com/codecov/example-go
set -e
echo "" > coverage.txt

PKGS="./ast/... ./linux/... ./noarch/... ./preprocessor/... ./program/... ./transpiler/... ./types/... ./util/... ./version/..."

echo "List of Go packages"
echo "$PKGS"
go test -v -tags=integration -coverprofile=profile.out -covermode=atomic  $PKGS
cat profile.out >> coverage.txt
rm profile.out

# check race
go test -tags=integration -run=TestIntegrationScripts/tests/ctype.c -race -v

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
