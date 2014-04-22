#!/bin/bash

cat | while read line; do
	url=`echo $line | cut -f 1 -d \& | cut -f 2 -d =`

	if (curl -m 5 -s --head "$url" |egrep -q "Content-Type: application/pdf"); then
		echo $line
	fi
done
