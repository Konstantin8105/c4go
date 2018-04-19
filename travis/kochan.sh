#!/bin/bash

set -e

# These steps are from the README to verify it can be installed and run as
# documented.
go build

export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
export C4GO=$C4GO_DIR/c4go

# Enviroment variable
export TEMP_FOLDER="/tmp/Kochan"

# Delete folder if exist
if [ -d "$TEMP_FOLDER" ]; then rm -Rf $TEMP_FOLDER; fi

# Variable for location of temp sqlite files
mkdir -p $TEMP_FOLDER

# Get all sources
git clone https://github.com/eugenetriguba/programming-in-c.git $TEMP_FOLDER

# List of all C files
FILE_LIST="$(find $TEMP_FOLDER -name "*.c" | \
	grep -v "5.9d.c" | \
	grep -v "5.9c.c" )"

while read -r fname; do
    echo $fname
    $C4GO transpile -o="$fname.go" "$fname"
	WARNINGS=`cat "$fname".go | grep "^// Warning" | sort | uniq | wc -l`
	echo "In file $fname : $WARNINGS warnings."
done <<< "$FILE_LIST"
