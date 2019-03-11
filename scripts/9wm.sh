#!/bin/bash

set -e

go build

# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export TEMP_FOLDER="/tmp/9wm"

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		git clone https://github.com/9wm/9wm.git $TEMP_FOLDER/
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app

# transpilation 
export FILE="$TEMP_FOLDER/9wm.go"
export FILES=`ls $TEMP_FOLDER/*.c | tr "\n" " "`
$C4GO transpile  -s                                          \
	             -o="$FILE"                                  \
				 $FILES

# show warnings comments in Go source
	echo "	***** warnings"
	WARNINGS=`cat $FILE | grep "^// Warning" | sort | uniq | wc -l`
	echo "		After transpiling : $WARNINGS warnings."
# show amount error from `go build`:
	WARNINGS_GO=`go build -o $TEMP_FOLDER/9wm.app -gcflags="-e" $FILE 2>&1 | wc -l`
	echo "		Go build : $WARNINGS_GO warnings"
# amount unsafe
	UNSAFE=`cat $FILE | grep unsafe | wc -l`
	echo "		Unsafe   : $UNSAFE"


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
				go build -o $f.app -gcflags="-e" $f 2>&1 | sort 
		done
fi
