#!/bin/bash

rm -rf eval4/results
mkdir eval4/results

for i in {0..9}
do
    go run github.com/victwj/freenet/eval4 > eval4/results/result$i.txt
done

python -u eval4/plot.py