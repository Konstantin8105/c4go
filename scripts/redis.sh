#!/bin/bash

set -e

go build

mkdir -p ./testdata/

# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export ARCHIVE="http://download.redis.io/releases/redis-5.0.4.tar.gz"
	export NAME="redis"
	export VERSION="redis-5.0.4"
	export TEMP_FOLDER="./testdata/$NAME"
	export GO_FILE="$TEMP_FOLDER/$NAME.go"
	export GO_APP="$TEMP_FOLDER/$NAME.app"

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		curl $ARCHIVE > $TEMP_FOLDER/$NAME.tar.gz
		tar -xf $TEMP_FOLDER/$NAME.tar.gz -C $TEMP_FOLDER/
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app

# transpilation of all projects
	echo "Transpile to $GO_FILE"
	$C4GO transpile                                          \
		-s                                                   \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/hiredis"   \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/jemalloc"  \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/linenoise" \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/lua"       \
		-o="$GO_FILE"                       \
		$TEMP_FOLDER/$VERSION/src/redis-cli.c

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


# Arguments menu
echo "    -s for show detail of Go build errors"
if [ "$1" == "-s" ]; then
	# show go build warnings	
		# c4go warnings
			cat $GO_FILE | grep "^// Warning" | sort | uniq
		# show amount error from `go build`:
			go build -o $GO_APP -gcflags="-e" $GO_FILE 2>&1
fi




