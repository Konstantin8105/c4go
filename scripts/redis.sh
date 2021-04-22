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

# AST tree
	echo "AST generate to $GO_FILE"
if [ "$1" == "-a" ]; then
	$C4GO ast                                                \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/hiredis"   \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/jemalloc"  \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/linenoise" \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/lua"       \
		$TEMP_FOLDER/$VERSION/src/redis-cli.c
fi

# transpilation of all projects
	echo "Transpile to Go"
	$C4GO transpile                                          \
		-s                                                   \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/hiredis"   \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/jemalloc"  \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/linenoise" \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/lua"       \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/src/"           \
		-o="$TEMP_FOLDER/redis_cli.go"                       \
		$TEMP_FOLDER/$VERSION/src/redis-cli.c


	$C4GO transpile                                          \
		-s                                                   \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/hiredis"   \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/jemalloc"  \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/linenoise" \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/lua"       \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/src/"           \
		-o="$TEMP_FOLDER/dict.go"                            \
		$TEMP_FOLDER/$VERSION/src/dict.c

	$C4GO transpile                                              \
		-s                                                       \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/hiredis"       \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/jemalloc"      \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/linenoise"     \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/deps/lua"           \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/src/"           \
		-o="$TEMP_FOLDER/sds.go"                                 \
		$TEMP_FOLDER/$VERSION/src/sds.c

	$C4GO transpile                                          \
		-o="$TEMP_FOLDER/linenoise.go"                       \
		-clang-flag="-I$TEMP_FOLDER/$VERSION/src/"           \
		$TEMP_FOLDER/$VERSION/deps/linenoise/linenoise.c     \
		$TEMP_FOLDER/$VERSION/deps/linenoise/example.c

# Arguments menu
echo "    -s for show detail of Go build errors"

# transpilation each file
	for f in $TEMP_FOLDER/*.go; do
			# iteration by C projects
				echo "***** transpilation folder : $f"
			# show warnings comments in Go source
				export FILE="$f"
				echo "Calculate warnings : $FILE"
				WARNINGS=`cat $FILE | grep "^// Warning" | sort | uniq | wc -l`
				echo "		After transpiling : $WARNINGS warnings."
			# show amount error from `go build`:
				echo "Build to $GO_APP file"
				WARNINGS_GO=`go build -o $GO_APP -gcflags="-e" $f 2>&1 | wc -l`
				echo "		Go build : $WARNINGS_GO warnings"
			# amount unsafe
				UNSAFE=`cat $FILE | grep "unsafe\." | wc -l`
				echo "		Unsafe   : $UNSAFE"

			if [ "$1" == "-s" ]; then
				# show go build warnings	
					# c4go warnings
						cat $f | grep "^// Warning" | sort | uniq
					# show amount error from `go build`:
						go build -o $GO_APP -gcflags="-e" $f 2>&1 && echo "OK" || echo "NOK"
			fi
	done

