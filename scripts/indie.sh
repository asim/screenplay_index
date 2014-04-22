#!/bin/bash


for tag in numbers a b c d e f g h i j k l m n o p q r s t u v w x y z; do
	curl -s "http://indiegroundfilms.com/resources/library/$tag/" | egrep 'href.*\.pdf' | while read line; do
		echo $line | sed -e 's@.*href="\([^"]\+\)"\([^>]\+\)\?>\([^<]\+\)<.*@url=\1\&title=\3@g'
	done
	sleep 5
done
