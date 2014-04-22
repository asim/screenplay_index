#!/bin/bash

curl -s http://www.simplyscripts.com/sitemap.html | egrep 'href.*\.html' | sed -e 's/.*href="\([^"]\+\)".*/\1/g'| cut -f 1 -d \# |sort | uniq > /tmp/simply.tmp

cat /tmp/simply.tmp | while read line; do
	sleep 0.3
	curl -m 10 -s $line | sed -e 's@</a>@</a>\n\n@g' |egrep '\.pdf"' | sed -e 's@.*href="\([^"]\+\.pdf\)"\([^>]\+\)\?>\([^<]\+\)<.*@url=\1\&title=\3@g' -e 's@.*href="\([^"]\+\.pdf\)"\([^>]\+\)\?>@url=\1\&title=@g' -e 's@<[^>]\+>@@g' | sort | uniq | while read x; do
		echo $x >> /tmp/simple.list.tmp
	done
done

cat  /tmp/simple.list.tmp | while read line; do
	url=`echo $line |cut -f 1 -d \& | cut -f 2 -d =`
	if (curl -m 10 -s --head $url|grep -q "Content-Type: application/pdf"); then
		echo $line | strings
	fi
done

rm -rf /tmp/simply.tmp /tmp/simple.list.tmp
