#!/bin/bash

set -e

mkdir -p ./testdata/

./scripts/9wm.sh
# ./scripts/cis71.sh
./scripts/ed.sh
./scripts/neateqn.sh
./scripts/neatmkfn.sh
./scripts/neatpost.sh
./scripts/neatrefer.sh
./scripts/neatroff.sh
./scripts/neatvi.sh
