#!/bin/bash

for f in $(find . -name go.mod)
do (
    cd $(dirname $f)
    go mod tidy -e
)

done
