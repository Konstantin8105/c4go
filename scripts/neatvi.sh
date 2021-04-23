#!/bin/bash

set -e

go build

mkdir -p ./testdata/

# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export GIT_SOURCE="https://github.com/aligrudi/neatvi.git"
	export NAME="neatvi"
	export TEMP_FOLDER="./testdata/$NAME"
	export GO_FILE="$TEMP_FOLDER/$NAME.go"
	export GO_APP="$TEMP_FOLDER/$NAME.app"
	export COMMIT=29b4fa028ebed5c2b2a9a7e13bce07ab2ae81587

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		git clone $GIT_SOURCE $TEMP_FOLDER
		cd $TEMP_FOLDER/
		git checkout $COMMIT
		cd ../../
		sed -i.bak '635,642d'   $TEMP_FOLDER/vi.c
		sed -i.bak 's/bufs/bufs_postfix/g'   $TEMP_FOLDER/reg.c
		sed -i.bak 's/static\ //g'   $TEMP_FOLDER/regex.c
		sed -i.bak '53,59d;9,24d'   $TEMP_FOLDER/uc.c
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app

# transpilation of all projects
	echo "Transpile to $GO_FILE"
	$C4GO transpile                         \
		-s                                  \
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
	# defers
		DEFER=`cat $GO_FILE | grep "defer func" | wc -l`
		echo "defer func           	 : $DEFER"


# Arguments menu
echo "    -s for show detail of Go build errors"
if [ "$1" == "-s" ]; then
	# show go build warnings	
		# c4go warnings
			cat $GO_FILE | grep "^// Warning" | sort | uniq
		# show amount error from `go build`:
			go build -o $GO_APP -gcflags="-e" $GO_FILE 2>&1
fi

