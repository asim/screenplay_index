#!/bin/bash

work=/home/asim/c/src/github.com/asim/screenplay_index

cd $work

echo -n > u

for j in scripts/joblo.sh scripts/august.sh scripts/awesome.sh scripts/bbc.sh scripts/drexel.sh scripts/indie.sh scripts/lee.sh scripts/screenplaydb.sh scripts/simply.sh ; do bash $j >> u; done

cat u | bash scripts/validator.sh | sort | uniq > r

./screenplay_index &

pid=$!

sleep 5

cat r | while read line; do curl -XPOST -d "$line" http://127.0.0.1:8081/_add; done

kill $pid
