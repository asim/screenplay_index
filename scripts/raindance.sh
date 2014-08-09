#!/bin/bash

curl -s http://www.raindance.org/site/scripts/ | grep "\.pdf" | sed 's@.*href="@title=@g' | sed 's@"> \?@\&url=http://www.raindance.org/site/scripts/@g' | sed 's@</a></li>$@@g' | sed 's@.pdf\&url@\&url@g' | while read line; do title=`echo $line | cut -f 1 -d \& | sed 's@_@ @g' |sed 's@%20@ @g'`; url=`echo $line | cut -f 2 -d \&`; echo "$url&$title" ; done
