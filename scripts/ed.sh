#!/bin/bash

set -e

go build

mkdir -p ./testdata/

# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export TEMP_FOLDER="./testdata/ed"
	export VERSION="ed-1.15"

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		mkdir -p $TEMP_FOLDER/$VERSION
		wget --no-check-certificate http://mirror.tochlab.net/pub/gnu/ed/$VERSION.tar.lz -P $TEMP_FOLDER
		echo "Please don't forget install: sudo apt-get install lzip"
		tar --lzip  -C $TEMP_FOLDER/$VERSION -xvf $TEMP_FOLDER/$VERSION.tar.lz
		sed -i.bak -e '92d'                                        $TEMP_FOLDER/$VERSION/$VERSION/signal.c
		sed -i.bak -e '92i   if( ioctl(0, TIOCGWINSZ, &ws) >= 0 )' $TEMP_FOLDER/$VERSION/$VERSION/signal.c
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app

# tranpilation
$C4GO transpile  -cpuprofile=./testdata/cpu.out              \
				 -s 										 \
	             -o="$TEMP_FOLDER/$VERSION.go"               \
				 -clang-flag="-DPROGVERSION=\"$VERSION\""    \
				 $TEMP_FOLDER/$VERSION/$VERSION/*.c

echo "Calculate warnings : $TEMP_FOLDER"
# show warnings comments in Go source
	export FILE="$TEMP_FOLDER/$VERSION.go"
	WARNINGS=`cat $FILE | grep "^// Warning" | sort | uniq | wc -l`
	echo "		After transpiling : $WARNINGS warnings."
# show amount error from `go build`:
	WARNINGS_GO=`go build -o $TEMP_FOLDER/$COUNTER.app -gcflags="-e" $FILE 2>&1 | wc -l`
	echo "		Go build : $WARNINGS_GO warnings"
# amount unsafe
	UNSAFE=`cat $FILE | grep "unsafe\." | wc -l`
	echo "		Unsafe   : $UNSAFE"
# amount Go code lines
	LINES=`wc $TEMP_FOLDER/$VERSION.go`
	echo "(lines,words,bytes)	 : $LINES"
	# defers
		DEFER=`cat $TEMP_FOLDER/$VERSION.go | grep "defer func" | wc -l`
		echo "defer func           	 : $DEFER"

# Arguments menu
echo "    -s for show detail of Go build errors"
if [ "$1" == "-s" ]; then
	# show go build warnings	
		for f in $TEMP_FOLDER/*.go ; do
			# iteration by Go files
				echo "	file : $f"
			# c4go warnings
				cat $f | grep "^// Warning" | sort | uniq
			# show amount error from `go build`:
				go build -o $f.app -gcflags="-e" $f 2>&1
		done
fi
