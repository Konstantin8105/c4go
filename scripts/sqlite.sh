#!/bin/bash

set -e

# These steps are from the README to verify it can be installed and run as
# documented.
go build

mkdir -p ./testdata/

export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
export C4GO=$C4GO_DIR/c4go

# This will have to be updated every so often to the latest version. You can
# find the latest version here: https://sqlite.org/download.html
export SQLITE3_FILE=sqlite-amalgamation-3250200

# Variable for location of temp sqlite files
SQLITE_TEMP_FOLDER="./testdata/SQLITE"
mkdir -p $SQLITE_TEMP_FOLDER

# Download/unpack SQLite if required.
if [ ! -e $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip ]; then
    curl http://sqlite.org/2018/$SQLITE3_FILE.zip > $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip
    unzip $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip -d $SQLITE_TEMP_FOLDER
fi

# SQLITE
$C4GO transpile  -s                                          \
                 -cpuprofile=./testdata/cpu.out              \
	             -o="$SQLITE_TEMP_FOLDER/sqlite.go"          \
				 -clang-flag="-DSQLITE_THREADSAFE=0"         \
				 -clang-flag="-DSQLITE_OMIT_LOAD_EXTENSION"  \
				 $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/shell.c   \
				 $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/sqlite3.c

# See profile file
# Run:
# go tool pprof ./testdata/cpu.out

# Show amount "Warning":
export GO_FILE="$SQLITE_TEMP_FOLDER/sqlite.go"
echo "Calculate warnings in file: $GO_FILE"

SQLITE_WARNINGS=`cat $GO_FILE | grep "^// Warning" | sort | uniq | wc -l`
echo "		After transpiling : $SQLITE_WARNINGS warnings."

# Show amount error from `go build`:
SQLITE_WARNINGS_GO=`go build -o $SQLITE_TEMP_FOLDER/sqlite.app $SQLITE_TEMP_FOLDER/sqlite.go -gcflags="-e" 2>&1 | wc -l`
echo "		Go build : $SQLITE_WARNINGS_GO warnings"

SQLITE_UNSAFE=`cat $SQLITE_TEMP_FOLDER/sqlite.go | grep "unsafe\." | wc -l`
echo "		Unsafe   : $SQLITE_UNSAFE"

# amount Go code lines
	LINES=`wc $SQLITE_TEMP_FOLDER/sqlite.go`
	echo "(lines,words,bytes)	 : $LINES"
# defers
	DEFER=`cat $SQLITE_TEMP_FOLDER/sqlite.go| grep "defer func" | wc -l`
	echo "defer func           	 : $DEFER"

# Arguments menu
echo "    -s for show detail of Go build errors"
if [ "$1" == "-s" ]; then
	# show go build warnings	
		# c4go warnings
			cat $SQLITE_TEMP_FOLDER/sqlite.go | grep "^// Warning" | sort | uniq
		# show amount error from `go build`:
			go build -o $SQLITE_TEMP_FOLDER/sqlite.app -gcflags="-e"  $SQLITE_TEMP_FOLDER/sqlite.go 2>&1
fi
