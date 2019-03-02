#!/bin/bash

set -e

go build

# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export NEAT_TEMP_FOLDER="/tmp/NEAT"

# prepare C code
    if [ ! -d $NEAT_TEMP_FOLDER ]; then
		mkdir -p $NEAT_TEMP_FOLDER
		git clone https://github.com/aligrudi/neateqn.git	$NEAT_TEMP_FOLDER/0
		sed '90,103d' /tmp/NEAT/0/box.c > /tmp/NEAT/0/box.c
		git clone https://github.com/aligrudi/neatrefer.git	$NEAT_TEMP_FOLDER/1
		git clone https://github.com/aligrudi/neatpost.git	$NEAT_TEMP_FOLDER/2
		sed '19,25d' 	/tmp/NEAT/2/dev.c 		> /tmp/NEAT/2/dev.c
		sed '24d'    	/tmp/NEAT/2/post.c		> /tmp/NEAT/2/post.c
		sed '754,756d' 	/tmp/NEAT/2/pdf.c		> /tmp/NEAT/2/pdf.c
		sed '338,343d'	/tmp/NEAT/2/pdf.c		> /tmp/NEAT/2/pdf.c
		git clone https://github.com/aligrudi/neatmkfn.git	$NEAT_TEMP_FOLDER/3
		sed '20,24d'	/tmp/NEAT/3/afm.c		> /tmp/NEAT/3/afm.c
		git clone https://github.com/aligrudi/neatroff.git	$NEAT_TEMP_FOLDER/4
		sed '25,31d'	/tmp/NEAT/4/dev.c		> /tmp/NEAT/4/dev.c
		sed '142,150d'	/tmp/NEAT/4/draw.c		> /tmp/NEAT/4/draw.c
		sed '135,140d'	/tmp/NEAT/4/draw.c		> /tmp/NEAT/4/draw.c
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $NEAT_TEMP_FOLDER/*.go
	rm -f $NEAT_TEMP_FOLDER/*.app

# transpilation of all projects
	COUNTER=0
	for f in $NEAT_TEMP_FOLDER/*; do
    	if [ -d ${f} ]; then
			# iteration by C projects
				echo "***** transpilation folder : $COUNTER"
        	# Will not run if no directories are available
				$C4GO transpile -s						\
					-clang-flag="-DTROFFFDIR=\"MMM\""	\
					-clang-flag="-DTROFFMDIR=\"WWW\""	\
					-o="$NEAT_TEMP_FOLDER/$COUNTER.go" $f/*.c
				let COUNTER=COUNTER+1
    	fi
	done

# show warnings comments in Go source
	echo "***** warnings"
	WARNINGS=`cat $NEAT_TEMP_FOLDER/*.go | grep "^// Warning" | sort | uniq | wc -l`
	echo "After transpiling : $WARNINGS warnings."

# show other warnings
for f in $NEAT_TEMP_FOLDER/*.go ; do
	# iteration by Go files
		echo "	file : $f"
	# show amount error from `go build`:
		WARNINGS_GO=`go build -o $f.app -gcflags="-e" $f 2>&1 | wc -l`
		echo "		Go build : $WARNINGS_GO warnings"
	# amount unsafe
		UNSAFE=`cat $f | grep unsafe | wc -l`
		echo "		Unsafe   : $UNSAFE"
done


# Arguments menu
echo "    -s for show detail of Go build errors"
if [ "$1" == "-s" ]; then
	# show go build warnings	
		for f in $NEAT_TEMP_FOLDER/*.go ; do
			# iteration by Go files
				echo "	file : $f"
			# show amount error from `go build`:
				go build -o $f.app -gcflags="-e" $f 2>&1 | sort 
		done
fi
