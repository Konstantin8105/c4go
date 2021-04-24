#!/bin/bash

set -e

go build

mkdir -p ./testdata/

# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export GIT_SOURCE="https://github.com/aligrudi/neatpost.git"
	export NAME="neatpost"
	export TEMP_FOLDER="./testdata/$NAME"
	export GO_FILE="$TEMP_FOLDER/$NAME.go"
	export GO_APP="$TEMP_FOLDER/$NAME.app"
	export COMMIT=393011f64e853e3e2dec6950aab8c90ef12832b9

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		git clone $GIT_SOURCE $TEMP_FOLDER
		cd $TEMP_FOLDER/
		git checkout $COMMIT
		echo "sed steps"
		# run 
		# ./scripts/neatpost.sh -s 2>&1| grep "ps.c" | grep redef
  		sed -i.bak '92,98   s/^/\/\/ /' font.c
 		sed -i.bak '110,122 s/^/\/\/ /' ps.c
 		sed -i.bak '124,127 s/^/\/\/ /' ps.c
 		sed -i.bak '129,132 s/^/\/\/ /' ps.c
 		sed -i.bak '134,138 s/^/\/\/ /' ps.c
 		sed -i.bak '140,144 s/^/\/\/ /' ps.c
 		sed -i.bak '147,152 s/^/\/\/ /' ps.c
 		sed -i.bak '153,157 s/^/\/\/ /' ps.c
 		sed -i.bak '159,162 s/^/\/\/ /' ps.c
 		sed -i.bak '164,174 s/^/\/\/ /' ps.c
 		sed -i.bak '187,193 s/^/\/\/ /' ps.c
 		sed -i.bak '196,367 s/^/\/\/ /' ps.c
 		sed -i.bak '34,41   s/^/\/\/ /' ps.c
 		sed -i.bak '369,375 s/^/\/\/ /' ps.c
 		sed -i.bak '377,397 s/^/\/\/ /' ps.c
 		sed -i.bak '43,46   s/^/\/\/ /' ps.c
 		sed -i.bak '48,57   s/^/\/\/ /' ps.c
 		sed -i.bak '500,524 s/^/\/\/ /' ps.c
 		sed -i.bak '59,79   s/^/\/\/ /' ps.c
 		sed -i.bak '8       s/^/\/\/ /' ps.c
 		sed -i.bak '11,13   s/^/\/\/ /' ps.c
 		sed -i.bak '176,177 s/^/\/\/ /' ps.c
 		sed -i.bak '82,89   s/^/\/\/ /' ps.c
 		sed -i.bak '91,108  s/^/\/\/ /' ps.c
		cd ../../
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
