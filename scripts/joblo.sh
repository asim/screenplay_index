#!/bin/bash

i=0
for j in {1..20}; do
	curl -q -s "http://www.joblo.com/movie-screenplays-scripts/archives/?mode=images&from=$i&year=&first_letter=&author=&rating=&order_by=" | \
	grep pdf | sed -e 's@.*<a title="\([^"]\+\)" href="\([^"]\+\)".*@url=\2\&title=\1@g' | sort | uniq
	i=$(($i+65))
	sleep 1
done
