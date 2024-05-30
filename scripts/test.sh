#!/bin/bash

set -e

# Initialize
mkdir -p ./testdata/

# github.com/Konstantin8105/c4go
# github.com/Konstantin8105/c4go/ast
# github.com/Konstantin8105/c4go/examples // ignore
# github.com/Konstantin8105/c4go/noarch
# github.com/Konstantin8105/c4go/preprocessor
# github.com/Konstantin8105/c4go/program
# github.com/Konstantin8105/c4go/scripts
# github.com/Konstantin8105/c4go/testdata // ignore
# github.com/Konstantin8105/c4go/tests    // ignore
# github.com/Konstantin8105/c4go/transpiler
# github.com/Konstantin8105/c4go/types
# github.com/Konstantin8105/c4go/util
# github.com/Konstantin8105/c4go/version

# Package list
export PKGS=$(go list -e ./... | grep -v testdata | grep -v examples | grep -v tests | grep -v vendor | tr '\n' ' ')
# export PKGS="github.com/Konstantin8105/c4go github.com/Konstantin8105/c4go/util"

# View
echo "PKGS       : $PKGS"

# Initialize
touch ./coverage.tmp

# Run tests
echo 'mode: atomic' > coverage.txt
echo "$PKGS" | xargs -n100 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp -coverpkg $(go list ./... | grep -v /vendor | tr "\n" ",") {} && tail -n +2 coverage.tmp >> coverage.txt || exit 255' && rm coverage.tmp
 
# Finilize
echo "End of coverage"
