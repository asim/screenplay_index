#!/bin/bash

urls=( http://leethomson.myzen.co.uk/ http://www.zen134237.zen.co.uk/ )

parse() {
	local url=$1
	local file=`echo $url | sed 's@[^a-zA-Z0-9]@@g'`

	curl -s $url > /tmp/lee.tmp/$file

	# current dir pdfs
	egrep '\.pdf' /tmp/lee.tmp/$file | sed -e 's@.*href="\([^"]\+\)">\([^<]\+\)</a>.*@url=\1\&title=\2@g' | while read line; do
		local title=`echo $line | cut -f 2 -d \& | sed -e 's@_@ @g' -e 's@\.pdf$@@g'`
		local uri=`echo $line | cut -f 1 -d \& | sed "s@^url=@url=$url@g"`
		echo "$uri&$title"
	done

	egrep '\[DIR\]' /tmp/lee.tmp/$file | sed 's@.*href="\([^"]\+\)".*@\1@g' |sort |uniq | while read line; do
		if [ "$line" == "/" ]; then
			continue
		fi

		sleep 0.2
		parse "$url$line"
	done

	rm -f /tmp/lee.tmp/$file
}

if [ -d /tmp/lee.tmp ]; then
	rm -rf /tmp/lee.tmp
fi

mkdir /tmp/lee.tmp

for url in ${urls[@]}; do
	parse $url
done

rm -rf /tmp/lee.tmp
