#!/bin/bash

for j in " " {2..11} ; do
	curl -s http://www.screenplaydb.com/film/all/$j/ | egrep "screenplaydb.com/film/(scripts|download)" | sed 's@.*film/scripts\([^>]\+\)>@@g' | sed 's@.*href="\([^"]\+\)".*@\1@g' | sed 's@<[^>]\+>@@g' | while read line; do

		if [ "$foo" ]; then
			u=$(curl -I -s $line |grep Location | grep pdf | sed 's/^Location: //g' |strings)
			if [ "$u" ]; then
				echo "url=$u&title=$foo"
			fi
			unset foo
		else
			foo=$line
		fi
	done
done
