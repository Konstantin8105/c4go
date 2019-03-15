#!/bin/bash

set -e

go build

# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export GIT_SOURCE="https://github.com/aligrudi/neatpost.git"
	export NAME="neatpost"
	export TEMP_FOLDER="/tmp/$NAME"
	export GO_FILE="$TEMP_FOLDER/$NAME.go"
	export GO_APP="$TEMP_FOLDER/$NAME.app"

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		git clone $GIT_SOURCE $TEMP_FOLDER
		sed -i.bak '92,98d' $TEMP_FOLDER/font.c
		sed -i.bak '163,173d;158,161d;152,156d;145,150d;139,143d;133,137d;128,131d;123,127d;109,121d;90,107d;80,88d;58,78d;47,56d;42,45d;33,40d;8d' $TEMP_FOLDER/ps.c
		sed -i.bak '219,232d;212,217d;204,210d;196,202d;179,194d;175,177d;131,173d' $TEMP_FOLDER/ps.c
		sed -i.bak '238,262d;123,129d;116,121d;109,114d;102,107d;95,100d;82,93d;73,80d;65,71d;57,63d' $TEMP_FOLDER/ps.c
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app

# transpilation of all projects
	echo "Transpile to $GO_FILE"
	$C4GO transpile                         \
		-clang-flag="-DTROFFFDIR=\"MMM\""	\
		-clang-flag="-DTROFFMDIR=\"WWW\""	\
		-o="$GO_FILE"                       \
		$TEMP_FOLDER/*.c

# show warnings comments in Go source
	echo "Calculate warnings in file: $GO_FILE"
	WARNINGS=`cat $GO_FILE | grep "^// Warning" | sort | uniq | wc -l`
	echo "After transpiling : $WARNINGS warnings."

# show other warnings
	# show amount error from `go build`:
		echo "Build to $GO_APP file"
		WARNINGS_GO=`go build -o $GO_APP -gcflags="-e" $GO_FILE 2>&1 | wc -l`
		echo "		Go build : $WARNINGS_GO warnings"
	# amount unsafe
		UNSAFE=`cat $GO_FILE | grep unsafe | wc -l`
		echo "		Unsafe   : $UNSAFE"


# Arguments menu
echo "    -s for show detail of Go build errors"
if [ "$1" == "-s" ]; then
	# show go build warnings	
		# c4go warnings
			cat $GO_FILE | grep "^// Warning" | sort | uniq
		# show amount error from `go build`:
			go build -o $GO_APP -gcflags="-e" $GO_FILE 2>&1 | sort 
fi
