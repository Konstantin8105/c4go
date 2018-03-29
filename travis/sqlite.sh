#!/bin/bash

set -e

# These steps are from the README to verify it can be installed and run as
# documented.
go build

# This will have to be updated every so often to the latest version. You can
# find the latest version here: https://sqlite.org/download.html
export SQLITE3_FILE=sqlite-amalgamation-3220000

# Variable for location of temp sqlite files
SQLITE_TEMP_FOLDER="/tmp/SQLITE"
mkdir -p $SQLITE_TEMP_FOLDER

# Download/unpack SQLite if required.
if [ ! -e $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip ]; then
    curl http://sqlite.org/2018/$SQLITE3_FILE.zip > $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip
    unzip $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip -d $SQLITE_TEMP_FOLDER
fi

# Clean generated files. This should not be required, but it's polite.
rm -f $SQLITE_TEMP_FOLDER/sqlite3.go $SQLITE_TEMP_FOLDER/shell.go

# Transpile the SQLite3 files.
echo "Transpiling shell.c..."
./c4go transpile -o=$SQLITE_TEMP_FOLDER/shell.go   $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/shell.c  

# sqlite3.c
echo "Transpiling sqlite3.c..."
./c4go transpile -o=$SQLITE_TEMP_FOLDER/sqlite3.go $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/sqlite3.c

# Show amount "Warning" in sqlite Go codes
SQLITE_WARNINGS=`cat $SQLITE_TEMP_FOLDER/sqlite3.go $SQLITE_TEMP_FOLDER/shell.go | grep "^// Warning" | sort | uniq | wc -l`
echo "In files (sqlite3.go and shell.go) summary : $SQLITE_WARNINGS warnings."

# SQLITE
c4go transpile -o="$SQLITE_TEMP_FOLDER/sqlite.go" -clang-flag="-DSQLITE_THREADSAFE=0" -clang-flag="-DSQLITE_OMIT_LOAD_EXTENSION" $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/shell.c $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/sqlite3.c

# Show amount "Warning":
SQLITE_WARNINGS=`cat $SQLITE_TEMP_FOLDER/sqlite.go | grep "^// Warning" | sort | uniq | wc -l`
echo "After transpiling shell.c and sqlite3.c together, have summary: $SQLITE_WARNINGS warnings."

# Show amount error from `go build`:
SQLITE_WARNINGS_GO=`go build $SQLITE_TEMP_FOLDER/sqlite.go 2>&1 | wc -l`
echo "In file sqlite.go summary : $SQLITE_WARNINGS_GO warnings in go build."

# Amount warning from gometalinter
echo "Calculation warnings by gometalinter"
GOMETALINTER_WARNINGS=`$GOPATH/bin/gometalinter $SQLITE_TEMP_FOLDER/sqlite.go 2>&1 | wc -l`
echo "Amount found warnings by gometalinter at 30 second : $GOMETALINTER_WARNINGS warnings."
