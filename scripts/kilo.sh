#!/bin/bash

set -e

go build

mkdir -p ./testdata/

# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export GIT_SOURCE="https://github.com/antirez/kilo.git"
	export NAME="kilo"
	export TEMP_FOLDER="./testdata/$NAME"
	export GO_FILE="$TEMP_FOLDER/$NAME.go"
	export GO_APP="$TEMP_FOLDER/$NAME.app"
	export COMMIT=d65f4c92e8ed405937a7bac3248d24fa6b40eb6f

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		git clone $GIT_SOURCE $TEMP_FOLDER
		cd $TEMP_FOLDER/
		git checkout $COMMIT
		cd ../../
		sed -i.bak '370iif(row->rsize > 0)' $TEMP_FOLDER/kilo.c
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app

# transpilation of all projects
	echo "Transpile to $GO_FILE"
	$C4GO transpile                         \
		-o="$GO_FILE"                       \
		$TEMP_FOLDER/kilo.c

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

# debugging
if [ "$1" == "-d" ]; then
	# try to run
		echo " move to folder"
			cd ./testdata/kilo/
		echo "step 1: create debug file"
			$C4GO debug kilo.c
		echo "step 2: prepare output data"
			echo "" > output.txt
			echo "" > output.g.txt
			echo "" > output.c.txt
		echo "step 3: prepare test script"
			echo -e 'Hello my dear friend\x0D\x13\x11' > script.txt
		echo "step 4: run Go application"
			$C4GO transpile -o=debug.kilo.go debug.kilo.c
			go build -o kilo.go.app	debug.kilo.go
			echo "" > debug.txt
			cat script.txt | ./kilo.go.app output.txt 2>&1 && echo "ok" || echo "not ok"
			cp output.txt output.g.txt
			echo "" > output.txt
			cp debug.txt debug.go.txt
			echo ""
			echo ""
			echo ""
			echo ""
			echo ""
			echo ""
			echo ""
			echo ""
			echo ""
			echo ""
			echo ""
			echo ""
			echo ""
			echo ""
		echo "step 5: run C application"
			clang -o kilo.c.app debug.kilo.c 2>&1
			echo "" > debug.txt
			cat script.txt | ./kilo.c.app output.txt  2>&1 && echo "ok" || echo "not ok"
			cp output.txt output.c.txt
			echo "" > output.txt
			cp debug.txt debug.c.txt
		echo "step 5"
			echo "-----------------------------"
			echo "debug"
			diff -y -t debug.c.txt debug.go.txt 2>&1  > debug.diff  && echo "ok" || echo "not ok"
			# cat debug.diff
			echo "-----------------------------"
			echo "output"
			diff -y -t output.c.txt output.g.txt 2>&1 > output.diff && echo "ok" || echo "not ok"
			cat output.diff
			echo "-----------------------------"
		echo "step 6: move back"
			cd ../../
fi

# Arguments menu
echo "    -s for show detail of Go build errors"
if [ "$1" == "-s" ]; then
	# show go build warnings	
		# c4go warnings
			cat $GO_FILE | grep "^// Warning" | sort | uniq
		# show amount error from `go build`:
			go build -o $GO_APP -gcflags="-e" $GO_FILE 2>&1
fi
