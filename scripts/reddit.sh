#!/bin/bash

urls() {
	curl -A "Mozilla/5.0 (compatible; scridxbot/0.1)" -s $1 | sed -e 's@</a>@\n\n@g' | grep "/r/Screenwriting" |sed -e 's@.*href="@@g' -e 's@".*@@g' -e 's@^/@http://www.reddit.com/@g' -e 's@\(?count=[0-9]\+&\)amp;@\1@g' | egrep "^http://" | egrep "http://www.reddit.com/r/Screenwriting/(comment|\?count)" | sort | uniq 
}

recurse() {
	while read line; do
		if (echo "$line" | egrep -q "http://www.reddit.com/r/Screenwriting/\?count=.*after"); then 
			sleep 5
			urls "$line" | recurse
			continue
		fi

		sleep 5
		curl -A "Mozilla/5.0 (compatible; scridxbot/0.1)" -s "$line" |sed -e 's@</a>@\n\n@g' | egrep "href.*\.pdf" | sed -e 's@.*href="@@g' -e 's@"[^>]\+>@ @g' -e 's@<[^>]\+>@@g' | egrep "\.pdf "
	done
}

urls "http://www.reddit.com/r/Screenwriting/" | recurse
