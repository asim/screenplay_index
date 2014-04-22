#!/bin/bash

# STATUS CODES

# 0 - OK
# 1 - Failed Download
# 2 - Not a PDF
# 3 - No metadata

url=$1
tmppdf=tmp.pdf
tmptxt=tmp.txt

if ! [ "${url: -4}" == ".pdf" ]; then
	echo "Not a pdf file"
	exit 2
fi

curl -m 10 -s $url > $tmppdf
if [ $? -ne 0 ]; then
	echo "Error retrieving $url"
	rm -f $tmppdf
	exit 1
fi

if ! (file -b $tmppdf | grep -q "PDF document"); then
	echo "Not a pdf file"
	rm -f $tmppdf
	exit 2
fi

pdftotext $tmppdf $tmptxt &> /dev/null
if [ $? -ne 0 ]; then
	echo "Error creating text file"
	rm -f $tmppdf $tmptxt
	exit 3
fi

meta=$(head -c 256 $tmptxt | tr "\\n" " ")
echo "$meta..."

rm -f $tmppdf $tmptxt

exit 0
