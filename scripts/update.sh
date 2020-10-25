#!/bin/bash

work=/home/asim/c/src/github.com/asim/screenplays

cd $work

echo -n > u

for j in scripts/joblo.sh scripts/august.sh scripts/awesome.sh scripts/bbc.sh scripts/drexel.sh scripts/indie.sh scripts/lee.sh scripts/screenplaydb.sh scripts/simply.sh ; do bash $j >> u; done

cat u | bash scripts/validator.sh | sort | uniq > r

./screenplays &

pid=$!

sleep 5

cat r | while read line; do 
	url=`echo $line | cut -f 1 -d \&`
	title=`echo $line | cut -f 2 -d \&`
	curl -XPOST --data-urlencode "$url" -d "$title" http://127.0.0.1:8081/_add
done

kill $pid
