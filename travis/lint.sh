#!/bin/bash

set -e

mkdir -p ./testdata/

file="./testdata/gofmt.list"
eval "find ./ -name '*.go' | grep -v 'testdata' | grep -v 'vendor' > $file"

while IFS= read -r line
do
	# Check go fmt first
	if [ -n "$(gofmt -l $line)" ]; then
		echo "Go code is not properly formatted. Use 'gofmt' for: $line"
		gofmt -d .
		exit 1
	fi
done < "$file"
