#!/bin/bash

set -e

go build

mkdir -p ./testdata/

# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export GIT_SOURCE="https://github.com/aligrudi/neatroff.git"
	export NAME="neatroff"
	export TEMP_FOLDER="./testdata/$NAME"
	export GO_FILE="$TEMP_FOLDER/$NAME.go"
	export GO_APP="$TEMP_FOLDER/$NAME.app"
	export COMMIT=08d6e9ceb2d7d420be1d916d32faeda18c4f27e0

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		git clone $GIT_SOURCE $TEMP_FOLDER
		cd $TEMP_FOLDER/
		git checkout $COMMIT
		cd ../../
		sed -i.bak '566,572d' $TEMP_FOLDER/font.c
		sed -i.bak '66,81d'   $TEMP_FOLDER/out.c
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app

# Arguments menu
echo "    -a for ast tree"
if [ "$1" == "-a" ]; then
	$C4GO ast                               \
		-clang-flag="-DTROFFFDIR=\"MMM\""	\
		-clang-flag="-DTROFFMDIR=\"WWW\""	\
		$TEMP_FOLDER/*.c
fi


# transpilation of all projects
	echo "Transpile to $GO_FILE"
	$C4GO transpile                         \
		-s									\
		-clang-flag="-DTROFFFDIR=\"MMM\""	\
		-clang-flag="-DTROFFMDIR=\"WWW\""	\
		-o="$GO_FILE"                       \
		$TEMP_FOLDER/*.c

# show warnings comments in Go source
	echo "Calculate warnings in file: $GO_FILE"
	WARNINGS=`cat $GO_FILE | grep "^// Warning" | sort | uniq | wc -l`
	echo "		After transpiling : $WARNINGS warnings."

# show other warnings
	# show amount error from `go build`:
		echo "Build to $GO_APP file"
		WARNINGS_GO=`go build -o $GO_APP -gcflags="-e" $GO_FILE 2>&1 | wc -l`
		echo "		Go build : $WARNINGS_GO warnings"
	# amount unsafe
		UNSAFE=`cat $GO_FILE | grep "unsafe\." | wc -l`
		echo "		Unsafe   : $UNSAFE"
	# amount Go code lines
		LINES=`wc $GO_FILE`
		echo "(lines,words,bytes)	 : $LINES"


# Arguments menu
echo "    -s for show detail of Go build errors"
if [ "$1" == "-s" ]; then
	# show go build warnings	
		# c4go warnings
			cat $GO_FILE | grep "^// Warning" | sort | uniq
		# show amount error from `go build`:
			go build -o $GO_APP -gcflags="-e" $GO_FILE 2>&1
fi

