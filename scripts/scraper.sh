#!/bin/bash

src=$1
file=$2

urlencode() {
    # urlencode <string>
 
    local length="${#1}"
    for (( i = 0 ; i < length ; i++ )); do
        local c="${1:i:1}"
        case "$c" in
            [a-zA-Z0-9.~_-]) printf "$c" ;;
            ' ') printf + ;;
            *) printf '%%%X' "'$c"
        esac
    done
}
 
urldecode() {
    # urldecode <string>
 
    local url_encoded="${1//+/ }"
    printf '%b' "${url_encoded//%/\x}"
}

dropbox() {
	sed 's@<a href@\n@g' $1 | sed 's@</a>@\n@g'  | grep pdf | awk '{print $1}' | sort | uniq | sed 's/=\?"//g' | sort | uniq | sed 's/www.dropbox.com/dl.dropboxusercontent.com/g' | grep pdf | sort | uniq | while read line; do title=`echo $line |awk -F / '{print $NF}'| sed 's/\.pdf$//g'`; title=`echo $(urldecode $title) | sed 's/\([^ ]\)_/\1 /g'`; echo "url=$line&title=$title"; done
}

if [ -z "$src" ]; then
	echo "need source {dropbox}"
	exit 1
fi

if [ -z "$file" ]; then
	echo "require file"
	exit 1
fi

if ! [ -f "$file" ]; then
	echo "cant find file"
	exit 1
fi

case "$src" in
	dropbox)
	dropbox $file
	;;
esac
