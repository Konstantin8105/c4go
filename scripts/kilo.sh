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
		sed -i.bak '53ivoid debug(char *msg){FILE*file;file=fopen("./debug.txt","a");if(file==NULL){exit(53);};fprintf(file,"%s\\n",msg);fclose(file);}'  $TEMP_FOLDER/kilo.c
		# add debug information
		sed -i.bak '726i{char buffer[500];sprintf(buffer,"line%d: filecol=%d",__LINE__,filecol);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '725i{char buffer[500];sprintf(buffer,"line%d: filerow=%d",__LINE__,filerow);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '1273i{char buffer[500];sprintf(buffer,"%d",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '1272i{char buffer[500];sprintf(buffer,"%d",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '1271i{char buffer[500];sprintf(buffer,"%d",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '1268i{char buffer[500];sprintf(buffer,"%d",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '1267i{char buffer[500];sprintf(buffer,"%d",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '1266i{char buffer[500];sprintf(buffer,"%d",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '1265i{char buffer[500];sprintf(buffer,"%d",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '1264i{char buffer[500];sprintf(buffer,"%d",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '1263i{char buffer[500];sprintf(buffer,"%d",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c

		sed -i.bak '1173i{char buffer[500];sprintf(buffer,"%d: key %d",__LINE__, c );debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '575i{char buffer[500];sprintf(buffer,"%d: sizeof : %d. numrows : %d.  E.row size : %d",__LINE__, sizeof(erow) ,E.numrows, sizeof(erow) * (E.numrows+1)  );debug(buffer);}' $TEMP_FOLDER/kilo.c
		
		sed -i.bak '710i{char buffer[500];sprintf(buffer,"%d: row %d. col %d. E.numrows %d",__LINE__, filerow, filecol, E.numrows );debug(buffer);}' $TEMP_FOLDER/kilo.c

		sed -i.bak '588i{char buffer[500];sprintf(buffer,"%d: len %d. at %d. chars: `%s`",__LINE__, len, at, E.row[at].chars );debug(buffer);}' $TEMP_FOLDER/kilo.c

		sed -i.bak '1175i{char buffer[500];sprintf(buffer,"%d: editorProcessKeypress",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '818i{char buffer[500];sprintf(buffer,"%d: editorSave",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '787i{char buffer[500];sprintf(buffer,"%d: editorOpen",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '747i{char buffer[500];sprintf(buffer,"%d: editorDelChar",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '708i{char buffer[500];sprintf(buffer,"%d: editorInsertNewLine",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '686i{char buffer[500];sprintf(buffer,"%d: editorInsertChar",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '604i{char buffer[500];sprintf(buffer,"%d: editorDelRow",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '574i{char buffer[500];sprintf(buffer,"%d: editorInsertRow",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '241i{char buffer[500];sprintf(buffer,"%d: editorReadKey",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c

		sed -i.bak '809i{char buffer[500];sprintf(buffer,"%d: loop in editorOpen : `%s` %d",__LINE__, line, linelen);debug(buffer);}' $TEMP_FOLDER/kilo.c
		
		sed -i.bak '804i{char buffer[500];sprintf(buffer,"%d: editorOpen: filenames: `%s` `%s`",__LINE__, filename, E.filename);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '815i{char buffer[500];sprintf(buffer,"%d: editorOpen: out of loop",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		
		sed -i.bak '546i{char buffer[500];sprintf(buffer,"%d: editorUpdateRow",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '371i{char buffer[500];sprintf(buffer,"%d: editorUpdateSyntax",__LINE__);debug(buffer);}' $TEMP_FOLDER/kilo.c
		sed -i.bak '576i{char buffer[500];sprintf(buffer,"%d: editorInsertRow: at %d, s: `%s`, len %d, E.numrows : %d",__LINE__,at,s,len,E.numrows);debug(buffer);}' $TEMP_FOLDER/kilo.c

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
	# try to run
		echo "try to run"
		echo "step 1"
			cd ./testdata/kilo/
		echo "step 2"
			echo "" > output.go.txt
			echo "" > output.c.txt
		echo "step 3"
			echo -e 'Hello my dear friend\x0D\x13\x11' > script.txt
		echo "step 4"
			echo "" > debug.txt
			cat script.txt | ./kilo.app   output.go.txt 2>&1 && echo "ok" || echo "not ok"
			cp debug.txt debug.go.txt	
			
			gcc -o kilo.c.app kilo.c
			echo "" > debug.txt
			cat script.txt | ./kilo.c.app output.c.txt  2>&1 && echo "ok" || echo "not ok"
			cp debug.txt debug.c.txt	
		echo "step 5"
			echo "-----------------------------"
			echo "debug"
			diff -y -t debug.c.txt debug.go.txt 2>&1  && echo "ok" || echo "not ok"
			echo "-----------------------------"
			echo "output"
			diff -y -t output.c.txt output.go.txt 2>&1  && echo "ok" || echo "not ok"
			echo "-----------------------------"
		echo "step 6"
			cd ../../

# Arguments menu
echo "    -s for show detail of Go build errors"
if [ "$1" == "-s" ]; then
	# show go build warnings	
		# c4go warnings
			cat $GO_FILE | grep "^// Warning" | sort | uniq
		# show amount error from `go build`:
			go build -o $GO_APP -gcflags="-e" $GO_FILE 2>&1
fi
