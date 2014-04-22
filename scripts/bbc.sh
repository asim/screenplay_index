#!/bin/bash


get() {
	curl -s $1 | sed 's@</a>@\n@g' | egrep "href.*writersroom/scripts/" | \
	sed 's@.*href="\([^"]\+\)">\(.\+\)@\1\t\2@g'| grep -v "<" | \
	while read url title; do
		sleep 0.2
		curl -s "http://www.bbc.co.uk$url" | sed 's@</a>@\n@g' | egrep "href.*\.pdf" |sed 's@.*href="\([^"]\+\)".*@\1@g' | while read line; do
			echo "url=$line&title=$title"
		done
	done
}

for page in {1..16}; do
	get http://www.bbc.co.uk/writersroom/scripts/search/page/$page
	sleep 1
done
