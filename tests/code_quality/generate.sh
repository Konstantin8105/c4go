#!/bin/bash

go build

# Generate code quality Go code
FILES='tests/code_quality/*.c'
for file in $FILES
do
  filename=$(basename "$file")
  ext="${filename#*.}"
  if [ "$ext" = "go.expected" ];  then
	  continue
  fi

  echo "Processing $file file..."
  filename=${file%.*}".go.expected"
  ./c4go transpile -o="$filename" -p="code_quality" $file
done
