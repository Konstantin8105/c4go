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

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		git clone $GIT_SOURCE $TEMP_FOLDER
		sed -i.bak '370iif(row->rsize > 0)' $TEMP_FOLDER/kilo.c
		# add debugging output
		sed -i.bak '243s/(1)/(243)/'     $TEMP_FOLDER/kilo.c
		sed -i.bak '787s/(1)/(787)/'     $TEMP_FOLDER/kilo.c
		sed -i.bak '1251s/(1)/(1251)/'   $TEMP_FOLDER/kilo.c
		sed -i.bak '1259s/(1)/(1259)/'   $TEMP_FOLDER/kilo.c
		# add debugging file
		sed -i.bak '53ivoid debug(char *msg){FILE*file;file=fopen("./debug.txt","a");if(file==NULL){exit(53);};fprintf(file,"%s\n",msg);fclose(file);}'  $TEMP_FOLDER/kilo.c
		# add debug information
		sed -i.bak '726i{char buffer[500];sprintf(buffer,"line726: filerow=%d",filerow);debug(buffer);}' $TEMP_FOLDER/kilo.c
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app

# transpilation of all projects
	echo "Transpile to $GO_FILE"
	$C4GO transpile                         \
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


# Arguments menu
echo "    -s for show detail of Go build errors"
if [ "$1" == "-s" ]; then
	# show go build warnings	
		# c4go warnings
			cat $GO_FILE | grep "^// Warning" | sort | uniq
		# show amount error from `go build`:
			go build -o $GO_APP -gcflags="-e" $GO_FILE 2>&1
fi
