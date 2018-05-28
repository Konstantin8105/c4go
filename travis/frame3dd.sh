#!/bin/bash

set -e

# These steps are from the README to verify it can be installed and run as
# documented.
go build

export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
export C4GO=$C4GO_DIR/c4go

FRAME3DD_TEMP_FOLDER="/tmp/FRAME3DD"
CODE_PATH="$FRAME3DD_TEMP_FOLDER/src"

if [ ! -d $FRAME3DD_TEMP_FOLDER ]; then
	mkdir -p -v $FRAME3DD_TEMP_FOLDER
	git clone -b Debug2 https://github.com/Konstantin8105/History_frame3DD.git $FRAME3DD_TEMP_FOLDER
fi

c4go transpile -o="$CODE_PATH/main.go" -clang-flag="-I$CODE_PATH/viewer/" -clang-flag="-I$CODE_PATH/microstran/" "$CODE_PATH/main.c" "$CODE_PATH/frame3dd.c" "$CODE_PATH/frame3dd_io.c" "$CODE_PATH/coordtrans.c" "$CODE_PATH/eig.c" "$CODE_PATH/HPGmatrix.c" "$CODE_PATH/HPGutil.c" "$CODE_PATH/NRutil.c"

# Show amount "Warning":
FRAME3DD_WARNINGS=`cat $CODE_PATH/main.go | grep "^// Warning" | sort | uniq | wc -l`
echo "After transpiling fraqme3dd, have summary: $FRAME3DD_WARNINGS warnings."

# Show amount error from `go build`:
FRAME3DD_GO=`go build $CODE_PATH/main.go 2>&1 | wc -l`
echo "In file sqlite.go summary : $FRAME3DD_GO warnings in go build."

# Amount warning from gometalinter
echo "Calculation warnings by gometalinter"
GOMETALINTER_WARNINGS=`$GOPATH/bin/gometalinter $CODE_PATH/main.go 2>&1 | wc -l`
echo "Amount found warnings by gometalinter at 30 second : $GOMETALINTER_WARNINGS warnings."

SQLITE_UNSAFE=`cat $CODE_PATH/main.go | grep unsafe | wc -l`
echo "Amount unsafe package using: $SQLITE_UNSAFE"
