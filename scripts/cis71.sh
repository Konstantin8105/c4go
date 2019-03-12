#!/bin/bash

set -e

go build

# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export TEMP_FOLDER="/tmp/CIS71"

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/hello.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/power2.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/homework1.c			-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/add2.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/addn.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/add.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/coins.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/factorial.c			-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/true.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/fibo.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/funcs.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/funcs2.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/scope1.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/scope2.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/array.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/array1.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/array2.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/misc.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/addresses.c			-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/codes.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/random.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/randompermute.c		-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/line.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/linear.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/shift.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/sieve.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/string1.c			-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/counts.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/cpfile.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/enum1.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/enum2.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/binary.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/selection.c			-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/bubble.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/number.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/cpintarray.c			-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/struct.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/studentarray.c		-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/merge.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/sortmerge.c			-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/clean.c				-P $TEMP_FOLDER
		wget --no-check-certificate https://cis.temple.edu/~giorgio/cis71/code/studentlist.c		-P $TEMP_FOLDER
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app

# transpilation each file
	COUNTER=0
	for f in $TEMP_FOLDER/*.c; do
				let COUNTER=COUNTER+1
			# iteration by C projects
				echo "***** transpilation folder : $f"
        	# Will not run if no directories are available
				$C4GO transpile 						\
					-o="$f.go" $f
			# show warnings comments in Go source
				export FILE="$f.go"
				echo "	***** warnings"
				WARNINGS=`cat $FILE | grep "^// Warning" | sort | uniq | wc -l`
				echo "		After transpiling : $WARNINGS warnings."
			# show amount error from `go build`:
				WARNINGS_GO=`go build -o $TEMP_FOLDER/$COUNTER.app -gcflags="-e" $FILE 2>&1 | wc -l`
				echo "		Go build : $WARNINGS_GO warnings"
			# amount unsafe
				UNSAFE=`cat $FILE | grep unsafe | wc -l`
				echo "		Unsafe   : $UNSAFE"
	done
