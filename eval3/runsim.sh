#!/bin/bash

rm -rf eval3/results
mkdir eval3/results

for i in {0..9}
do
    go run github.com/victwj/freenet/eval3 > eval3/results/result$i.txt
done

python -u eval3/plot.py