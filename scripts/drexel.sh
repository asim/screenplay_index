#!/bin/bash

url=http://www.pages.drexel.edu/~ina22/splaylib/

curl -s ${url}dwpscreenplay-library.htm | egrep "HREF.*\.pdf" | sed 's@.*HREF="\([^"]\+\)".*@\1@g' | egrep ^Screenplay | while read line; do
	script=`echo $line | sed -e "s@$url@@g"`
	title=`echo $script | sed -e 's@^Screenplay-@@g' -e 's@[-_]@ @g' -e 's@.pdf$@@g' -e 's@%20@ @g'`
	echo "url=$url$script&title=$title"
done
