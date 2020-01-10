#!/bin/bash

rm -rf eval1/results
mkdir eval1/results

for i in {0..9}
do
    go run github.com/victwj/freenet/eval1 > eval1/results/result$i.txt
done

python -u eval1/plot.py