#!/bin/bash

rm -rf eval2/results
mkdir eval2/results

for i in {0..9}
do
    go run github.com/victwj/freenet/eval2 > eval2/results/result$i.txt
done

python -u eval2/plot.py