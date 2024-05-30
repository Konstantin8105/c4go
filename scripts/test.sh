#!/bin/bash

# go install github.com/ory/go-acc@latest
#
# touch ./coverage.tmp
# echo 'mode: atomic' > coverage.txt
# go list ./... | grep -v /cmd | grep -v /vendor | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp -coverpkg $(go list ./... | grep -v /vendor | tr "\n" ",") {} && tail -n +2 coverage.tmp >> coverage.txt || exit 255' && rm coverage.tmp


# set -e
# 
# echo "" > coverage.txt
# 
# mkdir -p ./testdata/

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
# export PKGS="github.com/Konstantin8105/c4go github.com/Konstantin8105/c4go/ast github.com/Konstantin8105/c4go/noarch github.com/Konstantin8105/c4go/preprocessor github.com/Konstantin8105/c4go/program github.com/Konstantin8105/c4go/scripts github.com/Konstantin8105/c4go/transpiler github.com/Konstantin8105/c4go/types github.com/Konstantin8105/c4go/util github.com/Konstantin8105/c4go/version"
export PKGS='github.com/Konstantin8105/c4go github.com/Konstantin8105/c4go/ast github.com/Konstantin8105/c4go/version'
# export PKGS=$(go list -e ./... | grep -v testdata | grep -v examples | grep -v tests | grep -v vendor | tr '\n' ' ')

# Make comma-separated.
export PKGS_DELIM=$(echo "$PKGS" | tr ' ' ',')

echo "PKGS       : $PKGS"
echo "PKGS_DELIM : $PKGS_DELIM"

touch ./coverage.tmp
echo 'mode: atomic' > coverage.txt
go list ./...  | grep -v testdata | grep -v examples | grep -v tests | grep -v vendor | grep -v /cmd | grep -v /vendor | xargs -n100 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp -coverpkg $(go list ./... | grep -v /vendor | tr "\n" ",") {} && tail -n +2 coverage.tmp >> coverage.txt || exit 255' && rm coverage.tmp
 
# go test \
# 	-cover                \
# 	-covermode=atomic     \
# 	-timeout=30m          \
# 	-coverpkg=$PKGS_DELIM \
# 	-coverprofile=./testdata/pkg.coverprofile $PKGS

# Merge coverage profiles.
# COVERAGE_FILES=`ls -1 ./testdata/*.coverprofile 2>/dev/null | wc -l`
# if [ $COVERAGE_FILES != 0 ]; then
# 	# check program `gocovmerge` is exist
# 	if which gocovmerge >/dev/null 2>&1; then
# 		export FILES=$(ls testdata/*.coverprofile | tr '\n' ' ')
# 		echo "Combine next coverprofiles : $FILES"
# 		gocovmerge $FILES > coverage.txt
# 	fi
# fi

echo "End of coverage"
