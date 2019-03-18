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
SQLITE_TEMP_FOLDER="/tmp/SQLITE"
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
SQLITE_WARNINGS=`cat $SQLITE_TEMP_FOLDER/sqlite.go | grep "^// Warning" | sort | uniq | wc -l`
echo "After transpiling shell.c and sqlite3.c together, have summary: $SQLITE_WARNINGS warnings."

# Show amount error from `go build`:
SQLITE_WARNINGS_GO=`go build $SQLITE_TEMP_FOLDER/sqlite.go -gcflags="-e" 2>&1 | wc -l`
echo "In file sqlite.go summary : $SQLITE_WARNINGS_GO warnings in go build."

SQLITE_UNSAFE=`cat $SQLITE_TEMP_FOLDER/sqlite.go | grep unsafe | wc -l`
echo "Amount unsafe package using: $SQLITE_UNSAFE"
