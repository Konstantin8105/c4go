#!/bin/bash

set -e

# These steps are from the README to verify it can be installed and run as
# documented.
go build

export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
export C4GO=$C4GO_DIR/c4go

# Enviroment variable
export TEMP_FOLDER="/tmp/Deitel"

# Delete folder if exist
if [ -d "$TEMP_FOLDER" ]; then rm -Rf $TEMP_FOLDER; fi

# Variable for location of temp sqlite files
mkdir -p $TEMP_FOLDER

# Get all sources
git clone https://github.com/Emmetttt/C-Deitel-Book.git $TEMP_FOLDER

# List of all C files
FILE_LIST="$(find $TEMP_FOLDER -name "*.c"| \
	grep -v ".exe" | \
	grep -v ".md" )"

while read -r fname; do
    echo $fname
    $C4GO transpile -o="$fname.go" "$fname"
	WARNINGS=`cat "$fname".go | grep "^// Warning" | sort | uniq | wc -l`
	echo "In file $fname : $WARNINGS warnings."
done <<< "$FILE_LIST"


# List of all C files
FILE_LIST="$(find $TEMP_FOLDER -name "*.C"| \
	grep -v "E5.45.C" | \
	grep -v "06.14_const_type_qualifier.C" | \
	grep -v "E7.17.C" | \
	grep -v ".exe" | \
	grep -v ".md" )"

while read -r fname; do
    echo $fname
    $C4GO transpile -o="$fname.go" "$fname"
	WARNINGS=`cat "$fname".go | grep "^// Warning" | sort | uniq | wc -l`
	echo "In file $fname : $WARNINGS warnings."
done <<< "$FILE_LIST"
