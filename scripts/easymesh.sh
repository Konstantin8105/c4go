#!/bin/bash

set -e

go build

mkdir -p ./testdata/

# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export NAME="easymesh"
	export TEMP_FOLDER="./testdata/$NAME"
	export GO_FILE="$TEMP_FOLDER/$NAME.go"
	export GO_APP="$TEMP_FOLDER/$NAME.app"

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		cp ./tests/vendor/easymesh/easymesh.c $TEMP_FOLDER/easymesh.c

		sed -i.bak '1,33d'      $TEMP_FOLDER/$NAME.c
		sed -i.bak '140a(void)(i);'  $TEMP_FOLDER/$NAME.c
		sed -i.bak '449a(void)(d);(void)(ea);(void)(eb);'  $TEMP_FOLDER/$NAME.c
		sed -i.bak '634avoid '  $TEMP_FOLDER/$NAME.c
		sed -i.bak '647a(void)(e);(void)(xCa);(void)(yCa);(void)(px);(void)(py);'  $TEMP_FOLDER/$NAME.c

		sed -i.bak '252a(void)(en);(void)(n);(void)(ei);'  $TEMP_FOLDER/$NAME.c
		sed -i.bak '869a(void)(o);(void)(e2);'  $TEMP_FOLDER/$NAME.c
		sed -i.bak '740a(void)(iter);(void)(s);'  $TEMP_FOLDER/$NAME.c

		sed -i.bak '1050a(void)(N0);(void)(bound);(void)(xC);(void)(yC);'  $TEMP_FOLDER/$NAME.c
		sed -i.bak '1529a(void)(numb);(void)(x);(void)(y);(void)(ek);(void)(ej);(void)(ei);(void)(n);(void)(e);'  $TEMP_FOLDER/$NAME.c
		sed -i.bak '1655a(void)(numb);(void)(x);(void)(y);(void)(ek);(void)(ej);(void)(ei);(void)(n);(void)(e);'  $TEMP_FOLDER/$NAME.c

		sed -i.bak '1841a(void)(ans);(void)(g);'  $TEMP_FOLDER/$NAME.c

	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app

# transpilation of all projects
	echo "Transpile to $GO_FILE"
	$C4GO transpile                         \
		-clang-flag="-Wimplicit-int"        \
		-clang-flag="-Wreturn-type"         \
		-o="$GO_FILE"                       \
		$TEMP_FOLDER/$NAME.c

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

# debugging
if [ "$1" == "-d" ]; then
	# try to run
		echo " move to folder"
			cd ./testdata/$NAME/
		echo "step 1: create debug file"
			$C4GO debug $NAME.c
		echo "step 2: prepare output data"
		echo "step 3: prepare test script"
			echo '9 
0:  0.0  0.0    0.5    1
1:  5.0  0.0    0.5    2
2:  5.0  2.0    0.5    2
3:  4.0  3.0    0.5    3
4:  0.0  3.0    0.5    3
5:  1.0  1.0    0.9    4
6:  1.0  2.0    0.9    4
7:  2.0  2.0    0.9    4
8:  2.0  1.0    0.9    4
9 
0:  0  1    1
1:  1  2    2
2:  2  3    2
3:  3  4    3
4:  4  0    3
5:  5  6    4
6:  6  7    4
7:  7  8    4
8:  8  5    4' > inputc.d
			cat inputc.d > inputgo.d
		echo "step 4: run Go application"
			$C4GO transpile -o=debug.$NAME.go -clang-flag="-Wreturn-type" -clang-flag="-Wimplicit-int" debug.$NAME.c 2>&1 > c4go.output && echo "c4go ok" || echo "c4go not ok"
		echo "step 5: run C application"
			clang -o $NAME.c.app -lm debug.$NAME.c 2>&1 > clang.output && echo "clang ok" || echo "clang not ok"
			echo "" > debug.txt
			./$NAME.c.app inputc 2>&1 && echo "ok" || echo "not ok"
			cp debug.txt debug.c.txt
		echo "step 5: run Go application"
			go build -o $NAME.go.app	debug.$NAME.go
			echo "" > debug.txt
			./$NAME.go.app inputgo 2>&1 && echo "ok" || echo "not ok"
			cp debug.txt debug.go.txt
		echo "step 6"
			echo "-----------------------------"
			echo "debug"
			diff -y -t debug.c.txt debug.go.txt 2>&1  > debug.diff  && echo "ok" || echo "not ok"
			# cat debug.diff
			echo "-----------------------------"
		echo "step 7: move back"
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
