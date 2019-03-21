#!/bin/bash

set -e

mkdir -p ./testdata/

export OUTPUT_FILE="./testdata/scripts.txt"
export VERIFICATION_FILE="./scripts/scripts.txt"

echo "" > $OUTPUT_FILE

./scripts/9wm.sh		| grep -E 'warning|unsafe|Unsafe' | tee -a $OUTPUT_FILE
echo "" >> $OUTPUT_FILE
./scripts/cis71.sh	| grep -E 'warning|unsafe|Unsafe' | tee -a $OUTPUT_FILE
echo "" >> $OUTPUT_FILE
./scripts/ed.sh			| grep -E 'warning|unsafe|Unsafe' | tee -a $OUTPUT_FILE
echo "" >> $OUTPUT_FILE
./scripts/neateqn.sh	| grep -E 'warning|unsafe|Unsafe' | tee -a $OUTPUT_FILE
echo "" >> $OUTPUT_FILE
./scripts/neatmkfn.sh	| grep -E 'warning|unsafe|Unsafe' | tee -a $OUTPUT_FILE
echo "" >> $OUTPUT_FILE
./scripts/neatpost.sh	| grep -E 'warning|unsafe|Unsafe' | tee -a $OUTPUT_FILE
echo "" >> $OUTPUT_FILE
./scripts/neatrefer.sh	| grep -E 'warning|unsafe|Unsafe' | tee -a $OUTPUT_FILE
echo "" >> $OUTPUT_FILE
./scripts/neatroff.sh	| grep -E 'warning|unsafe|Unsafe' | tee -a $OUTPUT_FILE
echo "" >> $OUTPUT_FILE
./scripts/neatvi.sh		| grep -E 'warning|unsafe|Unsafe' | tee -a $OUTPUT_FILE
echo "" >> $OUTPUT_FILE

# Arguments menu
echo "    -u update scripts result"
if [ "$1" == "-u" ]; then
	cat $OUTPUT_FILE > $VERIFICATION_FILE
else
	diff $OUTPUT_FILE $VERIFICATION_FILE 2>&1
fi
