#!/bin/bash

set -e

# These steps are from the README to verify it can be installed and run as
# documented.
go build

export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
export C4GO=$C4GO_DIR/c4go

export GSL_SOURCE="http://mirror.tochlab.net/pub/gnu/gsl"
export GSL_FILE="gsl-2.4"


# Variable for location of temp sqlite files
GSL_TEMP_FOLDER="/tmp/GSL"
mkdir -p $GSL_TEMP_FOLDER

# Download/unpack SQLite if required.
if [ ! -e $GSL_TEMP_FOLDER/$GSL_FILE.tar.gz ]; then
    curl "$GSL_SOURCE/$GSL_FILE.tar.gz" > "$GSL_TEMP_FOLDER/$GSL_FILE.tar.gz"
	tar -C "$GSL_TEMP_FOLDER" -xzf "$GSL_TEMP_FOLDER/$GSL_FILE.tar.gz"
fi

cd $GSL_TEMP_FOLDER/$GSL_FILE/
chmod u+x configure
echo "" > /tmp/gcc.log
./configure CC=gccecho
make
go run ./gcc.log.go
