#!/bin/bash

workdir=`pwd`
tmpdir=/tmp/scraper.$$

mkdir -p $tmpdir/pdfs $tmpdir/txt $tmpdir/meta $tmpdir/m; pushd $tmpdir

cp $workdir/words .

curl -XGET "http://127.0.0.1:9200/scripts/script/_search?size=15000&pretty=true" | \
grep -B 1 url | egrep -v '"_type" : "script",' | grep -B 1 -v '"meta"' | egrep -v '^\-\-$' | \
perl -e 'while (<>) { s/^\s+"_id" : "([^"]+)",/$1/g; s/^\s+"_score" : .+, "_source" : //g; print }' | \
while read line; do 
	if [ "$foo" ]; then 
		echo "$foo\t$line"
		unset foo 
	else
		foo=$line
	fi
done | grep -v '"meta"' > urls

#split -l 200 urls ura

#for j in ura*; do
#	cat $j| perl -e 'while (<>) { s/{"id".+"url":"//g ; s/"}$//g;print }' | while read id url ; do
#		if [ -f "$workdir/pdfs/$id.pdf" ]; then 
#			echo exists
#		#else
#			#wget --connect-timeout=10 -O "pdfs/$id.pdf" "$url" --no-check-certificate
#		fi
#	done &
#done

#wait ${!}

cat urls | perl -e 'while (<>) { s/{"id".+"url":"//g ; s/"}$//g;print }' | while read id url ; do
	if ! [ -f "txt/$id.txt" ]; then
		if [ -f "$workdir/../pdfs/$id.pdf" ]; then
			pdftotext "$workdir/../pdfs/$id.pdf" "txt/$id.txt"
		fi
	fi
done

pushd txt
for j in *; do name=${j%.txt}; strings -- "$j" | head -c 155 | perl -e 'while (<>) { s/\n/ /g ; s/\s+/ /g ;print}' > "../meta/$name.meta"; done; popd

export GOPATH=/Users/asim/checkouts
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin

go run $workdir/f.go $tmpdir

popd

echo dir is $tmpdir/m
